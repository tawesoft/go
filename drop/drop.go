package drop

import (
    "context"
    "fmt"
    "os"
    "os/exec"
    "os/signal"
    "runtime"
    "syscall"
)

// Drop works in two ways, depending on the privileges of the current process.
//
// As the superuser (root), drop executes a new copy of the running program as
// the given user, blocks until that child program exits, then returns.
//
// As the child process, drop returns without blocking. The current stdio
// streams, and zero-or-more "inheritable" files, are persisted and inherited
// by the new process.
//
// If the first return argument is true, the caller is the root process and
// should exit immediately (e.g. by returning from main(), with os.Exit(), etc.)
//
// The second return argument is a function. It is nil if the third argument
// is an error. In the child process, this calls the Close() method on each of
// the supplied Inheritable files. The caller may either Close() each
// Inheritable manually or typically simply defer this returned function. The
// returned []error argument from this function contains the errors returned by
// Inheritable Close() methods. In the root process, this returned function is
// always a no-op and does not need to be called.
func Drop(username string, files ... Inheritable,) (bool, func() []error, error) {
    if runtime.GOOS == "windows" { panic("drop does not work on Windows") }

    // no way to exit without supervising in a way that is safe against PID
    // reuse.
    const supervise = true

    closer := func(files ... Inheritable) []error {
        errors := make([]error, 0)
        for _, i := range files {
            err := i.Close()
            if err != nil { errors = append(errors, err) }
        }
        return errors
    }

    handles := make([]*os.File, 0, len(files))

    if IsSuperuser() { // drop
        uid, gid, groups, err := UserLookup(username)
        if err != nil { return false, nil, fmt.Errorf("user lookup error for %s: %v", username, err) }
        if uid == 0 { return false, nil, fmt.Errorf("cannot drop to root") }

        defer closer(files...)

        // convert []int groups to []uint32
        // for cmd.SysProcAttr.Credential
        ugroups := make([]uint32, len(groups))
        for i, _ := range groups {
            ugroups[i] = uint32(groups[i])
        }

        for _, file := range files {
            handle, err := file.Open()
            if err != nil {
                return false, nil, fmt.Errorf("error opening file while privileged: %v", err)
            }
            handles = append(handles, handle)
        }

        args := os.Args
        cmd := exec.Command(args[0], args[1:]...)

        if supervise {
            cmd.Stdin  = os.Stdin
            cmd.Stdout = os.Stdout
            cmd.Stderr = os.Stderr
        }

        cmd.ExtraFiles = handles
        cmd.SysProcAttr = &syscall.SysProcAttr{}
        cmd.SysProcAttr.Credential = &syscall.Credential{
            Uid: uint32(uid),
            Gid: uint32(gid),
            Groups: ugroups,
        }

        // Let the parent process recover if the child process is killed
        // so that we can
        // e.g. so that we can close any open sockets
        ctx, cancel := context.WithCancel(context.Background())
        defer cancel()

        go func(ctx context.Context) {
            sigchan := make(chan os.Signal, 1)
            signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

            select {
                case <- sigchan:
                case <- ctx.Done():
            }

            signal.Stop(sigchan)
        }(ctx)

        err = cmd.Start()
        if err != nil {
            return false, nil, fmt.Errorf("error dropping privileges: %v", err)
        }

        if supervise {
            err = cmd.Wait()
            if err != nil {
                return false, nil, fmt.Errorf("child process exited with error: %v", err)
            }
        }

        closer(files...)
        return true, func() []error { return []error{} }, nil
    } else { // inherit

        for i := 0; i < len(files); i++ {
            name := files[i].String()
            handle := os.NewFile(uintptr(3 + i), name)
            if handle == nil {
                closer(files...)
                return false, nil, fmt.Errorf("missing file handle for %s", name)
            }
            handles = append(handles, handle)
        }

        for i := 0; i < len(files); i++ {
            name := files[i].String()
            err := files[i].Inherit(handles[i])
            if err != nil {
                closer(files...)
                return false, nil, fmt.Errorf("error inheriting file %s %v", name, err)
            }
        }

        return false, func() []error { return closer(files...) }, nil
    }
}
