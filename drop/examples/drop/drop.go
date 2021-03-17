// Opens privileged files and ports as root, then drops privileges
package main

import (
    "fmt"
    "net"
    "os"

    "tawesoft.co.uk/go/drop"
)

// Define structures and methods that meets the start.Inheritable interface
// by implementing Name(), Open(), and Inherit()...

// InheritableFile is a file handle that survives a dropping of
// privileges.
type InheritableFile struct {
    Path  string
    Flags int // e.g.  os.O_RDWR|os.O_CREATE
    Perm  os.FileMode // e.g. 0600

    handle *os.File
}

func (h InheritableFile) Name() string {
    return h.Path
}

func (h InheritableFile) Open() (*os.File, error) {
    return os.OpenFile(h.Path, h.Flags, h.Perm)
}

func (h *InheritableFile) Inherit(f *os.File) error {
    f.Seek(0, 0)
    h.handle = f
    return nil
}

// InheritableNetListener is a net.Listener that survives a dropping of
// privileges.
type InheritableNetListener struct {
    Network string
    Address string
    handle net.Listener
}

func (h InheritableNetListener) Name() string {
    return fmt.Sprintf("(%s) %s", h.Network, h.Address)
}

func (h InheritableNetListener) Open() (*os.File, error) {
    nl, err := net.Listen(h.Network, h.Address)
    if err != nil { return nil, err }
    defer nl.Close()

    fl, err := nl.(*net.TCPListener).File()
    if err != nil { return nil, err }
    return fl, nil
}

func (h *InheritableNetListener) Inherit(f *os.File) error {
    defer f.Close()
    fl, err := net.FileListener(f)
    if err != nil { return err }
    h.handle = fl
    return nil
}

func main() {
    if len(os.Args) < 2 {
        panic(fmt.Sprintf("USAGE: sudo %s username\n", os.Args[0]))
    }

    // what user to drop to
    username := os.Args[1]

    // resources to be opened as root and persist after privileges are dropped
    privilegedFile := &InheritableFile{"/tmp/privileged-file-example", os.O_RDWR|os.O_CREATE, 0600, nil}
    privilegedPort := &InheritableNetListener{"tcp4", "127.0.0.1:81", nil}

    // If the program is run as root, open privileged resources as root, then
    // start a child process as `username` that inherits these resources and
    // the parent process's stdio, and immediately exit.
    //
    // If the program is run as non-root, inherit these resources and continue.
    shouldExit, err := drop.Drop(username, privilegedFile, privilegedPort)
    if err != nil {
        panic(fmt.Sprintf("error dropping privileges (try running as root): %v", err))
    }
    if shouldExit { return }

    // At this point, the program is no longer running as root, but it still
    // has access to these privileged resources.

    // do things with privilegedFile
    privilegedFile.handle.WriteString("hello world\n")
    privilegedFile.handle.Close()

    // do things with privilegedPort
    privilegedPort.handle.Close()
}
