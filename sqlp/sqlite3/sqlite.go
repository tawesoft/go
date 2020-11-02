// Package sqlite implements a mattn/go-sqlite3 database with simple setup
// of things like utf8 collation and tawesoft.co.uk/go/dev/sqlp features.
package sqlite3

import (
    "database/sql"
    "fmt"
    "strings"
    
    "github.com/mattn/go-sqlite3"
    "tawesoft.co.uk/go/sqlp"
)

// Features implements tawesoft.co.uk/go/dev/sqlp features for SQLite3
var Features = sqlp.Features{
    EscapeIdentifier: escapeIdentifier, // double quote
    EscapeString: escapeString, // single quote
    IsUniqueConstraintError: isUniqueConstraintError,
}

// Returns true iff err is a unique constraint error for SQLite
func isUniqueConstraintError(err error) bool {
    if sqliteErr, ok := err.(sqlite3.Error); ok {
        return (sqliteErr.Code == sqlite3.ErrConstraint) &&
            (sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique)
    }
    return false
}

func escapeIdentifier(s string) (string, error) {
    s = strings.ReplaceAll(s, `"`, `""`)
    return `"`+s+`"`, nil
}

func escapeString(s string) (string, error) {
    s = strings.ReplaceAll(s, `'`, `''`)
    return `'`+s+`'`, nil
}

// A Collation is a names and a collation-aware string comparison function
type Collation struct {
    Name string
    Cmp func(string, string) int
}

// Utf8Collation is a common Collation provided as a convenience
var Utf8Collation = Collation{"utf8", strings.Compare}

// JournalMode is a type of SQLite journal mode
type JournalMode string
const (
    JournalModeDelete   = JournalMode("DELETE")
    JournalModeTruncate = JournalMode("TRUNCATE")
    JournalModePersist  = JournalMode("PERSIST")
    JournalModeMemory   = JournalMode("MEMORY")
    JournalModeWAL      = JournalMode("WAL")
    JournalModeOff      = JournalMode("OFF")
)

// Config can be used to configure the SQLite connection
type Config struct {
    // ForeignKeys enables/disables the SQLite foreign_keys pragma
    ForeignKeys bool
    
    // SecureDelete enables/disables the SQLite secure_delete pragma
    SecureDelete bool
    
    // JournalMode, defaults to JournalModeWAL
    JournalMode JournalMode
}

// Register creates a new Sqlite3 database driver named DriverName
// (e.g. "sqlite3_withCollations") that implements the given Collations and
// Extensions.
//
// Currently extensions is a placeholder argument and is not implemented.
//
// Must only be called once with a given driverName, otherwise Go will panic.
func Register(driverName string, collations []Collation, extensions []interface{}) {
    // see https://godoc.org/github.com/mattn/go-sqlite3#hdr-Connection_Hook
    sql.Register(driverName,
        &sqlite3.SQLiteDriver{
            ConnectHook: func(conn *sqlite3.SQLiteConn) error {
                for _, col := range collations {
                    err := conn.RegisterCollation(col.Name, col.Cmp)
                    if err != nil {
                        return fmt.Errorf("error registering %s database collation %s: %+v",
                            driverName, col.Name, err)
                    }
                }
                return nil
            },
    })
}

func onOff(x bool) string {
    if x { return "ON" }
    return "OFF"
}

// Opener returns an Open function with the config argument already applied
func Opener(config Config) (func(string, string) (*sql.DB, error)) {
    return func(driverName string, dataSource string) (*sql.DB, error) {
        return Open(driverName, dataSource, config)
    }
}

// Open opens/creates an SQLite database with pragmas from config
//
// DriverName should match the name used in Register (or leave blank for default)
//
// DataSource is the sqlite connection string
// e.g. "file:/var/sqlite/example.sql" or ":memory:"
func Open(driverName string, dataSource string, config Config) (*sql.DB, error) {
    
    if driverName == "" {
        driverName = "sqlite3"
    }
    
    db, err := sql.Open(driverName, dataSource)
    if err != nil { return nil, err }
    
    var stmt = `
        PRAGMA foreign_keys = `+onOff(config.ForeignKeys)+`;
        PRAGMA secure_delete = `+onOff(config.SecureDelete)+`;
        PRAGMA journal_mode = `+string(config.JournalMode)+`;
    `
    _, err = db.Exec(stmt)
    if err != nil {
        db.Close();
        return nil, fmt.Errorf("error setting database PRAGMAs: %+v", err)
    }
    
    return db, nil
}
