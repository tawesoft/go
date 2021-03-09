package grace

import (
    "context"
    "os"
    "os/signal"
    "sync"
    "syscall"
)

// Process defines a long-lived cancellable process
type Process interface {
    // Start launches a process. This shouldn't block, so launch goroutines or
    // other programs as necessary. Call done() when that goroutine returns
    // or that process has terminated. If Start returns an error, then done()
    // must not be called.
    Start(done func()) error

    // Stop stops a process. Derive from the context as necessary e.g.
    // use context.WithTimeout
    Stop(context.Context) error
}

func shutdown(ctx context.Context, processes []Process) []error {
    var errors []error

    for _, process := range processes {
        if err := process.Stop(ctx); err != nil {
            errors = append(errors, err)
        }
    }

    return errors
}

// Run starts each process, waits for any signal in `notify` or
// until ctx is cancelled, then cancels each process. It blocks until all
// processes have been stopped with Stop() and all process Start() functions
// mark themselves as completely done().
//
// The first return value is the os.Signal received.
//
// The second return value is a list of errors for processes that returned
// errors when being cancelled.
//
// The third return value returns any startup error.
//
// The second return value may be non-empty even when the third return value is
// non-nil when there is both a startup error and an error stopping any
// previously started processes e.g. if process one starts but process two
// fails, then process one needs to be cancelled but may also run into an error
// cancelling.
func Run(ctx context.Context, processes []Process, signals []os.Signal) (os.Signal, []error, error) {

    var wg sync.WaitGroup
    closechan := make(chan struct{sig os.Signal; errs []error}, 1)

    donefn := func() {
        wg.Done()
    }

    // start each process
    for i, process := range processes {
        err := process.Start(donefn)
        if err != nil {
            errors := shutdown(ctx, processes[0:i])
            wg.Wait()
            return syscall.Signal(0), errors, err
        }
        wg.Add(1)
    }

    // start the signal listener
    go func(ctx context.Context) {
        // block until interrupt signal, cancel signal, or context cancelled
        var result os.Signal
        sigchan := make(chan os.Signal, 1)
        signal.Notify(sigchan, signals...)

        select {
            case result = <- sigchan:
            case <- ctx.Done():
        }

        errors := shutdown(ctx, processes)
        wg.Wait()
        closechan <- struct{sig os.Signal; errs []error}{result, errors}
    }(ctx)

    // blocks until safe to exit
    result := <- closechan
    return result.sig, result.errs, nil
}
