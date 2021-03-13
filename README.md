[![Tawesoft](https://www.tawesoft.co.uk/media/0/logo-240r.png)](https://tawesoft.co.uk/go)
================================================================================

A monorepo for small Go modules maintained by [Tawesoft®](https://www.tawesoft.co.uk/go)

This is permissively-licensed open source software but exact licenses may vary between modules.

Download
--------

    go get -u tawesoft.co.uk/go

Contents
--------


### dialog - simple cross-platform messagebox

Package dialog implements simple cross platform native MessageBox/Alert
dialogs for Go.

```go
import "tawesoft.co.uk/go/dialog"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_dialog] ∙ [docs][docs_dialog] ∙ [src][src_dialog] | [MIT-0][copy_dialog] | ✔ yes |

[home_dialog]: https://tawesoft.co.uk/go/dialog
[src_dialog]:  https://github.com/tawesoft/go/tree/master/dialog
[docs_dialog]: https://www.tawesoft.co.uk/go/doc/dialog
[copy_dialog]: https://github.com/tawesoft/go/tree/master/dialog/LICENSE.txt
### drop - drop privileges and inherit handles

Package drop implements the ability to start a process as root, open
privileged resources as files, drop privileges to become a given user account,
and inherit file handles across the dropping of privileges.

```go
import "tawesoft.co.uk/go/drop"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_drop] ∙ [docs][docs_drop] ∙ [src][src_drop] | [MIT][copy_drop] | candidate |

[home_drop]: https://tawesoft.co.uk/go/drop
[src_drop]:  https://github.com/tawesoft/go/tree/master/drop
[docs_drop]: https://www.tawesoft.co.uk/go/doc/drop
[copy_drop]: https://github.com/tawesoft/go/tree/master/drop/LICENSE.txt
### email - format multipart MIME email

Package email implements the formatting of multipart MIME e-mail messages,
including Unicode headers, attachments, HTML email, and plain text.

```go
import "tawesoft.co.uk/go/email"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_email] ∙ [docs][docs_email] ∙ [src][src_email] | [MIT][copy_email] | candidate |

[home_email]: https://tawesoft.co.uk/go/email
[src_email]:  https://github.com/tawesoft/go/tree/master/email
[docs_email]: https://www.tawesoft.co.uk/go/doc/email
[copy_email]: https://github.com/tawesoft/go/tree/master/email/LICENSE.txt
### glcaps - read and check OpenGL capabilities

Package glcaps provides a nice interface to declare OpenGL capabilities you
care about, including minimum required extensions or capabilities. Glcaps has
no dependencies and is agnostic to the exact OpenGL binding used.

```go
import "tawesoft.co.uk/go/glcaps"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_glcaps] ∙ [docs][docs_glcaps] ∙ [src][src_glcaps] | [MIT][copy_glcaps] | ✔ yes |

[home_glcaps]: https://tawesoft.co.uk/go/glcaps
[src_glcaps]:  https://github.com/tawesoft/go/tree/master/glcaps
[docs_glcaps]: https://www.tawesoft.co.uk/go/doc/glcaps
[copy_glcaps]: https://github.com/tawesoft/go/tree/master/glcaps/LICENSE.txt
### grace - start and gracefully shutdown processes

Package grace implements a simple way to start multiple long-lived processes
(e.g. goroutines) with cancellation, signal handling and graceful shutdown.

```go
import "tawesoft.co.uk/go/grace"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_grace] ∙ [docs][docs_grace] ∙ [src][src_grace] | [MIT][copy_grace] | candidate |

[home_grace]: https://tawesoft.co.uk/go/grace
[src_grace]:  https://github.com/tawesoft/go/tree/master/grace
[docs_grace]: https://www.tawesoft.co.uk/go/doc/grace
[copy_grace]: https://github.com/tawesoft/go/tree/master/grace/LICENSE.txt
### loader - concurrent dependency graph solver

Package loader implements the ability to define a graph of tasks and
dependencies, classes of synchronous and concurrent workers, and limiting
strategies, and solve the graph incrementally or totally.

```go
import "tawesoft.co.uk/go/loader"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_loader] ∙ [docs][docs_loader] ∙ [src][src_loader] | [MIT][copy_loader] | ✘ **no** |

[home_loader]: https://tawesoft.co.uk/go/loader
[src_loader]:  https://github.com/tawesoft/go/tree/master/loader
[docs_loader]: https://www.tawesoft.co.uk/go/doc/loader
[copy_loader]: https://github.com/tawesoft/go/tree/master/loader/LICENSE.txt
### log - uniformly configurable loggers

Package log provides a common way to quickly configure a logging implementation
with file rotation, syslog, console output, etc. for some popular logging
implementations such as zerolog.

```go
import "tawesoft.co.uk/go/log"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_log] ∙ [docs][docs_log] ∙ [src][src_log] | [MIT-0][copy_log] | candidate |

[home_log]: https://tawesoft.co.uk/go/log
[src_log]:  https://github.com/tawesoft/go/tree/master/log
[docs_log]: https://www.tawesoft.co.uk/go/doc/log
[copy_log]: https://github.com/tawesoft/go/tree/master/log/LICENSE.txt
### log/zero - easy-config zerolog

Package log/zero makes it trivial to configure a zerolog logger with syslog,
rotating file, and/or console output using the same uniform configuration
interface.

```go
import "tawesoft.co.uk/go/log/zero"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_log/zero] ∙ [docs][docs_log/zero] ∙ [src][src_log/zero] | [MIT-0][copy_log/zero] | candidate |

[home_log/zero]: https://tawesoft.co.uk/go/log/zero
[src_log/zero]:  https://github.com/tawesoft/go/tree/master/log/zero
[docs_log/zero]: https://www.tawesoft.co.uk/go/doc/log/zero
[copy_log/zero]: https://github.com/tawesoft/go/tree/master/log/zero/LICENSE.txt
### lxstrconv - locale-aware number parsing

Package lxstrconv is an attempt at implementing locale-aware parsing of
numbers that integrates with golang.org/x/text.

```go
import "tawesoft.co.uk/go/lxstrconv"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_lxstrconv] ∙ [docs][docs_lxstrconv] ∙ [src][src_lxstrconv] | [MIT][copy_lxstrconv] | ✘ **no** |

[home_lxstrconv]: https://tawesoft.co.uk/go/lxstrconv
[src_lxstrconv]:  https://github.com/tawesoft/go/tree/master/lxstrconv
[docs_lxstrconv]: https://www.tawesoft.co.uk/go/doc/lxstrconv
[copy_lxstrconv]: https://github.com/tawesoft/go/tree/master/lxstrconv/LICENSE.txt
### operator - operators as functions

Package operator implements logical, arithmetic, bitwise and comparison
operators as functions (like the Python operator module). Includes unary,
binary, and n-ary functions with overflow checked variants.

```go
import "tawesoft.co.uk/go/operator"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_operator] ∙ [docs][docs_operator] ∙ [src][src_operator] | [MIT-0][copy_operator] | ✔ yes |

[home_operator]: https://tawesoft.co.uk/go/operator
[src_operator]:  https://github.com/tawesoft/go/tree/master/operator
[docs_operator]: https://www.tawesoft.co.uk/go/doc/operator
[copy_operator]: https://github.com/tawesoft/go/tree/master/operator/LICENSE.txt
### queue - same-process durable message queue

Package queue implements simple, durable/ACID, same-process message queues
supporting at-least-once or exactly-once delivery and best-effort ordering
by priority and/or time.

```go
import "tawesoft.co.uk/go/queue"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_queue] ∙ [docs][docs_queue] ∙ [src][src_queue] | [MIT][copy_queue] | candidate |

[home_queue]: https://tawesoft.co.uk/go/queue
[src_queue]:  https://github.com/tawesoft/go/tree/master/queue
[docs_queue]: https://www.tawesoft.co.uk/go/doc/queue
[copy_queue]: https://github.com/tawesoft/go/tree/master/queue/LICENSE.txt
### router - general purpose (HTTP, etc.) router

Package router is a general purpose router of methods (e.g. HTTP "GET") and
paths (e.g. "/user/123/profile") to some value e.g. a controller.

```go
import "tawesoft.co.uk/go/router"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_router] ∙ [docs][docs_router] ∙ [src][src_router] | [MIT][copy_router] | candidate |

[home_router]: https://tawesoft.co.uk/go/router
[src_router]:  https://github.com/tawesoft/go/tree/master/router
[docs_router]: https://www.tawesoft.co.uk/go/doc/router
[copy_router]: https://github.com/tawesoft/go/tree/master/router/LICENSE.txt
### sqlp - SQL database extras

Package sqlp ("SQL-plus" or "squelp!") defines helpful interfaces and
implements extra features for Go SQL database drivers. Specific driver
extras are implemented in the subdirectories.

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
### sqlp/sqlite3 - SQLite3 database extras

Package sqlite enchances a mattn/go-sqlite3 database with simple setup
of things like utf8 collation and tawesoft.co.uk/go/sqlp features.

```go
import "tawesoft.co.uk/go/sqlp/sqlite3"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_sqlp/sqlite3] ∙ [docs][docs_sqlp/sqlite3] ∙ [src][src_sqlp/sqlite3] | [MIT][copy_sqlp/sqlite3] | candidate |

[home_sqlp/sqlite3]: https://tawesoft.co.uk/go/sqlp/sqlite3
[src_sqlp/sqlite3]:  https://github.com/tawesoft/go/tree/master/sqlp/sqlite3
[docs_sqlp/sqlite3]: https://www.tawesoft.co.uk/go/doc/sqlp/sqlite3
[copy_sqlp/sqlite3]: https://github.com/tawesoft/go/tree/master/sqlp/sqlite3/LICENSE.txt
### variadic - helpers for variadic functions

Package variadic implements features that make it easier to work with
variadic functions.

```go
import "tawesoft.co.uk/go/variadic"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_variadic] ∙ [docs][docs_variadic] ∙ [src][src_variadic] | [MIT][copy_variadic] | candidate |

[home_variadic]: https://tawesoft.co.uk/go/variadic
[src_variadic]:  https://github.com/tawesoft/go/tree/master/variadic
[docs_variadic]: https://www.tawesoft.co.uk/go/doc/variadic
[copy_variadic]: https://github.com/tawesoft/go/tree/master/variadic/LICENSE.txt
### ximage - extended image types

Package ximage implements Red, RG, and RGB images matching the core
image interface.

```go
import "tawesoft.co.uk/go/ximage"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_ximage] ∙ [docs][docs_ximage] ∙ [src][src_ximage] | [BSD-3-Clause][copy_ximage] | ✔ yes |

[home_ximage]: https://tawesoft.co.uk/go/ximage
[src_ximage]:  https://github.com/tawesoft/go/tree/master/ximage
[docs_ximage]: https://www.tawesoft.co.uk/go/doc/ximage
[copy_ximage]: https://github.com/tawesoft/go/tree/master/ximage/LICENSE.txt
### ximage/xcolor - extended color types

Package xcolor implements Red, RedGreen, and RGB color models matching the core
image/color interface.

```go
import "tawesoft.co.uk/go/ximage/xcolor"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_ximage/xcolor] ∙ [docs][docs_ximage/xcolor] ∙ [src][src_ximage/xcolor] | [BSD-3-Clause][copy_ximage/xcolor] | ✔ yes |

[home_ximage/xcolor]: https://tawesoft.co.uk/go/ximage/xcolor
[src_ximage/xcolor]:  https://github.com/tawesoft/go/tree/master/ximage/xcolor
[docs_ximage/xcolor]: https://www.tawesoft.co.uk/go/doc/ximage/xcolor
[copy_ximage/xcolor]: https://github.com/tawesoft/go/tree/master/ximage/xcolor/LICENSE.txt

Links
-----

* Home: [tawesoft.co.uk/go](https://tawesoft.co.uk/go)
* Docs hub: [tawesoft.co.uk/go/doc/](https://www.tawesoft.co.uk/go/doc/)
* Repository: [github.com/tawesoft/go](https://github.com/tawesoft/go)
* Or [search "tawesoft"](https://pkg.go.dev/search?q=tawesoft) on [go.dev](https://go.dev/)

Support
-------

### Free and Community Support

* [GitHub issues](https://github.com/tawesoft/go/issues)
* Email open-source@tawesoft.co.uk (feedback welcomed, but support is "best
 effort")

### Commercial Support

Open source software from Tawesoft® backed by commercial support options.

Email open-source@tawesoft.co.uk or visit [tawesoft.co.uk/products/open-source-software](https://www.tawesoft.co.uk/products/open-source-software)
to learn more.
