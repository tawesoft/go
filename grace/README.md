# grace - start and gracefully shutdown processes

```shell script
go get -u "tawesoft.co.uk/go"
```

```go
import "tawesoft.co.uk/go/grace"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_grace] ∙ [docs][docs_grace] ∙ [src][src_grace] | [MIT][copy_grace] | candidate |

[home_grace]: https://tawesoft.co.uk/go/grace
[src_grace]:  https://github.com/tawesoft/go/tree/master/grace
[docs_grace]: https://www.tawesoft.co.uk/go/doc/grace
[copy_grace]: https://github.com/tawesoft/go/tree/master/grace/LICENSE.txt

## About

Package grace implements a simple way to start multiple long-lived processes
(e.g. goroutines) with cancellation, signal handling and graceful shutdown.


## Examples


Start HTTP servers on multiple ports with graceful shutdown
```go
package main

import (
    "context"
    "fmt"
    "net"
    "net/http"
    "os"
    "syscall"
    "time"

    "tawesoft.co.uk/go/grace"
)

// Implement grace.Process
type HttpServer struct {
    Addr string // e.g. ":8080"
    Srv *http.Server
    Ln net.Listener
}

func (s *HttpServer) Start(done func()) error {
    // non-async startup
    s.Srv = &http.Server{}
    ln, err := net.Listen("tcp", s.Addr)
    if err != nil { return err }
    s.Ln = ln

    // async server
    go func() {
        defer done()

        fmt.Printf("Serve %s\n", s.Addr)

        // blocks until Shutdown
        err := s.Srv.Serve(s.Ln)

        time.Sleep(50 * time.Millisecond)

        // ErrServerClosed on graceful close
        if err == http.ErrServerClosed {
            fmt.Printf("server closed normally\n")
        } else {
            fmt.Printf("server error: %v\n", err)
        }
    }()

    return nil
}

func (s *HttpServer) Stop(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, time.Second * 1)
    defer cancel()
    fmt.Printf("Shutdown %s\n", s.Addr)
    return s.Srv.Shutdown(ctx) // closes s.Ln
}

func main() {
    servers := []grace.Process{
        &HttpServer{Addr: ":8081"},
        &HttpServer{Addr: ":8082"},
        &HttpServer{Addr: ":8083"},
        &HttpServer{Addr: ":8084"},
        // &HttpServer{Addr: ":8084"}, // uncomment this to try out startup errors
        &HttpServer{Addr: ":8085"},
    }

    signals := []os.Signal{
        syscall.SIGINT,
        syscall.SIGKILL,
        syscall.SIGTERM,
    }

    // blocks until cancelled, interrupted, terminated, or killed, until
    // all servers have shutdown, and all start functions have marked
    // themselves as completely `done()` (e.g. so they have time to clean up)
    signal, errs, err := grace.Run(context.Background(), servers, signals)

    if err != nil {
        fmt.Printf("Startup error: %v\n", err)
    } else {
        fmt.Printf("Shutdown complete. Recieved %s signal\n", signal)
    }

    if len(errs) > 0 {
        fmt.Printf("Shutdown errors: %v\n", errs)
    }
}
```

## Getting Help

This package is part of [tawesoft.co.uk/go](https://www.tawesoft.co.uk/go),
a monorepo for small Go modules maintained by Tawesoft®.
Check out that URL for more information about other Go modules from
Tawesoft plus community and commercial support options.