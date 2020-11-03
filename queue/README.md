# queue - same-process durable message queue

```shell script
go get "tawesoft.co.uk/go/"
```

```go
import "tawesoft.co.uk/go/queue"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_queue] ∙ [docs][docs_queue] ∙ [src][src_queue] | [MIT][copy_queue] | candidate |

[home_queue]: https://tawesoft.co.uk/go/queue
[src_queue]:  https://github.com/tawesoft/go/tree/master/dialog
[docs_queue]: https://godoc.org/tawesoft.co.uk/go/queue
[copy_queue]: https://github.com/tawesoft/go/tree/master/queue/LICENSE.txt

## About

Package queue implements simple, durable/ACID, same-process message queues
supporting at-least-once or exactly-once delivery and best-effort ordering
by priority and/or time.

## Examples

See examples folder.

## Getting Help

This package is part of [tawesoft.co.uk/go](https://www.tawesoft.co.uk/go),
a monorepo for small Go modules maintained by Tawesoft®.
Check out that URL for more information about other Go modules from
Tawesoft plus community and commercial support options.
