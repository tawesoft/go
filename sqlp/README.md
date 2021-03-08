# sqlp - SQL database extras

```shell script
go get "tawesoft.co.uk/go/"
```

```go
import "tawesoft.co.uk/go/sqlp"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_sqlp] ∙ [docs][docs_sqlp] ∙ [src][src_sqlp] | [MIT][copy_sqlp] | candidate |

[home_sqlp]: https://tawesoft.co.uk/go/sqlp
[src_sqlp]:  https://github.com/tawesoft/go/tree/master/sqlp
[docs_sqlp]: https://www.tawesoft.co.uk/go/doc/sqlp
[copy_sqlp]: https://github.com/tawesoft/go/tree/master/sqlp/LICENSE.txt

## About

Package sqlp ("SQL-plus" or "squelp!") defines helpful interfaces and
implements extra features for Go SQL database drivers. Specific driver
extras are implemented in the subdirectories.


## Features


* Open a SQLite database with foreign keys, UTF8 collation, etc. made easy
to avoid copy+pasting the same boilerplate into each project.

* "Missing" essentials like escaping an SQL column name
(https://github.com/golang/go/issues/18478) or examining an SQL error for
properties such as IsUniqueConstraintError when inserting duplicate items

* Interfaces like Queryable which is implemented by all of sql.DB, sql.Tx
and sql.Stmt, for performing queries regardless of if they are in a
transaction or not.


## Driver extras


* tawesoft.co.uk/go/sqlp/sqlite3 (mattn/go-sqlite3)

## Getting Help

This package is part of [tawesoft.co.uk/go](https://www.tawesoft.co.uk/go),
a monorepo for small Go modules maintained by Tawesoft®.
Check out that URL for more information about other Go modules from
Tawesoft plus community and commercial support options.