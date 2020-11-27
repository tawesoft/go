// tawesoft.co.uk/go/drop
// 
// Copyright © 2020 Tawesoft Ltd <open-source@tawesoft.co.uk>
// Copyright © 2020 Ben Golightly <ben@tawesoft.co.uk>
// 
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction,  including without limitation the rights
// to use,  copy, modify,  merge,  publish, distribute, sublicense,  and/or sell
// copies  of  the  Software,  and  to  permit persons  to whom  the Software is
// furnished to do so, subject to the following conditions:
// 
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
// 
// THE SOFTWARE IS PROVIDED  "AS IS",  WITHOUT WARRANTY OF ANY KIND,  EXPRESS OR
// IMPLIED,  INCLUDING  BUT  NOT LIMITED TO THE WARRANTIES  OF  MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE  AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
// AUTHORS  OR COPYRIGHT HOLDERS  BE LIABLE  FOR ANY  CLAIM,  DAMAGES  OR  OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
// 

// Package drop implements the ability to start a process as root, open
// privileged resources as files, drop privileges to become a given user account,
// and inherit file handles across the dropping of privileges.
// 
// Examples
// 
// Opens privileged files and ports as root, then drops privileges
//
//     package main
//     
//     import (
//         "fmt"
//         "net"
//         "os"
//         "time"
//         
//         "tawesoft.co.uk/go/drop"
//     )
//     
//     // Define structures and methods that meets the start.Inheritable interface
//     // by implementing Name(), Open(), and Inherit()...
//     
//     // InheritableFile is a file handle that survives a dropping of
//     // privileges.
//     type InheritableFile struct {
//         Path  string
//         Flags int // e.g.  os.O_RDWR|os.O_CREATE
//         Perm  os.FileMode // e.g. 0600
//         
//         handle *os.File
//     }
//     
//     func (h InheritableFile) Name() string {
//         return h.Path
//     }
//     
//     func (h InheritableFile) Open() (*os.File, error) {
//         return os.OpenFile("/tmp/privileged-file-example", os.O_RDWR|os.O_CREATE, 0600)
//     }
//     
//     func (h *InheritableFile) Inherit(f *os.File) error {
//         f.Seek(0, 0)
//         h.handle = f
//         return nil
//     }
//     
//     // InheritableNetListener is a net.Listener that survives a dropping of
//     // privileges.
//     type InheritableNetListener struct {
//         Network string
//         Address string
//         handle net.Listener
//     }
//     
//     func (h InheritableNetListener) Name() string {
//         return fmt.Sprintf("(%s) %s", h.Network, h.Address)
//     }
//     
//     func (h InheritableNetListener) Open() (*os.File, error) {
//         nl, err := net.Listen(h.Network, h.Address)
//         if err != nil { return nil, err }
//         defer nl.Close()
//         
//         fl, err := nl.(*net.TCPListener).File()
//         if err != nil { return nil, err }
//         return fl, nil
//     }
//     
//     func (h *InheritableNetListener) Inherit(f *os.File) error {
//         defer f.Close()
//         fl, err := net.FileListener(f)
//         if err != nil { return err }
//         h.handle = fl
//         return nil
//     }
//     
//     func main() {
//         if len(os.Args) < 2 {
//             panic(fmt.Sprintf("USAGE: sudo %s username\n", os.Args[0]))
//         }
//         
//         // what user to drop to
//         username := os.Args[1]
//         
//         // resources to be opened as root and persist after privileges are dropped
//         privilegedFile := &InheritableFile{"/tmp/privileged-file-example", os.O_RDWR|os.O_CREATE, 0600, nil}
//         privilegedPort := &InheritableNetListener{"tcp4", "127.0.0.1:81", nil}
//         
//         // If the program is run as root, open privileged resources as root, then
//         // start a child process as `username` that inherits these resources and
//         // the parent process's stdio, and immediately exit.
//         //
//         // If the program is run as non-root, inherit these resources and continue.
//         err := drop.Drop(username, privilegedFile, privilegedPort)
//         if err != nil {
//             panic(fmt.Sprintf("error dropping privileges (try running as root): %v", err))
//         }
//         
//         // At this point, the program is no longer running as root, but it still
//         // has access to these privileged resources.
//         
//         // do things with privilegedFile
//         privilegedFile.handle.WriteString("hello world\n")
//         privilegedFile.handle.Close()
//         
//         // do things with privilegedPort
//         privilegedPort.handle.Close()
//         
//         time.Sleep(10)
//     }
//
//
// Package Information
//
// License: MIT (see LICENSE.txt)
//
// Stable: candidate
//
// For more information, documentation, source code, examples, support, links,
// etc. please see https://www.tawesoft.co.uk/go and 
// https://www.tawesoft.co.uk/go/drop
//
//     2020-11-27
//     
//         * Drop() functionality has been moved to tawesoft.co.uk/go/drop with
//           changes to Inheritables from a struct to an interface.
//     
package drop // import "tawesoft.co.uk/go/drop"

// Code generated by internal. DO NOT EDIT.
// Instead, edit DESC.txt and run mkdocs.sh.