package drop

import (
    "fmt"
    "net"
    "os"
    "strings"
)

// Inheritable interface describes a File-based handle that can be passed
// between a parent and child process across a dropping of privileges.
type Inheritable interface {
    Name() string
    Open() (*os.File, error)
    Inherit(*os.File) error
    Close()
}

// InheritableFile is a file handle that survives a dropping of
// privileges.
type InheritableFile struct {
    Path  string
    Flags int // e.g.  os.O_RDWR|os.O_CREATE
    Perm  os.FileMode // e.g. 0600
    Handle *os.File // nil until Inherit returns
}

func (h InheritableFile) Name() string {
    return h.Path
}

func (h InheritableFile) Open() (*os.File, error) {
    return os.OpenFile(h.Path, h.Flags, h.Perm)
}

func (h *InheritableFile) Inherit(f *os.File) error {
    f.Seek(0, 0)
    h.Handle = f
    return nil
}

func (h InheritableFile) Close() {
    h.Handle.Close()
}

// InheritableNetListener is a net.Listener that survives a dropping of
// privileges.
//
// Note that On JS and Windows, the File method of most Listeners are not
// implemented, so this will not work.
type InheritableNetListener struct {
    Network string
    Address string
    Handle net.Listener // nil until Inherit returns
}

func (h InheritableNetListener) Name() string {
    return fmt.Sprintf("(%s) %s", h.Network, h.Address)
}

func (h InheritableNetListener) Open() (*os.File, error) {
    nl, err := net.Listen(h.Network, h.Address)
    if err != nil { return nil, err }
    defer nl.Close()

    if strings.HasPrefix(h.Network, "tcp") {
        if fl, err := nl.(*net.TCPListener).File(); err != nil {
            return fl, nil
        } else {
            return nil, fmt.Errorf("error obtaining TCPListener File handle: %+v", err)
        }
    } else if strings.HasPrefix(h.Network, "unix") {
        if fl, err := nl.(*net.UnixListener).File(); err != nil {
            return fl, nil
        } else {
            return nil, fmt.Errorf("error obtaining UnixListener File handle: %+v", err)
        }
    }

    return nil, fmt.Errorf("unsupported network %v", h.Network)
}

func (h *InheritableNetListener) Inherit(f *os.File) error {
    defer f.Close()
    fl, err := net.FileListener(f)
    if err != nil { return err }
    h.Handle = fl
    return nil
}

func (h InheritableNetListener) Close() {
    h.Handle.Close()
}
