package loader

// Strategy allows temporarily delaying of a task based on other currently
// progressing tasks e.g. see the examples folder for an implementation that
// avoids multiple concurrent connections to a single host.
//
// Note: do not share the same instance of a Strategy across two consumers
// without making it thread safe e.g. use of mutexes.
//
// Note: the Strategy must give a constant answer for the same sequence of
// Start and End methods and arguments i.e. must not depend on an external
// factor such as time or user input but must only depend on some innate
// property of the currently accepted tasks.
//
// Note: if not used carefully, a Strategy, or the interaction between two
// Strategies, may lead to deadlocks. This may be the case if there is any
// construction of a Task dependency graph (considering only those Tasks that
// may be delayed by a Strategy), or any subgraph thereof formidable by
// removing leaf nodes, where the Strategy's lower bound on the number of
// simultaneous tasks is less than or equal to the number of leaf nodes minus
// the number of vertex disjoint subgraphs of that graph.
type Strategy interface {
    // Start takes the result of a task Info() and returns true if the task is
    // accepted, or false if the task must be temporarily delayed.
    Start(info interface{}) bool

    // End takes the result of a task Info() and registers that the task has
    // completed processing.
    End(info interface{})
}

// ConsumerID uniquely identifies a consumer in a given Loader.
//
// See the NewConsumer() method on the Loader type.
type ConsumerID int

type consumer interface {
    sendTask(idx int, resultChan chan consumerWorkResult) []consumerWorkResult
}

type consumerWorkResult struct {
    idx int // task index
    result interface{} // result of Load() function
    err error // result of Load() function
}

// mtConsumer is a concurrent implementation
type mtConsumer struct {
    id ConsumerID
    concurrency int
    strategy Strategy
    pending []int  // index of tasks to send to workers
    current int    // index of next task to immediately send to first available worker
    jobs chan int  // index of task sent to workers; < 0 to kill workers and consumer
}

func (c *mtConsumer) sendTask(idx int, resultChan chan consumerWorkResult) []consumerWorkResult {
    var results []consumerWorkResult

    // this thorny code is because a consumer might finish jobs even before
    // we're finished sending them, so we have to receive and queue at the
    // same time

    for {
        select {
            case c.jobs <- idx: // send
               // OK
               // fmt.Printf("send %d ok\n", idx)
               return results
            case result := <- resultChan: // have to receive first
                results = append(results, result)
                // fmt.Printf("early recieve %d\n", result.idx)
        }
    }
}

func mtConsumerWorker(dag *dag, consumerID ConsumerID, workerID int, jobs <-chan int, results chan<- consumerWorkResult) {
    for j := range jobs {
        if j < 0 { break } // kill signal

        // fmt.Printf("worker (%d, %d) started job %d\n", consumerID, workerID, j)

        result, err := dag.nodes[j].Load(dag.inputs(j)...)
        //dag.results[j] = result

        results <- consumerWorkResult{
            idx: j,
            result: result,
            err: err,
        }

        // fmt.Printf("worker (%d, %d) returned job %d\n", consumerID, workerID, j)
    }
}

// manager sends and collects work and results from worker routines
func (c *mtConsumer) manager(dag *dag, results chan<- consumerWorkResult) {
    workerJobs    := make(chan int) // maybe buffer of size c.concurrency
    workerResults := make(chan consumerWorkResult)
    availableWorkers := c.concurrency

    for i := 1; i <= c.concurrency; i++ {
        go mtConsumerWorker(dag, c.id, i, workerJobs, workerResults)
    }

    for {

        // block until a job or a result appears
        select {
            case idx := <- c.jobs:
                if idx < 0 {
                    // kill
                    for i := 1; i <= c.concurrency; i++ {
                        workerJobs <- 0
                    }
                    return
                }
                c.pending = append(c.pending, idx)
            case result := <- workerResults:
                // notify strategy we're done with this task
                c.strategy.End(dag.nodes[result.idx].info())

                // pass result back up to loader
                // (we have to capture it first in order to wake)
                availableWorkers++
                results <- result
        }

        // nothing immediately queued, find the first acceptable pending task
        // according to strategy that is to become immediately queued
        if c.current < 0 {
            for i := 0; i < len(c.pending); i++ {
                idx := c.pending[i]
                task := dag.nodes[idx]
                if c.strategy.Start(task.info()) {
                    c.current = idx
                    c.pending = intArrayDeleteElement(c.pending, i)
                    break
                }
            }
        }

        // attempt to send immediately queued task, if any
        if (c.current >= 0) && (availableWorkers > 0) {
            // a worker is guaranteed to be ready (although in practice may be
            // a few nanoseconds away due to loop overhead), so block waiting
            // for it rather than use a non-blocking select
            workerJobs <- c.current
            c.current = -1
            availableWorkers--
        }
    }
}

// seqConsumer is a non-concurrent / sequential implementation
type seqConsumer struct {
    id ConsumerID
    strategy Strategy
    pending []int  // index of tasks to send to workers
    current int    // index of next task to immediately send to first available worker
}

func (c *seqConsumer) sendTask(idx int, _ chan consumerWorkResult) []consumerWorkResult {
    c.pending = append(c.pending, idx)
    return nil
}

func (c *seqConsumer) manage(dag *dag, tasksByIndex *[]*Task, registerResult func(int, interface{}, error) error) {
    for len(c.pending) > 0 {
        // nothing immediately queued, find the first acceptable pending task
        // according to strategy that is to become immediately queued
        if c.current < 0 {
            for i := 0; i < len(c.pending); i++ {
                idx := c.pending[i]
                task := (*tasksByIndex)[idx]
                if c.strategy.Start(task.info()) {
                    c.current = idx
                    c.pending = intArrayDeleteElement(c.pending, i)
                    break
                }
            }
        }

        // run immediately queued task, if any
        if c.current >= 0 {
            task := (*tasksByIndex)[c.current]
            result, err := task.Load(dag.inputs(c.current)...)
            c.strategy.End((*tasksByIndex)[c.current].info())

            // dag.results[c.current] = result
            registerResult(c.current, result, err)
            c.current = -1
        }
    }
}
