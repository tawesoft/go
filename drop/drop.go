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
// the given user, blocks until that child program exits, then returns.
//
// The current stdio streams, and zero-or-more "inheritable" files, are
// persisted and inherited by the new process.
//
// As a non-superuser, drop inherits the files from the calling process and
// continues without blocking or exiting.
//
// If the first return argument is true, the caller should exit immediately
// (e.g. by returning from main(), with os.Exit(), etc.)
func Drop(username string, files ... Inheritable,) (bool, error) {
    return dropArg(username, true, files...)
}

// Detach works like Drop, however in the case of a superuser creating a new
// running program, it does not block until the child process exits.
//
// Also, Stdin, Stdout, and Stderr are not attached - you will want to use
// logging, instead.
//
// You'll probably want some mechanism (e.g. PID files) for terminating
// long-lived child processes too, because they will no longer be interruptable
// from the terminal.
//
// If the first return argument is true, the caller should exit immediately
// (e.g. by returning from main(), with os.Exit(), etc.)
//
// Not sure there's any way to do this that is safe against PID reuse! So
// disabled for now
//
//func Detach(username string, files ... Inheritable) (bool, error) {
//    return dropArg(username, false, files...)
//}

func dropArg(username string, supervise bool, files ... Inheritable) (bool, error) {
    if runtime.GOOS != "linux" {
        return false, fmt.Errorf("unsupported: Drop only works on Linux")
    }

    closeAll := func(handles []*os.File) {
        for _, i := range handles {
            i.Close()
        }
    }

    handles := make([]*os.File, 0, len(files))

    if isSuperuser() { // drop
        uid, gid, groups, err := userLookup(username)
        if err != nil { return false, fmt.Errorf("user lookup error for %s: %v", username, err) }
        if uid == 0 { return false, fmt.Errorf("cannot drop to root") }

        for _, file := range files {
            handle, err := file.Open()
            if err != nil {
                closeAll(handles)
                return false, fmt.Errorf("error opening file while privileged: %v", err)
            }
            handles = append(handles, handle)
        }

        args := os.Args
        cmd        := exec.Command(args[0], args[1:]...)

        if supervise {
            cmd.Stdin  = os.Stdin
            cmd.Stdout = os.Stdout
            cmd.Stderr = os.Stderr
        }

        cmd.ExtraFiles = handles
        cmd.SysProcAttr = &syscall.SysProcAttr{}
        cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uid, Gid: gid, Groups: groups}

        err = cmd.Start()
        if err != nil {
            closeAll(handles)
            return false, fmt.Errorf("error dropping privileges: %v", err)
        }

        closeAll(handles)

        if supervise {
            err = cmd.Wait()
            if err != nil {
                return false, fmt.Errorf("child process exited with error: %v", err)
            }
        }

        return true, nil
    } else { // inherit

        for i := 0; i < len(files); i++ {
            handle := os.NewFile(uintptr(3 + i), files[i].Name())
            if handle == nil {
                closeAll(handles)
                return false, fmt.Errorf("error inheriting file %s", files[i].Name())
            }
            handles = append(handles, handle)
        }

        for i := 0; i < len(files); i++ {
            err := files[i].Inherit(handles[i])
            if err != nil {
                closeAll(handles)
                return false, fmt.Errorf("error inheriting file %s: %v", files[i].Name(), err)
            }
        }
    }

    return false, nil
}
