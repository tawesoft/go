// Opens privileged files and ports as root, then drops privileges
package main

import (
    "fmt"
    "io/ioutil"
    "net"
    "os"
    "time"
    
    "tawesoft.co.uk/go/start"
)

func main() {
    if len(os.Args) < 2 {
        panic(fmt.Sprintf("USAGE: sudo %s username\n", os.Args[0]))
    }
    
    var (
        ln net.Listener
        rootfile *os.File
        username = os.Args[1]
    )
    
    // If the program is run as root, open privileged resources as root, then
    // start a child process as `username` that inherits these resources and
    // the parent process's stdio, and immediately exit.
    //
    // If the program is run as non-root, inherit these resources and continue.
    err := start.Drop(username, []start.Inheritable{
        {
            Name: "privileged-port",
            Open: func() (*os.File, error) {
                nl, err := net.Listen("tcp4", "127.0.0.1:81")
                defer nl.Close()
                if err != nil { return nil, err }
                fl, err := nl.(*net.TCPListener).File()
                if err != nil { return nil, err }
                return fl, nil
            },
            
            Set: func(f *os.File) error {
                defer f.Close()
                fl, err := net.FileListener(f)
                if err != nil { return err }
                ln = fl
                return nil
            },
        },
        
        {
            Name: "privileged-file",
            Open: func() (*os.File, error) {
                f, err := os.OpenFile("/tmp/privileged-file-example", os.O_RDWR|os.O_CREATE, 0600)
                if err != nil { return nil, err }
                f.Write([]byte("this file is only readable by root!\n"))
                return f, err
            },
        
            Set: func(f *os.File) error {
                rootfile = f
                rootfile.Seek(0, 0)
                return nil
            },
        },
    })
    if err != nil { panic(err) }
    
    
    // At this point, the program is no longer running as root, but it still
    // has access to these privileged resources.
    
    defer rootfile.Close()
    contents, err := ioutil.ReadAll(rootfile)
    if err != nil { panic(err) }
    fmt.Printf("read file result: %s", string(contents))
    
    time.Sleep(10)
}
