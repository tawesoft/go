# grace - start and gracefully shutdown processes

```shell script
go get "tawesoft.co.uk/go/"
```

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

## About

Package grace implements a simple way to start multiple long-lived processes
(e.g. goroutines) with cancellation, signal handling and graceful shutdown.

## Getting Help

This package is part of [tawesoft.co.uk/go](https://www.tawesoft.co.uk/go),
a monorepo for small Go modules maintained by Tawesoft®.
Check out that URL for more information about other Go modules from
Tawesoft plus community and commercial support options.