package queue

import (
    "database/sql"
    "fmt"
    "os"
    "time"
    
    "tawesoft.co.uk/go/sqlp"
    "tawesoft.co.uk/go/sqlp/sqlite3"
    "tawesoft.co.uk/go/variadic"
)

type itemSqliteService struct {
    dbname  string
    dbpath  string
    db      *sql.DB
    msgSvc  messageSqliteService
}
var nilItemSqliteService = itemSqliteService{}

type itemSqliteNew struct {
    Priority   int // default 0
    Created    int64 // time.Time // UTC
    RetryAfter int64 // time.Time // UTC
}
func (p itemSqliteNew) fields() string {
    return "priority, created, retryAfter"
}
func (p itemSqliteNew) placeholders() string {
    return "?, ?, ?"
}
func (p itemSqliteNew) values() []interface{} {
    return []interface{}{p.Priority, p.Created, p.RetryAfter}
}

type itemSqlite struct {
    ID ItemID
    Priority int
    Message string
    Attempt int
    Created int64 // time.Time // UTC
    RetryAfter int64 // time.Time // UTC
}
func (p itemSqlite) fields() string {
    return "items.id, items.priority, messages.data, items.attempt, items.created, items.retryAfter"
}
func (p itemSqlite) placeholders() string {
    return "?, ?, ?, ?, ?, ?"
}
func (p itemSqlite) values() []interface{} {
    return []interface{}{p.ID, p.Priority, p.Message,p.Attempt, p.Created, p.RetryAfter}
}
func (p *itemSqlite) pointers() []interface{} {
    return []interface{}{&p.ID, &p.Priority, &p.Message, &p.Attempt, &p.Created, &p.RetryAfter}
}
func (p itemSqlite) toItem() Item {
    return Item{
        ID:         p.ID,
        Priority:   p.Priority,
        Message:    p.Message,
        Attempt:    p.Attempt,
        Created:    time.Unix(p.Created, 0),
        RetryAfter: time.Unix(p.RetryAfter, 0),
    }
}

func initItemSqliteService(db *sql.DB, dbname string, file string) error {
    _, err := db.Exec(`
        ATTACH DATABASE `+file+` as `+ dbname +`;
        VACUUM `+ dbname +`;

        CREATE TABLE IF NOT EXISTS `+ dbname +`.items (
            id         INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
            priority   INTEGER NOT NULL DEFAULT 0,
            attempt    INTEGER NOT NULL DEFAULT 0,
            created    INTEGER NOT NULL, -- epoch seconds
            retryAfter INTEGER NOT NULL  -- epoch seconds
        );

        CREATE INDEX IF NOT EXISTS `+ dbname +`.items_idx_sort ON items(priority DESC, retryAfter, created, id);
    `)
    if err != nil {
        return fmt.Errorf("error initialising item table: %+v", err)
    }
    return nil
}

// newItemSqliteService creates a new itemService implemented by a
// itemSqliteService
func newItemSqliteService(db *sql.DB, name string, rawpath string) (itemSqliteService, error) {
    var err error
    
    name, err = sqlite3.Features.EscapeIdentifier(name)
    if err != nil { return nilItemSqliteService, err }
    path, err := sqlite3.Features.EscapeString(rawpath)
    if err != nil { return nilItemSqliteService, err }
    
    err = initItemSqliteService(db, name, path)
    if err != nil { return nilItemSqliteService, err }
    
    msgSvc, err := newMessageSqliteService(db, name)
    if err != nil { return nilItemSqliteService, err }
    
    return itemSqliteService{
        dbname:  name,
        dbpath:  rawpath,
        db:      db,
        msgSvc:  msgSvc,
    }, nil
}

func (s itemSqliteService) CreateItem(newItem NewItem) error {
    err := func() error {
        tx, err := s.db.Begin()
        if err != nil { return err }
        defer tx.Rollback()
        
        i := itemSqliteNew{
            Priority:   newItem.Priority,
            Created:    newItem.Created.Unix(),
            RetryAfter: newItem.RetryAfter.Unix(),
        }
        
        query := `INSERT INTO `+ s.dbname +`.items(`+i.fields()+`) VALUES (`+i.placeholders()+`)`
        
        result, err := tx.Exec(query, i.values()...)
        if err != nil { return err }
        
        if n, ok := sqlp.RowsAffectedBetween(result, 1, 1); !ok {
            return fmt.Errorf("rows affected %d != 1", n)
        }
        
        itemID, err := result.LastInsertId()
        if err != nil { return err }
        
        err = s.msgSvc.Create(tx, ItemID(itemID), newItem.Message)
        if err != nil { return err }
        
        err = tx.Commit()
        if err != nil { return err }
        
        return nil
    }()
    
    if err != nil {
        return fmt.Errorf("error inserting new item: %+v", err)
    }
    
    return nil
}

func (s itemSqliteService) PeekItems(
    n int,
    minPriority int,
    due time.Time,
    excluding []ItemID,
) ([]Item, error) {
    items, err := func () ([]Item, error) {
        var i itemSqlite
        
        query := `
        SELECT
            `+ i.fields() +`
        FROM
            `+ s.dbname +`.items,
            `+ s.dbname +`.messages
        WHERE
            items.id = messages.id
            AND
                items.priority >= ?
            AND
                items.retryAfter <= ?
            AND
                items.id NOT IN (`+ sqlp.RepeatString("?", len(excluding)) +`)
        ORDER BY
            items.priority DESC,
            items.retryAfter,
            items.created,
            items.id
        LIMIT ?`
        
        args := variadic.FlattenExcludingNils(minPriority, due.Unix(), excluding, n)
        rows, err := s.db.Query(query, args...)
        if err != nil { return nil, err }
        defer rows.Close()
        
        items := make([]Item, 0, n)
        
        for rows.Next() {
            if err := rows.Scan(i.pointers()...); err != nil {
                return nil, err
            }
            items = append(items, i.toItem())
        }
        if err := rows.Err(); err != nil {
            return nil, err
        }
        
        return items, nil
    }()
    
    if err != nil {
        return nil, fmt.Errorf("error selecting items: %+v", err)
    }
    
    return items, nil
}

func (s itemSqliteService) RetryItem(id ItemID, priority int, due time.Time) error {
    query := `
        UPDATE
            `+ s.dbname +`.items
        SET
            priority = ?,
            retryAfter = ?
        WHERE
            id = ?
    `
    
    result, err := s.db.Exec(query, priority, due.Unix(), id)
    if err != nil {
        return fmt.Errorf("error updating %s item %d: %+v", s.dbname, id, err)
    }
    if _, ok := sqlp.RowsAffectedBetween(result, 1, 1); !ok {
        return fmt.Errorf("error updating %s item %d: rows affected != 1", s.dbname, id)
    }
    return nil
}

func (s itemSqliteService) DeleteItem(id ItemID) error {
    query := `
        DELETE FROM
            `+ s.dbname +`.items
        WHERE
            id = ?
    `
    
    result, err := s.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("error deleting %s item %d: %+v", s.dbname, id, err)
    }
    if _, ok := sqlp.RowsAffectedBetween(result, 1, 1); !ok {
        return fmt.Errorf("error deleting %s item %d: rows affected != 1", s.dbname, id)
    }
    return nil
}

func (s itemSqliteService) Close() error {
    _, err := s.db.Exec(`DETACH DATABASE `+ s.dbname)
    if err != nil {
        return fmt.Errorf("error closing %s.item table: %+v", s.dbname, err)
    }
    return nil
}

func (s itemSqliteService) Delete() error {
    return os.Remove(s.dbpath)
}
