# log - uniformly configurable loggers

```shell script
go get -u "tawesoft.co.uk/go"
```

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

## About

Package log provides a common way to quickly configure a logging implementation
with file rotation, syslog, console output, etc. for some popular logging
implementations such as zerolog.

This package defines the configuration interface, which is json-encodable.

Loggers are concretely implemented by the packages in the subfolder e.g.
tawesoft.co.uk/go/log/zerolog.

The package also wraps the stdlib syslog as an interface without it being a
compile-time constraint so that it can be imported on platforms that don't
support syslog (like Windows), giving a runtime error if used instead.

## Getting Help

This package is part of [tawesoft.co.uk/go](https://www.tawesoft.co.uk/go),
a monorepo for small Go modules maintained by Tawesoft®.
Check out that URL for more information about other Go modules from
Tawesoft plus community and commercial support options.