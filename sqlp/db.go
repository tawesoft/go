// Package sqlp ("SQL-plus" or "squelp!") defines helpful interfaces and
// implements extra features for Go SQL database drivers. Specific drivers are
// implemented in the subdirectories.
//
// Features include:
//
// * Open a SQLite database with foreign keys, UTF8 collation, etc. made easy
//   to avoid copy+pasting the same boilerplate into each project.
//
// * "Missing" essentials like escaping an SQL column name
//   (https://github.com/golang/go/issues/18478) or examining an SQL error for
//   properties such as IsUniqueConstraintError when inserting duplicate items
//
// * Interfaces like Queryable which is implemented by all of sql.DB, sql.Tx
//   and sql.Stmt, for performing queries regardless of if they are in a
//   transaction or not.
//
package sqlp

import (
    "context"
    "database/sql"
    "fmt"
    "syscall"
)

// Queryable is an interface describing the intersection of the methods
// implemented by sql.DB, sql.Tx, and sql.Stmt.
type Queryable interface {
    Exec(query string, args ...interface{}) (sql.Result, error)
    ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
    Query(query string, args ...interface{}) (*sql.Rows, error)
    QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
    QueryRow(query string, args ...interface{}) *sql.Row
    QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// Transactionable is an interface describing a Queryable that can start
// (possibly nested) transactions. A sql.Tx is not transactable, but a sql.DB
// is.
type Transactionable interface {
    Queryable
    Begin() (*sql.Tx, error)
    BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// Features are a common set of features with driver-specific implementations.
// e.g. `import tawesoft.co.uk/go/sqlp/sqlite3` and use sqlite3.Features
type Features struct {
    EscapeIdentifier func(s string) (string, error)
    EscapeString func(s string) (string, error)
    IsUniqueConstraintError func(err error) bool
}

// IsUniqueConstraintError returns true iff the database driver implements a
// check for unique constraint errors and the error indicates a statement
// did not execute because it would violate a uniqueness constraint e.g. when
// attempting to insert a duplicate item.
func IsUniqueConstraintError(f *Features, err error) bool {
    if f.IsUniqueConstraintError == nil { return false }
    return f.IsUniqueConstraintError(err)
}

// EscapeIdentifier returns a SQL identifier (such as a column name) for a
// given database driver.
func EscapeIdentifier(f *Features, s string) (string, error) {
    if f.EscapeIdentifier == nil { return s, fmt.Errorf("not implemented") }
    return f.EscapeIdentifier(s)
}

// EscapeString returns a quoted string literal for a given database driver.
func EscapeString(f *Features, s string) (string, error) {
    if f.EscapeString == nil { return s, fmt.Errorf("not implemented") }
    return f.EscapeString(s)
}

// OpenMode wraps a database opener (e.g. a sqlite3.Opener) with a syscall to
// set the file permissions to a unix mode when the file is created (e.g. mode
// 0006 for user read/write only) and, additionally, checks the connection using
// db.Ping().
//
// Note - NOT safe to be used concurrently with other I/O due to use of syscall
func OpenMode(
    driverName string,
    dataSource string,
    mode int,
    Opener func(string, string) (*sql.DB, error),
) (*sql.DB, error) {
    oldMask := syscall.Umask(mode)
    defer syscall.Umask(oldMask)
    
    db, err := Opener(driverName, dataSource)
    if err != nil {
        return nil, fmt.Errorf("error opening %s database %s: %+v",
            driverName, dataSource, err)
    }
    
    // For a file data source, ensure creation.
    // For a network source, ensure connection.
    err = db.Ping()
    if err != nil {
        db.Close()
        return nil, fmt.Errorf("error connecting to %s database %s: %+v",
            driverName, dataSource, err)
    }
    
    return db, nil
}

// RowsAffectedBetween returns true iff the result rows affected is not an
// error and falls between min and max (inclusive). Otherwise, returns false and
// the first argument is the number of rows actually affected or -1 and false
// if there was an error counting how many rows were affected (which should
// only ever happen if there is a bug in your code e.g. trying to count rows
// affected by a DDL command (such as a CREATE TABLE) or not checking a
// previous error and using an invalid result).
func RowsAffectedBetween(result sql.Result, min int, max int) (int64, bool) {
    affected, err := result.RowsAffected()
    if err != nil { return -1, false }
    if affected < int64(min) { return affected, false }
    if affected > int64(max) { return affected, false }
    return 0, true
}
