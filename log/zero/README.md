# log/zero - easy-config zerolog

```shell script
go get "tawesoft.co.uk/go/"
```

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

## About

Package log/zero makes it trivial to configure a zerolog logger with syslog,
file, and/or console output.


## See https://github.com/rs/zerolog



## See https://www.tawesoft.co.uk/go/doc/log



## Examples



```go
package main

import (
    "time"

    "tawesoft.co.uk/go/log"
    "tawesoft.co.uk/go/log/zero"
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

    logger, closer, err := zero.New(cfg)
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