# router - general purpose (HTTP, etc.) router

## About

Package router is a general purpose router of methods (e.g. HTTP "GET") and paths (e.g. "/user/123/profile") to
some value e.g. a controller.

Supports named routes, route parameters, constructing a path from a route, etc.

Although built with HTTP routing in mind, this is a general purpose implementation that can route to any type
of value - it is not limited to HTTP handlers.

|  Links  | License | Stable? | 
|:-------:|:-------:|:-------:| 
| [home][home_] ∙ [docs][docs_] ∙ [src][src_] | [MIT][copy_] | ✔ yes |

[home_]: https://tawesoft.co.uk/go/router
[src_]:  https://github.com/tawesoft/go/tree/master/router
[docs_]: https://godoc.org/tawesoft.co.uk/go/router
[copy_]: https://github.com/tawesoft/go/tree/master/router/_COPYING.md

## Download

```shell script
go get -u tawesoft.co.uk/go
```

## Import

```
import tawesoft.co.uk/go/router
```