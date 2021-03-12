package log

import (
    "io"
)

// Syslog is syslog.Writer as an interface
type Syslog interface {
    io.Writer
    Close() error
    Debug(m string) error
    Info(m string) error
    Warning(m string) error
    Err(m string) error
    Emerg(m string) error
    Crit(m string) error
}
