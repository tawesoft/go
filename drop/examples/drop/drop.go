// Opens privileged files and ports as root, then drops privileges
package main

import (
    "fmt"
    "os"

    "tawesoft.co.uk/go/drop"
)

// Define structures and methods that meets the drop.Inheritable interface
// by implementing Name(), Open(), Inherit() and Close()...

// These are implemented here for example purposes - but you can use
// the builtins from drop instead.

// InheritableFile is a file handle that survives a dropping of
// privileges.
type InheritableFile struct {
    Path  string
    Flags int // e.g.  os.O_RDWR|os.O_CREATE
    Perm  os.FileMode // for os.O_CREATE, e.g. 0600

    handle *os.File
}

func (h InheritableFile) String() string {
    return h.Path
}

func (h *InheritableFile) Open() (*os.File, error) {
    f, err := os.OpenFile(h.Path, h.Flags, h.Perm)
    h.handle = f
    return f, err
}

func (h *InheritableFile) Inherit(f *os.File) error {
    f.Seek(0, 0)
    h.handle = f
    return nil
}

func (h *InheritableFile) Close() error {
    return h.handle.Close()
}

func main() {
    if len(os.Args) < 2 {
        panic(fmt.Sprintf("USAGE: sudo %s username\n", os.Args[0]))
    }

    // what user to drop to
    username := os.Args[1]

    // resources to be opened as root and persist after privileges are dropped
    privilegedFile := &InheritableFile{"/tmp/privileged-file-example", os.O_RDWR|os.O_CREATE, 0600, nil}
    privilegedPort := drop.NewInheritableTCPListener(":81")

    // If the program is run as root, open privileged resources as root, then
    // start a child process as `username` that inherits these resources and
    // the parent process's stdio, and immediately exit.
    //
    // If the program is run as non-root, inherit these resources and continue.
    shouldExit, closer, err := drop.Drop(username, privilegedFile, privilegedPort)
    if err != nil {
        panic(fmt.Sprintf("error dropping privileges (try running as root): %v", err))
    }
    defer closer()
    if shouldExit { return }

    // At this point, the program is no longer running as root, but it still
    // has access to these privileged resources.

    // do things with privilegedFile
    privilegedFile.handle.WriteString("hello world\n")
    privilegedFile.Close()

    // do things with privilegedPort
    // privilegedPort.handle ....
    privilegedPort.Close()
}
