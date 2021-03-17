# drop - drop privileges and inherit handles

```shell script
go get -u "tawesoft.co.uk/go"
```

```go
import "tawesoft.co.uk/go/drop"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_drop] ∙ [docs][docs_drop] ∙ [src][src_drop] | [MIT][copy_drop] | candidate |

[home_drop]: https://tawesoft.co.uk/go/drop
[src_drop]:  https://github.com/tawesoft/go/tree/master/drop
[docs_drop]: https://www.tawesoft.co.uk/go/doc/drop
[copy_drop]: https://github.com/tawesoft/go/tree/master/drop/LICENSE.txt

## About

Package drop implements the ability to start a process as root, open
privileged resources as files, drop privileges to become a given user account,
and inherit file handles across the dropping of privileges.


## Examples


Opens privileged files and ports as root, then drops privileges
```go
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
```

## Changes

### 2021-03-17

* Drop() now returns a (bool, error) 2-tuple. The first return value,
if true, indicates that the caller should immediately exit.

### 2020-11-27

* Drop() functionality has been moved to tawesoft.co.uk/go/drop with
changes to Inheritables from a struct to an interface.


## Getting Help

This package is part of [tawesoft.co.uk/go](https://www.tawesoft.co.uk/go),
a monorepo for small Go modules maintained by Tawesoft®.
Check out that URL for more information about other Go modules from
Tawesoft plus community and commercial support options.