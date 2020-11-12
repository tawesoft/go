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
    
    // Priority orders items in the queue so that due items with a higher
    // priority come before other due items, even if they were due sooner.
    Priority int // default 0, limit +/- math.MaxInt16
    
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

func (i Item) String() string {
    return i.Message
}

type NewItem struct {
    Message    string
    Priority   int
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
    // CreateItem places a new Item in the queue
    CreateItem(item NewItem) error
    
    // PeekItems returns up to `n` items with a priority >= `minPriority`,
    // and with a due time >= `due`. Returned items are ordered by (highest
    // priority, earliest due, earliest created, lowest ID). Items with IDs
    // in `excluding` (which may be nil or empty) are not included.
    PeekItems(n int, minPriority int, due time.Time, excluding []ItemID) ([]Item, error)
    
    // RetryItem reorders an item in the queue at a later `due` time and a
    // given priority.
    RetryItem(id ItemID, priority int, due time.Time) error
    
    // DeleteItem removes an item from the queue.
    DeleteItem(ItemID) error
    
    // Close any resources such as database handles.
    Close() error
    
    // Delete removes a queue database from disk. If you've opened an in-memory
    // database then don't try deleting it!
    Delete() error
}
