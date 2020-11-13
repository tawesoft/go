package start

import (
    "context"
    "os"
    "os/signal"
    "syscall"
    "time"
)

// Process defines a long-lived async process
type Process struct {
    Name string
    Start func()
    Cancel func(context.Context) error
    CancelTimeout time.Duration
}

// Start starts each process in order, and cancels all running processes on
// interrupt. It blocks until all processes are gracefully closed (if possible).
//
// The return value is either nil, or a list of errors for processes that
// returned shutdown errors
func Start(processes []Process) []error {
    
    var (
        errors []error
        closechan = make(chan interface{}, 1)
    )
    
    // start the signal listener
    go func() {
        // block until interrupt signal
        sigchan := make(chan os.Signal, 1)
        signal.Notify(sigchan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
        <- sigchan

        // gracefully shutdown each process
        for _, process := range processes {
            ctx, cancel := context.WithTimeout(context.Background(), process.CancelTimeout)
            if err := process.Cancel(ctx); err != nil {
                errors = append(errors, err)
            }
            cancel()
        }
        
        // all idle connections now closed and socket closed
        closechan <- nil
    }()
    
    // start each process
    for _, process := range processes {
        go process.Start()
    }
    
    // blocks until safe to exit
    <- closechan
    
    return errors // maybe nil
}
