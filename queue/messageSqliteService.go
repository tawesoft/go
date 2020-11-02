package queue

import (
    "database/sql"
    "fmt"
    
    "tawesoft.co.uk/go/sqlp"
)

type messageSqliteService struct {dbname string}
var nilMessageSqliteService = messageSqliteService{}

func newMessageSqliteService(db *sql.DB, dbname string) (messageSqliteService, error) {
    err := initMessageSqliteService(db, dbname)
    if err != nil { return nilMessageSqliteService, err }
    
    return messageSqliteService{dbname}, nil
}

func initMessageSqliteService(db *sql.DB, dbname string) error {
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS `+ dbname +`.messages (
            id   INTEGER PRIMARY KEY NOT NULL REFERENCES items(id) ON DELETE CASCADE,
            data BLOB NOT NULL
        );
    `)
    if err != nil {
        return fmt.Errorf("error initialising messages table: %+v", err)
    }
    return nil
}

func (s messageSqliteService) Create(q sqlp.Queryable, itemID ItemID, message string) error {
    query := `INSERT INTO `+ s.dbname +`.messages(id, data) VALUES (?, ?)`
    
    result, err := q.Exec(query, itemID, message)
    if err != nil {
        return fmt.Errorf("error inserting message: %+v", err)
    }
    
    if n, ok := sqlp.RowsAffectedBetween(result, 1, 1); !ok {
        return fmt.Errorf("error creating message %d: rows affected %d != 1", itemID, n)
    }
    
    return nil
}

