package queue

import (
    "database/sql"
    "fmt"
    
    "tawesoft.co.uk/go/sqlp/sqlite3"
)

type queueSqliteService struct {
    db *sql.DB
    uuidSvc UUIDService
}

// NewQueueSqliteService creates a new QueueService implemented by a SQLite
// backend that persists queues to individual database files.
//
// The SQLite backend may place a limit on the number of attached queue
// databases per connection (default 7).
//
// SQLite is used in SecureDelete mode so that deleted items are overwritten
// by zeros on disk to protect possibly sensitive data.
//
// A queue databases is VACUUMed when first attached by OpenQueue
func NewQueueSqliteService(uuidSvc UUIDService) (QueueService, error) {
    
    db, err := sqlite3.Open("sqlite3", ":memory:", sqlite3.Config{
        ForeignKeys:  true,
        SecureDelete: true,
        JournalMode:  sqlite3.JournalModeWAL,
    })
    if err != nil {
        return nil, fmt.Errorf("error opening main database: %+v", err)
    }
    
    return queueSqliteService{
        db:      db,
        uuidSvc: uuidSvc,
    }, nil
}

func (s queueSqliteService) OpenQueue(name string, path string) (Queue, error) {
    return newItemSqliteService(s.db,s.uuidSvc,  name, path)
}

func (s queueSqliteService) Close() error {
    return s.db.Close()
}

