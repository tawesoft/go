package loader

import (
    "fmt"
    "runtime"
    "time"
)

type defaultStrategyT struct {}
var defaultStrategy *defaultStrategyT
func (s *defaultStrategyT) Start(_ interface{}) bool { return true }
func (s *defaultStrategyT) End(_ interface{}) {}

// Progress holds the result of a Loader's Load or LoadAll methods and
// represents progress in completing all Task items.
type Progress struct {
    Completed int
    Remaining int
    Done bool
}

type workResult struct {
    value interface{}
    err error
}

// Loader is used to manage a graph of Task items to be completed synchronously
// or concurrently with different types of work divided among a set of
// Consumer objects.
type Loader struct {
    namedResults map[string]workResult // todo index instead
    progress Progress
    consumers []consumer // all consumers
    seqConsumers []*seqConsumer // subset
    dag         dag
    resultsChan chan consumerWorkResult
}

// Returns a new Loader
func New() *Loader {

    // Maximum number of tasks that can complete without blocking between
    // successive calls to Load. Assume as low as 16 Hz between calls (~64ms),
    // and all tasks taking 0.25ms, assume ideal fully linear multithreaded
    // speedup by adding logical CPUs, and assume no other bottlenecks.
    resultBufferSize := 64 * 4 * runtime.NumCPU();

    l := &Loader{
        namedResults: make(map[string]workResult),
        resultsChan:  make(chan consumerWorkResult, resultBufferSize),
        //results: make(chan consumerWorkResult), // unbuffered
    }

    // Default sequential consumer with zero-value ConsumerID
    l.NewConsumer(0, nil)

    return l
}

// Returns a named result from a Task where Keep is true and the Name is unique
// across all Tasks where Keep is True.
func (l *Loader) Result(name string) (workResult, bool) {
    value, exists := l.namedResults[name]
    return value, exists
}

func (l *Loader) MustResult(name string) interface{} {
    return l.namedResults[name].value
}

// NewConsumer creates a task consumer that performs the (possibly
// asynchronous) completion of tasks at a given level of concurrency (e.g.
// number of goroutines) and returns an opaque ID that uniquely identifies that
// consumer with the active Loader.
//
// A concurrency of zero means that the consumer's tasks will be performed
// sequentially on the same thread as the Loader's Load() or LoadAll() methods.
//
// The strategy argument allows control over temporarily delaying a task.
// Strategy may be nil to always accept.
//
// The special ConsumerID of zero corresponds to a default builtin consumer
// that has a concurrency of zero and a nil strategy.
func (l *Loader) NewConsumer(concurrency int, strategy Strategy) ConsumerID {
    if strategy == nil {
        strategy = defaultStrategy
    }

    var c consumer
    id := ConsumerID(len(l.consumers))

    if concurrency > 0 {
        mtc := &mtConsumer{
            id:          id,
            concurrency: concurrency,
            strategy:    strategy,
            pending:     make([]int, 0),
            current:     -1,

            jobs:       make(chan int), // maybe buffer of size c.concurrency
        }

        go mtc.manager(&l.dag, l.resultsChan)
        c = mtc
    } else {
        seqc := &seqConsumer{
            id: id,
            strategy:    strategy,
            pending:     make([]int, 0),
            current:     -1,
        }
        c = seqc
        l.seqConsumers = append(l.seqConsumers, seqc)
    }

    l.consumers = append(l.consumers, c)
    return id
}

func (l *Loader) sendPendingTasksToConsumers() []consumerWorkResult {
    var results []consumerWorkResult

    for _, idx := range l.dag.pending {
        task := l.dag.nodes[idx]
        consumer := l.consumers[task.Consumer]
        rs := consumer.sendTask(idx, l.resultsChan)
        if rs != nil {
            results = append(results, rs...)
        }
    }

    // empty pending tasks
    l.dag.pending = l.dag.pending[:0]

    return results
}

// TODO callers need to check this for error
func (l *Loader) registerResult(idx int, result interface{}, err error) error {
    l.progress.Remaining--
    l.progress.Completed++

    if err == nil {
        l.dag.registerResult(idx, result)
    } else {
        // if err is not fatal, this can be a nil input for a parent task
        l.dag.registerResult(idx, nil)
    }

    // stash the result if its a kept task result with a unique name
    // TODO store the index instead
    task := l.dag.nodes[idx]
    if task.Keep {
        name := task.Name
        _, alreadyExists := l.namedResults[name]
        if alreadyExists {
            return fmt.Errorf("kept named result %q already exists (but must be unique)", name)
        }
        l.namedResults[name] = workResult{
            value: result,
            err:   err,
        }
    }

    return nil
}

// Load completes as many loading tasks as possible within the time budget. If
// idle while waiting for concurrent results, it may return early.
//
// See also the LoadAll() method.
func (l *Loader) Load(budget time.Duration) (Progress, error) {

    earlyResults := l.sendPendingTasksToConsumers()

    start := time.Now()

    outer:
    for {
        // TODO move pending tasks here

        // while there is sequential work to be done...
        for _, c := range l.seqConsumers {
            c.manage(&l.dag, &l.dag.nodes, l.registerResult)
        }

        if earlyResults != nil {
            for _, result := range earlyResults {
                l.registerResult(result.idx, result.result, result.err)
            }
            earlyResults = nil
        }

        // pick up any completed jobs from any consumer, without blocking
        select {
            case result := <- l.resultsChan:
                l.registerResult(result.idx, result.result, result.err)
            default:
                // no progress possible yet, so break early
                break outer
        }

        // continue unless budget exceeded...
        elapsed := time.Now().Sub(start)
        if elapsed >= budget { break }
    }

    l.progress.Done = (l.progress.Remaining == 0)

    return l.progress, nil
}

// LoadAll completes all loading tasks and blocks until finished
func (l *Loader) LoadAll() (Progress, error) {

    earlyResults := l.sendPendingTasksToConsumers()

    // TODO
    earlyResults = earlyResults

    for l.progress.Remaining > 0 {
        // while there is sequential work to be done...
        for _, c := range l.seqConsumers {
            c.manage(&l.dag, &l.dag.nodes, l.registerResult)
        }

        // if there are still concurrent tasks left...
        if l.progress.Remaining == 0 { break }

        // then any further progress is blocked on a concurrent task
        select {
            case result := <- l.resultsChan:
                l.registerResult(result.idx, result.result, result.err)
        }
    }

    l.progress.Done = true
    return l.progress, nil
}

// Close TODO
func (l *Loader) Close() {
    // TODO
}

// identify returns a fully qualified string representation of a task's
// position in the hierarchy of tasks e.g. for locating errors.
//
// For example, "parentTask1.parentTask2.<anonymous-task>.taskName"
//
// The string should be considered opaque but human-readable. Do not rely on
// the format being fixed or parsable.
func (l *Loader) identify(idx int) {

}

// Add adds a graph of Tasks to the Loader to later process with its Load() or
// LoadAll() methods. Tasks may be added even during or after a Load() loop.
//
// Because Task names, even at the top level of the array, are scoped to this
// function, Two Tasks across a Loader Add() boundary cannot refer to each
// other by name. If this behaviour is desired, append to a Task array and send
// the combined array in one call to Add.
//
// An error is generated if a named Task requirement is not in scope. In this
// event, the state of the task dependency graph is undefined and no methods
// on the Loader, other than Close, may be called.
func (l *Loader) Add(tasks []Task) error {
    err := l.dag.add(tasks)
    l.progress.Remaining += len(l.dag.nodes)
    l.progress.Done = false
    return err
}

// Task is a (possibly recursively nested) unit of work that is to be performed
// (possibly concurrently), subject to dependencies and limits.
type Task struct {
    // Optional name of the task. Used to reference a task as a dependency of
    // another and to retrieve its result, if kept. Does not have to be unique
    // (it is scoped to its subtasks and successor siblings in one call to a
    // Loader Add method), unless it is kept (see the Keep field, below).
    Name string

    // Keep indicates that the task's Load() result will be available from
    // the Loader Result() method by the task Name. If Keep is true, the
    // task's Name must be globally unique across all Tasks kept by the loader.
    // If Keep is false, the task's name does not have to be unique, even if a
    // kept Task has the same name.
    Keep bool

    // Load performs the (possibly asynchronous) completion of the task e.g.
    // reading a file from disk, a unit of computation, etc.
    //
    // The results argument is the ordered results of the tasks in
    // RequiresNamed (if any) followed by the ordered results of the tasks in
    // RequiresDirect (if any).
    Load func(results ... interface{}) (interface{}, error)

    // Free performs the (possibly asynchronous) removal of a task's Load()
    // result e.g. releasing memory. May be nil.
    Free func(i interface{})

    // RequiresNamed is an array of names of tasks required to complete first
    // as a dependency of this task. May be nil.
    //
    // Note that a required named task must be defined before a task can depend
    // on it (e.g. by appearing earlier in the array passed to Loader Add()).
    RequiresNamed []string

    // RequiresDirect is an array of tasks required to complete first as a
    // dependency of this task. May be nil. These "subtasks" are in a new scope
    // for naming purposes.
    RequiresDirect []Task

    // Consumer performs the asynchronous completion of tasks at a given level
    // of concurrency. Use an ID returned by a Loader Consumer() method. May be
    // zero, in which case the task is completed in the same thread as the
    // caller Loader's Load() or LoadAll() methods.
    Consumer ConsumerID

    // Info returns some value for a task understood by the given Strategy.
    // May be nil. Must be a constant function i.e. always return the same
    // value for a given Task.
    Info func() interface{}
}

// info() calls a Task's Info() method, or if not defined acts as if the
// Info() method does nothing except return nil.
func (t *Task) info() interface{} {
    if t.Info != nil { return t.Info() }
    return nil
}

// NamedTask is used to reference another task by name as a subtask i.e. in
// a Task's RequiresDirect instead of RequiresNamed.
//
// TODO NamedTasks should be simplified in the DAG to remove the node entirely.
func NamedTask(name string) Task {
    return Task{
        Name: fmt.Sprintf("<loader NamedTask(%q) passthrough>", name),
        RequiresNamed: []string{name},
        Load: func(results ... interface{}) (interface{}, error) {
            if len(results) < 1 { return "<NamedTask TODO>", nil } // panic("oops") }
            return results[0], nil
        },
        Consumer: 0, // special maybe -1 for high priority
    }
}

