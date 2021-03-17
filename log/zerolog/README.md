# log/zerolog - easy-config zerolog

```shell script
go get -u "tawesoft.co.uk/go"
```

```go
import "tawesoft.co.uk/go/log/zerolog"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_log/zerolog] ∙ [docs][docs_log/zerolog] ∙ [src][src_log/zerolog] | [MIT-0][copy_log/zerolog] | candidate |

[home_log/zerolog]: https://tawesoft.co.uk/go/log/zerolog
[src_log/zerolog]:  https://github.com/tawesoft/go/tree/master/log/zerolog
[docs_log/zerolog]: https://www.tawesoft.co.uk/go/doc/log/zerolog
[copy_log/zerolog]: https://github.com/tawesoft/go/tree/master/log/zerolog/LICENSE.txt

## About

Package log/zerolog makes it trivial to configure a zerolog logger with syslog,
rotating file, and/or console output using the same uniform configuration
interface.


## See https://github.com/rs/zerolog


Log rotation provided by https://gopkg.in/natefinch/lumberjack.v2/


## See https://www.tawesoft.co.uk/go/doc/log



## Examples



```go
package main

import (
    "time"

    "tawesoft.co.uk/go/log"
    "tawesoft.co.uk/go/log/zerolog"
)

func main() {
    cfg := log.Config{
        Syslog: log.ConfigSyslog{
            Enabled:  true,
            Network:  "", // local
            Address:  "", // local
            Priority: log.LOG_ERR | log.LOG_DAEMON,
            Tag:      "example",
        },
        File:   log.ConfigFile{
            Enabled:          true,
            Mode:             0600,
            Path:             "example.log",
            Rotate:           true,
            RotateCompress:   true,
            RotateMaxSize:    64 * 1024 * 1024, // 64MB
            RotateKeepAge:    30 * 24 * time.Hour,
            RotateKeepNumber: 32, // 32 * 64 MB = 2 GB max storage (before compression)
        },
        Stderr: log.ConfigStderr{
            Enabled: true,
            Color:   true,
        },
    }

    logger, closer, err := zerolog.New(cfg)
    if err != nil { panic(err) }
    defer closer()

    logger.Info().Msg("Hello world!")
}
```

## Getting Help

This package is part of [tawesoft.co.uk/go](https://www.tawesoft.co.uk/go),
a monorepo for small Go modules maintained by Tawesoft®.
Check out that URL for more information about other Go modules from
Tawesoft plus community and commercial support options.