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

NOTE: This package has only been tested on Linux. YMMV.

NOTE: This package WILL NOT WORK on Windows.

WARNING: if a process opens a config file as root, that file must be writable
by root or system accounts only. The safest way to do this is change it to
be root-owned with permissions 0644 (or 0600).


## Examples


Opens privileged files and ports as root, then drops privileges
```go
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
```

## Changes

### 2021-07-09

* The Inheritable interface has changed. It now has a Close() method. The
Name() method has also been renamed String() to satisfy the stringer
interface.

* The Drop() function now returns an extra value before the error value.
This `closer` can be used by the child process to close all Inheritable
handles. Alternatively, it is possible to ignore this and close each
handle by calling their Close() method.

* The package now exports the builtins InheritableFile and
InheritableNetListener that implement the Inheritable interface for
Files and net.Listeners. These are created by the functions
NewInheritableFile, NewInheritableTCPListener and
NewInheritableUnixListener.

* Drop() no longer panics on non-Linux platforms. However, it has only been
tested on Linux so YMMV. It will continue to panic on Windows. Listeners
also cannot be inherited on the JS platform target as they are not backed
by files.

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