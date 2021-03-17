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


## Examples



```go
package main

import (
    "encoding/json"
    "fmt"
    "time"

    "tawesoft.co.uk/go/log"
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

    encodedCfg, err := json.Marshal(cfg)
    fmt.Printf("Encoded config: %s\n\n", string(encodedCfg))

    var decodedCfg log.Config
    err = json.Unmarshal(encodedCfg, &decodedCfg)
    if err != nil { panic(err) }
    fmt.Printf("Decoded encoded config: %+v\n", decodedCfg)

    if cfg != decodedCfg { panic("not equal!") }
}
```

## Getting Help

This package is part of [tawesoft.co.uk/go](https://www.tawesoft.co.uk/go),
a monorepo for small Go modules maintained by Tawesoft®.
Check out that URL for more information about other Go modules from
Tawesoft plus community and commercial support options.