package drop

import (
    "errors"
    "fmt"
    "net"
    "os"
    "strings"
)

// Inheritable interface describes a File-based handle that can be passed
// between a parent and child process across a dropping of privileges.
type Inheritable interface {
    // String returns some description of the resource.
    String() string

    // Open is called by the root process
    Open() (*os.File, error)

    // Inherit is called by the child process
    Inherit(*os.File) error

    // Close closes any resource generated by a call to Open or Inherit.
    Close() error
}

// InheritableFile is a file handle that survives a dropping of
// privileges.
type InheritableFile struct {
    path   string
    flags  int         // e.g.  os.O_RDWR|os.O_CREATE
    uid    int
    gid    int
    mode   os.FileMode // e.g. 0600
    handle *os.File    // nil until Inherit returns
}

// NewInheritableFile returns an InhertiableFile that wraps a file handle that
// survives a dropping of privileges.
//
// The parameters uid, gid, and mode, unless equal to -1, -1, or 0, set the
// user, group, and permissions of the file after it has been opened.
//
// WARNING: if these values are supplied from a config file, that config file
// should be writable to root or system accounts only - otherwise, an attacker
// may edit the config file in such a way as to set the permissions of
// arbitrary files.
func NewInheritableFile(path string, flags int, uid int, gid int, mode os.FileMode) *InheritableFile {
    return &InheritableFile{
        path:   path,
        flags:  flags,
        uid:    uid,
        gid:    gid,
        mode:   mode,
    }
}

func (h InheritableFile) String() string {
    return fmt.Sprintf("<InheritableFile %q>", h.path)
}

func (h *InheritableFile) Open() (*os.File, error) {
    f, err := os.OpenFile(h.path, h.flags, 0000)
    if err != nil { return nil, err }

    err = rootSetMode(h.path, h.uid, h.gid, h.mode)
    if err != nil {
        f.Close()
        return nil, err
    }

    h.handle = f

    return f, nil
}

func (h *InheritableFile) Inherit(f *os.File) error {
    f.Seek(0, 0)
    h.handle = f
    return nil
}

func (h InheritableFile) Handle() *os.File {
    return h.handle
}

func (h InheritableFile) Close() error {
    if h.handle == nil { return nil }
    return h.handle.Close()
}

// InheritableNetListener is a net.Listener that survives a dropping of
// privileges.
//
// Note that On JS and Windows, the File method of most Listeners are not
// implemented, so this will not work.
type InheritableNetListener struct {
    network string
    address string
    uid  int
    gid  int
    mode os.FileMode
    handle net.Listener
    fileHandle *os.File
}

func NewInheritableTCPListener(address string) *InheritableNetListener {
    return &InheritableNetListener{
        network: "tcp",
        address: address,
    }
}

// NewInheritableUnixListener returns an InheritableNetListener for a UNIX socket.
//
// The parameters uid, gid, and mode, unless equal to -1, -1, or 0, set the
// user, group, and permissions of the socket after it has been opened.
//
// WARNING: if these values are supplied from a config file, that config file
// should be writable to root or system accounts only - otherwise, an attacker
// may edit the config file in such a way as to set the permissions of
// arbitrary files.
func NewInheritableUnixListener(address string, uid int, gid int, mode os.FileMode) *InheritableNetListener {
    return &InheritableNetListener{
        network: "unix",
        address: address,
        uid:     uid,
        gid:     gid,
        mode:    mode,
    }
}

func (h InheritableNetListener) String() string {
    return fmt.Sprintf("<InheritableNetListener (%s) %q>", h.network, h.address)
}

func (h *InheritableNetListener) Open() (*os.File, error) {
    nl, err := net.Listen(h.network, h.address)
    if err != nil { return nil, err }

    h.handle = nl

    if strings.HasPrefix(h.network, "tcp") {
        if fl, err := nl.(*net.TCPListener).File(); err != nil {
            nl.Close()
            return nil, fmt.Errorf("error obtaining TCPListener File handle: %+v", err)
        } else {
            return fl, nil
        }
    } else if strings.HasPrefix(h.network, "unix") {
        rootSetMode(h.address, h.uid, h.gid, h.mode)

        if fl, err := nl.(*net.UnixListener).File(); err != nil {
            nl.Close()
            return nil, fmt.Errorf("error obtaining UnixListener File handle: %+v", err)
        } else {
            return fl, nil
        }
    }

    nl.Close()
    return nil, fmt.Errorf("unsupported network %v", h.network)
}

func (h *InheritableNetListener) Inherit(f *os.File) error {
    fl, err := net.FileListener(f)
    if err != nil { return err }
    h.handle = fl
    return nil
}

func (h InheritableNetListener) Handle() net.Listener {
    return h.handle
}

func (h *InheritableNetListener) Close() error {
    var errs []string
    if h.fileHandle != nil {
        err := h.fileHandle.Close()
        if err != nil {
            errs = append(errs, fmt.Sprintf("error closing file handle: %v", err))
        }
    }
    if h.handle != nil {
        err := h.handle.Close()
        if err != nil {
            errs = append(errs, fmt.Sprintf("error closing listener: %v", err))
        }
    }

    if errs == nil { return nil }
    return errors.New(strings.Join(errs, "; "))
}
