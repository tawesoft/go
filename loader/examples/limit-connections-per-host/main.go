// Configure the loader to limit concurrent connections per host
package main

import (
    "fmt"
    "net/url"
    "runtime"
    "time"

    "tawesoft.co.uk/go/loader"
)

// NetStrategy is a loader.Strategy for limiting the number of concurrent
// connections to a single host
type NetStrategy struct {
    // Limit concurrent connections to a single host
    // Firefox uses 8, Chrome uses 6.
    MaxConcurrentConnectionsPerHost int

    // A count of concurrent connections by hostname
    Hosts map[string]int
}

// Start returns true if the task may start. Checks the current connections by
// hostname to see if it exceeds the limit or not.
func (s *NetStrategy) Start(info interface{}) bool {
    name := info.(string)
    count, _ := s.Hosts[name]

    if count >= s.MaxConcurrentConnectionsPerHost {
        fmt.Printf("Temporarily delaying task due to too many connections to host %s\n", name)
        return false
    }

    s.Hosts[name] = count + 1
    return true
}

// End indicates the task has completed, so we no longer have to count it
// towards the limit.
func (s *NetStrategy) End(info interface{}) {
    name := info.(string)
    count := s.Hosts[name]
    s.Hosts[name] = count - 1
}

func init() {
    // for main thread only code (like OpenGL)
    runtime.LockOSThread()
}

func main() {
    netStrategy := &NetStrategy{
        MaxConcurrentConnectionsPerHost: 1, // e.g. 8
        Hosts: make(map[string]int),
    }

    ldr := loader.New()

    consumerDisk := ldr.NewConsumer(1, nil) // try 4 for SSD
    consumerCPU  := ldr.NewConsumer(runtime.NumCPU(), nil)
    consumerNet  := ldr.NewConsumer(16, netStrategy) // Firefox uses 256!

    // a loader Task for reading a file from disk
    loadFile := func(keep bool, name string, path string) loader.Task {
        return loader.Task{
            Name: name,
            Keep: keep,
            Consumer: consumerDisk,
            Load: func(_ ... interface{}) (interface{}, error) {
                // pretend to read a file
                //time.Sleep(time.Millisecond * 50)
                return fmt.Sprintf("I am file %s!", path), nil
            },
        }
    }

    // a loader Task for downloading a file
    loadNet := func(keep bool, name string, path string) loader.Task {
        u, err := url.Parse(path)
        if err != nil { panic(err) }
        hostname := u.Hostname()

        return loader.Task{
            Name: name,
            Keep: keep,
            Info: func() interface{} {
                return hostname
            },
            Consumer: consumerNet,
            Load: func(_ ... interface{}) (interface{}, error) {
                // pretend to read a file
                return fmt.Sprintf("I am network file %s!", path), nil
            },
        }
    }
    loadNet = loadNet // TODO

    // Concatenate the results of two named tasks
    loadConcatNamed := func(keep bool, name string, u string, v string) loader.Task {
        return loader.Task{
            Name: name,
            Keep: keep,
            RequiresNamed: []string{u, v},
            Consumer: consumerCPU,
            Load: func(results ... interface{}) (interface{}, error) {
                if len(results) < 2 { return "TODO inputs", nil }
                left := results[0].(string)
                right := results[1].(string)
                return left + " " + right, nil
            },
        }
    }

    // Concatenate the results of two subtasks
    loadConcatSubtasks := func(keep bool, name string, u loader.Task, v loader.Task) loader.Task {
        return loader.Task{
            Name: name,
            Keep: keep,
            RequiresDirect: []loader.Task{u, v},
            Consumer: consumerCPU,
            Load: func(results ... interface{}) (interface{}, error) {
                if len(results) < 2 { return "TODO inputs", nil }
                left := results[0].(string) // u
                right := results[1].(string) // v
                return left + " " + right, nil
            },
        }
    }

    tasks := []loader.Task{

        // load two files as named tasks, then concatenate them
        loadFile(false, "A",        "A.txt"),
        loadFile(true,  "B",        "B.txt"),
        loadConcatNamed(true, "AB", "A", "B"),

        // equivalently...
        loadConcatSubtasks(true, "AB2",
            loader.NamedTask("A"),
            loader.NamedTask("B"),
        ),

        // load two files as anonymous subtasks, then concatenate them
        loadConcatSubtasks(true, "CD",
            loadFile(false, "(C)",  "C.txt"), // anonymous
            loadFile(true, "D", "D.txt"), // named but locally scoped
        ),

        loadFile(false, "D", "D2.txt"), // globally scoped

        // depends on task AB and task CD
        loadConcatNamed(true, "ABCD", "AB", "CD"),

        // refer to a named task alongside an anonymous subtask
        loadConcatSubtasks(true, "ABCDEEFA",
            loader.NamedTask("ABCD"),
            loadConcatSubtasks(false, "(EEFA)",
                loadFile(false, "E", "E.txt"), // locally scoped
                loadConcatSubtasks(false, "(EFA)",
                    loader.NamedTask("E"), // in parent scope
                    loadConcatSubtasks(false, "(FA)",
                        loadFile(false, "(F)",  "F.txt"),
                        loader.NamedTask("D"), // globally scoped
                    ),
                ),
            ),
        ),

        // load files off the internet, with limits on multiple connections to
        // a single host
        loadNet(true, "products.json", "https://www.example.net/products.json"),
        loadNet(true, "servers.json",  "https://www.example.net/servers.json"),
        loadNet(true, "news.json",     "https://www.example.net/news.json"),
        loadNet(true, "news.xml",      "https://anotherhost.example.org/news.xml"),


        // load with no concurrency at all
        {
            Name: "sequential-task-1",
            Consumer: 0, // implicit
            Load: func(results ... interface{}) (interface{}, error) {
                return "I am sequential task one!", nil
            },
        },
        {
            Name: "sequential-task-2",
            Consumer: 0, // implicit
            Load: func(results ... interface{}) (interface{}, error) {
                return "I am sequential task two!", nil
            },
        },

    }

    ldr.Add(tasks)

    // progress, err := ldr.LoadAll()
    // fmt.Printf("%+v %s\n", progress, err)

    for {
        progress, err := ldr.Load(50 * time.Millisecond)
        if err != nil {
            fmt.Printf("Load error: %s\n", err)
            break
        }
        fmt.Printf("... %+v\n", progress)
        if progress.Done { break }

        time.Sleep(16 * time.Millisecond)
    }

    fmt.Printf(
`Here are some results:\n
    B:             %s
    D:             %s
    AB:            %s
    AB2:           %s
    CD:            %s
    ABCD:          %s
    ABCDEEFA:      %s
    products.json: %s
    servers.json:  %s
    news.json:     %s
    news.xml:      %s 
`,
        ldr.MustResult("B"),
        ldr.MustResult("D"),
        ldr.MustResult("AB"),
        ldr.MustResult("AB2"),
        ldr.MustResult("CD"),
        ldr.MustResult("ABCD"),
        ldr.MustResult("ABCDEEFA"),
        ldr.MustResult("products.json"),
        ldr.MustResult("servers.json"),
        ldr.MustResult("news.json"),
        ldr.MustResult("news.xml"),
    )
}
