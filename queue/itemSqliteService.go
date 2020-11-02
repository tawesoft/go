package queue

import (
    "database/sql"
    "fmt"
    "os"
    "time"
    
    "tawesoft.co.uk/go/sqlp"
    "tawesoft.co.uk/go/sqlp/sqlite3"
)

type itemSqliteService struct {
    dbname  string
    dbpath  string
    db      *sql.DB
    uuidSvc UUIDService
    msgSvc  messageSqliteService
}
var nilItemSqliteService = itemSqliteService{}

type itemSqliteNew struct {
    UUID       []byte
    Created    int64 // time.Time // UTC
    RetryAfter int64 // time.Time // UTC
}
func (p itemSqliteNew) fields() string {
    return "uuid, created, retryAfter"
}
func (p itemSqliteNew) placeholders() string {
    return "?, ?, ?"
}
func (p itemSqliteNew) values() []interface{} {
    return []interface{}{p.UUID, p.Created, p.RetryAfter}
}

type itemSqlite struct {
    ID ItemID
    UUID []byte
    Message string
    Attempt int
    Created int64 // time.Time // UTC
    RetryAfter int64 // time.Time // UTC
}
func (p itemSqlite) fields() string {
    return "items.id, messages.data, items.uuid, items.attempt, items.created, items.retryAfter"
}
func (p itemSqlite) placeholders() string {
    return "?, ?, ?, ?, ?, ?"
}
func (p itemSqlite) values() []interface{} {
    return []interface{}{p.ID, p.Message, p.UUID, p.Attempt, p.Created, p.RetryAfter}
}
func (p *itemSqlite) pointers() []interface{} {
    return []interface{}{&p.ID, &p.Message, &p.UUID, &p.Attempt, &p.Created, &p.RetryAfter}
}
func (p itemSqlite) toItem() Item {
    return Item{
        ID:         p.ID,
        UUID:       p.UUID,
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
            uuid       BLOB NOT NULL,
            attempt    INTEGER NOT NULL DEFAULT 0,
            created    INTEGER NOT NULL, -- epoch seconds
            retryAfter INTEGER NOT NULL  -- epoch seconds
        );

        CREATE UNIQUE INDEX IF NOT EXISTS `+ dbname +`.items_idx_nonce ON items(uuid);
        CREATE        INDEX IF NOT EXISTS `+ dbname +`.items_idx_retryAfter ON items(retryAfter);
    `)
    if err != nil {
        return fmt.Errorf("error initialising item table: %+v", err)
    }
    return nil
}

// newItemSqliteService creates a new itemService implemented by a
// itemSqliteService
func newItemSqliteService(db *sql.DB, uuidSvc UUIDService, name string, rawpath string) (itemSqliteService, error) {
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
        uuidSvc: uuidSvc,
        msgSvc:  msgSvc,
    }, nil
}

func (s itemSqliteService) CreateItem(newItem NewItem) error {
    err := func() error {
        tx, err := s.db.Begin()
        if err != nil { return err }
        defer tx.Rollback()
        
        uuid, err := s.uuidSvc.Generate()
        if err != nil { return fmt.Errorf("UUID error: %+v", err) }
        
        i := itemSqliteNew{
            UUID:       uuid,
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

func (s itemSqliteService) PeekItems(n int, due time.Time) ([]Item, error) {
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
                items.retryAfter <= ?
        ORDER BY
            items.retryAfter,
            items.id
        LIMIT ?`
        
        rows, err := s.db.Query(query, due.Unix(), n)
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

/*
// selectRow is a general purpose SELECT for a single Session
func (s *SessionDBService) selectRow(query string, args...interface{}) *Session {
    result := &SessionDB{}
    query = fmt.Sprintf("SELECT %s FROM sessions\n%s\nLIMIT 1", sessionDBFields, query)
    err := s.db.QueryRow(query, args...).Scan(
        &result.ID, &result.Key, &result.User, &result.RememberMe,
        &result.Location, &result.UserAgentGuess,
        &result.CreatedAt, &result.RefreshedAt, &result.DeletedAt)
    if err == sql.ErrNoRows { return nil }
    if err != nil { panic(err) }
    
    return result.ToSession()
}

func (s *SessionDBService) ByID(id SessionID) *Session {
    session := s.selectRow(`WHERE id == ?`, id)
    return session
}
 */

func (s itemSqliteService) Close() error {
    _, err := s.db.Exec(`
        DETACH DATABASE `+s.dbname +`;
    `)
    if err != nil {
        return fmt.Errorf("error closing %s.item table: %+v", s.dbname, err)
    }
    return nil
}

func (s itemSqliteService) Delete() error {
    return os.Remove(s.dbpath)
}
