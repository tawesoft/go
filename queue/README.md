# queue - same-process durable message queue

```shell script
go get -u "tawesoft.co.uk/go"
```

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

## About

Package queue implements simple, durable/ACID, same-process message queues
with best-effort ordering by priority and/or time.


## Examples


See examples folder.

## Changes

### 2021-07-06

* The Queue RetryItem method now takes an Attempt parameter.

* Calling the Delete() method on a Queue now attempts to avoid deleting
an in-memory database opened as ":memory:".


## Getting Help

This package is part of [tawesoft.co.uk/go](https://www.tawesoft.co.uk/go),
a monorepo for small Go modules maintained by Tawesoft®.
Check out that URL for more information about other Go modules from
Tawesoft plus community and commercial support options.