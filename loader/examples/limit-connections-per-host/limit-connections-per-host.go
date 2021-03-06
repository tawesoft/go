// Configure the Loader with a Strategy to limit concurrent connections per host
package main

import (
    "fmt"
    "math/rand"
    "net/url"
    "runtime"
    "strings"
    "time"

    "tawesoft.co.uk/go/loader"
)

// interactive, if true, means we display a progress in real time.
// If false, we block until everything has loaded
const interactive = true

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
        fmt.Printf("Temporarily delaying connection to %s due to too many connections to host\n", name)
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
    rand.Seed(time.Now().UnixNano())
}

func main() {
    ldr := loader.New()

    // Initialise a NetStrategy as a method of limiting concurrent connections
    // to a single host. For our example, 2 concurrent connections.
    netStrategy := &NetStrategy{
        MaxConcurrentConnectionsPerHost: 2, // (Chrome uses 6, Firefox uses 8)
        Hosts: make(map[string]int),
    }

    // Define consumerNet on the loader as a class of worker for network files.
    // Allows up to 5 simultaneous downloads (Firefox uses 256!) but the
    // strategy will limit concurrent connections to a single host.
    consumerNet := ldr.NewConsumer(5, netStrategy)

    // Define consumerCPU on the loader as a class of worker for CPU-bound
    // tasks.
    consumerCPU := ldr.NewConsumer(runtime.NumCPU(), nil)

    // A helper function that returns a loader.Task for downloading a file
    // concurrently with consumerNet
    loadNet := func(path string) loader.Task {
        u, err := url.Parse(path)
        if err != nil { panic(err) }
        hostname := u.Hostname()

        return loader.Task{
            // Info is used by the consumer's netStrategy
            Info: func() interface{} {
                return hostname
            },
            Consumer: consumerNet,

            Load: func(_ ... interface{}) (interface{}, error) {
                // pretend to read a file
                time.Sleep(time.Millisecond * time.Duration(rand.Intn(400)))
                return fmt.Sprintf("I am network file %s!", path), nil
            },
        }
    }

    // A helper function that returns a loader.Task that does something with
    // its subtasks
    loadService := func(name string, tasks ... loader.Task) loader.Task {
        return loader.Task{
            Name: name,
            Keep: true,
            RequiresDirect: tasks,
            Consumer: consumerCPU,

            Load: func(inputs ... interface{}) (interface{}, error) {
                inputStrings := make([]string, 0)
                for _, input := range inputs {
                    inputStrings = append(inputStrings, input.(string))
                }
                result := fmt.Sprintf("I'm task %s and I have the following inputs: %s",
                    name, strings.Join(inputStrings, ", "))
                return result, nil
            },
        }
    }

    tasks := []loader.Task{
        // load files off the internet, with limits on multiple connections to
        // a single host, and do something with the results
        loadService("example.net API",
            loadNet("https://www.example.net/products.json"),
            loadNet("https://www.example.net/servers.json"),
            loadNet("https://www.example.net/news.json"),
        ),
        loadService("anotherhost API",
            loadNet("https://anotherhost.example.org/friends.json"),
            loadNet("https://anotherhost.example.org/recommendations.json"),
            loadNet("https://anotherhost.example.org/notifications.json"),
        ),
    }

    ldr.Add(tasks)

    // We can either load incrementally with a realtime progress bar
    if interactive {
        lastComplete := -1
        for {
            progress, err := ldr.Load(50 * time.Millisecond)
            if err != nil {
                fmt.Printf("Load error: %s\n", err)
                break
            }

            if progress.Completed != lastComplete {
                lastComplete = progress.Completed
                fmt.Printf("Progress: %d/%d\n", progress.Completed, progress.Total)
            }

            if progress.Done { break }

            time.Sleep(16 * time.Millisecond)
        }

    // Or just block until everything has finished loading
    } else {
        ldr.LoadAll()
    }

    // Get results
    fmt.Println(ldr.MustResult("example.net API"))
    fmt.Println(ldr.MustResult("anotherhost API"))
}
