package queue

import (
    "time"
)

type UUIDService interface {
    Generate() ([]byte, error)
}

type ItemID int64

type Item struct {
    // ID uniquely identifies the item in a single queue
    ID ItemID
    
    // UUID uniquely identifies the item across multiple systems and can be
    // used e.g. as a cryptographic nonce to prevent replays / duplicates
    UUID []byte
    
    // Message is any user-supplied string of bytes
    Message string
    
    // Attempt records how many times the given item has been unsuccessfully
    // processed and put back in the queue
    Attempt int
    
    // Created is the time the item was first added to the queue
    // (use `time.Now().UTC()`)
    Created time.Time // UTC
    
    // RetryAfter is the earliest time the queue will attempt to process the
    // item
    RetryAfter time.Time // UTC
}

type NewItem struct {
    Message    string
    Created    time.Time
    RetryAfter time.Time
}

// QueueService defines an interface for the creation of queues.
//
// One such implementation is the NewQueueSqliteService which provides a
// reliable persistent queues backed by SQLite databases.
type QueueService interface{
    // OpenQueue opens (or creates) a new queue with a given name, backed
    // by a file at the given path (the SQLite backend also supports
    // ":memory:" as a target path)
    OpenQueue(name string, path string) (Queue, error)
    
    // Close any resources such as database handles. You should individually
    // close any open queues first.
    Close() error
}

type Queue interface {
    CreateItem(item NewItem) error
    PeekItems(n int, due time.Time) ([]Item, error)
    //RetryItem(ItemID, time.Time) error
    //DeleteItem(ItemID) error
    
    Close() error
    Delete() error
}
