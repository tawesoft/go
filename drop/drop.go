package drop

import (
    "fmt"
    "os"
    "os/exec"
    "runtime"
    "syscall"
)

// Inheritable
type Inheritable interface {
    Name() string
    Open() (*os.File, error)
    Inherit(*os.File) error
}

// Drop works in two ways, depending on the privileges of the current process.
//
// As the superuser (root), drop executes a new copy of the running program as
// the given user, then exits. The current stdio streams, and zero-or-more
// "inheritable" files, are persisted and inherited by the new process.
//
// As a non-superuser, drop inherits the files from the calling process.
//
// NOTE: it is not a good idea to open any resource, other than those passed
// as Inheritable, prior to a call to drop. Open as few resources, and do as
// little as possible, while root. Call Drop as early as possible.
func Drop(
    username string,
    files ... Inheritable,
) error {
    if runtime.GOOS != "linux" {
        return fmt.Errorf("unsupported: Drop only works on Linux")
    }

    closeAll := func(handles []*os.File) {
        for _, i := range handles {
            i.Close()
        }
    }

    handles := make([]*os.File, 0, len(files))

    if isSuperuser() { // drop
        uid, gid, groups, err := userLookup(username)
        if err != nil { return fmt.Errorf("user lookup error for %s: %v", username, err) }
        if uid == 0 { return fmt.Errorf("cannot drop to root") }

        for _, file := range files {
            handle, err := file.Open()
            if err != nil {
                closeAll(handles)
                return fmt.Errorf("error opening file while privileged: %v", err)
            }
            handles = append(handles, handle)
        }

        args := os.Args
        cmd        := exec.Command(args[0], args[1:]...)
        cmd.Stdin  = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        cmd.ExtraFiles = handles
        cmd.SysProcAttr = &syscall.SysProcAttr{}
        cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uid, Gid: gid, Groups: groups}

        err = cmd.Start()
        if err != nil {
            closeAll(handles)
            return fmt.Errorf("error dropping privileges: %v", err)
        }

        closeAll(handles)

        // close immediately
        // os.Exit(0)

        // alternatively supervise (waste of memory, use system services)
        err = cmd.Wait()
        if err != nil {
            return fmt.Errorf("child process exited with error: %v", err)
        }
        os.Exit(0)

    } else { // inherit

        for i := 0; i < len(files); i++ {
            handle := os.NewFile(uintptr(3 + i), files[i].Name())
            if handle == nil {
                closeAll(handles)
                return fmt.Errorf("error inheriting file %s", files[i].Name())
            }
            handles = append(handles, handle)
        }

        for i := 0; i < len(files); i++ {
            err := files[i].Inherit(handles[i])
            if err != nil {
                closeAll(handles)
                return fmt.Errorf("error inheriting file %s: %v", files[i].Name(), err)
            }
        }
    }

    return nil
}
