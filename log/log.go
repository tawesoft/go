package log

import (
    "fmt"
    "os"
    "time"

    "github.com/mattn/go-isatty"
)

// ConfigSyslog configures a syslog logger
type ConfigSyslog struct {
    Enabled bool
    Network string // See syslog.Dial
    Address string // See syslog.Dial
    Priority Priority // See syslog.Dial
    Tag string // See syslog.Dial
}

func (c ConfigSyslog) Dial() (Syslog, error) {
    if !c.Enabled { return nil, fmt.Errorf("cannot connect to syslog: not enabled") }
    return dial(c.Network, c.Address, c.Priority, c.Tag)
}

// ConfigFile configures a file logger, with optional file rotation
type ConfigFile struct {
    Enabled bool

    // Mode to use when creating the file e.g. 0644, 0600
    Mode os.FileMode

    // Path to write the current (non-rotated) file. Rotated files appear
    // in the same directory.
    Path string

    // If Rotate is true, logs are rotated (e.g. like logrotate) once they
    // get to a certain size.
    Rotate bool

    // If RotateCompress is true, rotated log files are compressed (read them
    // with zless, zcat, or gunzip, for example)
    RotateCompress bool

    // A log is rotated if it would be bigger than RotateMaxSize (in bytes)
    RotateMaxSize int // bytes

    // If non-zero, delete any rotated logs older than RotateKeepAge
    RotateKeepAge time.Duration

    // If non-zero, keep only this many rotated logs and delete any exceeding
    // the limit of RotateKeepNumber.
    RotateKeepNumber int
}

// ConfigStdio configures a logger which writes to Stderr
type ConfigStderr struct {
    Enabled bool

    // If Color is true, output is colourised iff Stderr is attached to
    // a terminal.
    Color bool
}

// ShouldColorize returns true if the output should be colourised (if possible)
// for a given output (e.g. os.Stderr). This is true when both the config Color
// field is true and the output is a terminal.
func (c ConfigStderr) ShouldColorize(output *os.File) bool {
    return c.Color && isatty.IsTerminal(os.Stderr.Fd())
}

type Config struct{
    Syslog ConfigSyslog
    File   ConfigFile
    Stderr ConfigStderr
}
