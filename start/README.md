# start - system process setup

```shell script
go get "tawesoft.co.uk/go/"
```

```go
import "tawesoft.co.uk/go/start"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_start] ∙ [docs][docs_start] ∙ [src][src_start] | [MIT][copy_start] | candidate |

[home_start]: https://tawesoft.co.uk/go/start
[src_start]:  https://github.com/tawesoft/go/tree/master/start
[docs_start]: https://godoc.org/tawesoft.co.uk/go/start
[copy_start]: https://github.com/tawesoft/go/tree/master/start/LICENSE.txt

## About

Package start implements helpful features for starting a (system) process
including dropping privileges (while inheriting privileged file handles) and
managing multiple goroutines with graceful shutdown.

## Getting Help

This package is part of [tawesoft.co.uk/go](https://www.tawesoft.co.uk/go),
a monorepo for small Go modules maintained by Tawesoft®.
Check out that URL for more information about other Go modules from
Tawesoft plus community and commercial support options.