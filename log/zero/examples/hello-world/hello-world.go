package main

import (
    "encoding/json"
    "fmt"
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

    encodedCfg, err := json.Marshal(cfg)
    fmt.Println("Encoded config: ", string(encodedCfg))

    logger, closer, err := zero.New(cfg)
    if err != nil { panic(err) }
    defer closer()

    logger.Info().Msg("Hello world!")
}
