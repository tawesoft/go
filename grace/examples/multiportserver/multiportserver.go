// Start HTTP servers on multiple ports with graceful shutdown
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

    // blocks until cancelled, interrupted, terminated, or killed
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
