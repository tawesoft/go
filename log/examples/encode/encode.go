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
