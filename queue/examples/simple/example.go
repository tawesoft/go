// Simple example of creating queues, adding items, peeking at due items.
package main

import (
    "fmt"
    "os"
    "time"
    
    "tawesoft.co.uk/go/queue"
)

func Must(err error) {
    if err != nil { panic(err) }
}

func MustQueueService(s queue.QueueService, err error) queue.QueueService {
    if err != nil { panic(err) }
    return s
}

func MustQueue(q queue.Queue, err error) queue.Queue {
    if err != nil { panic(err) }
    return q
}

func main() {
    // Delete any existing files to start from scratch just for the sake of the
    // demonstration.
    os.Remove("q1.db")
    os.Remove("q2.db")
    
    // Give each item in the queue a unique ID e.g. if we're sending things to
    // a remote API and want to prevent duplicates.
    uuidService := queue.NewUUID4Service()
    
    // Alternatively, if you don't care about this:
    // uuidService := queue.NewUUIDNuullService()
    
    // Create a service that can create queues. We'll use the Sqlite
    // implementation with each queue backed by its own file as an attached
    // database.
    queueService := MustQueueService(queue.NewQueueSqliteService(uuidService))
    defer queueService.Close()
    
    // Create two seperate queues backed by different files
    queue1 := MustQueue(queueService.OpenQueue("q1", "q1.db"))
    defer queue1.Close()
    
    queue2 := MustQueue(queueService.OpenQueue("q2", "q2.db"))
    defer queue2.Close()
    
    // SQLite queues don't have to be persisted to disk and can also be
    // in-memory only
    queue3 := MustQueue(queueService.OpenQueue("q3", ":memory:"))
    defer queue3.Close()
    
    // Place some items in the queues due at different times in the future
    Must(queue1.CreateItem(queue.NewItem{
        Message:    "Hello world! (Queue 1; in 5 seconds)",
        Created:    time.Now().UTC(),
        RetryAfter: time.Now().UTC().Add(time.Second * 5),
    }))
    
    Must(queue1.CreateItem(queue.NewItem{
        Message:    "Hello world! (Queue 1; in 60 seconds)",
        Created:    time.Now().UTC(),
        RetryAfter: time.Now().UTC().Add(time.Second * 60),
    }))
    
    Must(queue2.CreateItem(queue.NewItem{
        Message:    "Hello world! (Queue 2)",
        Created:    time.Now().UTC(),
        RetryAfter: time.Now().UTC().Add(time.Second * 5),
    }))
    
    Must(queue3.CreateItem(queue.NewItem{
        Message:    "Hello world! (Queue 3)",
        Created:    time.Now().UTC(),
        RetryAfter: time.Now().UTC().Add(time.Second * 5),
    }))
    
    // Look up some items in the queue that are due in the future.
    // At 15 seconds in, the first item in Queue 1 should be due, but the
    // second item is not yet due.
    future := time.Now().UTC().Add(time.Second * 15)
    items, err := queue1.PeekItems(5, future)
    Must(err)
    for _, item := range items {
        fmt.Printf("got item: %+v\n", item)
    }
}
