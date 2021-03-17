package log

import (
    "encoding/json"
    "fmt"
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/mattn/go-isatty"
    "golang.org/x/text/language"
    "tawesoft.co.uk/go/humanizex"
)

// ConfigSyslog configures a syslog logger
type ConfigSyslog struct {
    Enabled  bool
    Network  string // See syslog.Dial
    Address  string // See syslog.Dial
    Priority Priority // See syslog.Dial
    Tag      string // See syslog.Dial
}

type nativeConfigSyslog struct {
    Enabled  bool
    Network  string
    Address  string
    Priority string
    Tag      string
}

// Dial connects to a local or remote syslog.
func (c ConfigSyslog) Dial() (Syslog, error) {
    if !c.Enabled { return nil, fmt.Errorf("cannot connect to syslog: not enabled") }
    return dial(c.Network, c.Address, c.Priority, c.Tag)
}

func (c ConfigSyslog) MarshalJSON() ([]byte, error) {
    return json.Marshal(nativeConfigSyslog{
        Enabled:   c.Enabled,
        Network:   c.Network,
        Address:   c.Address,
        Priority:  c.Priority.String(),
        Tag:       c.Tag,
    })
}

func (c *ConfigSyslog) UnmarshalJSON(data []byte) error {
    var n nativeConfigSyslog
    err := json.Unmarshal(data, &n)
    if err != nil { return err }

    priority, err := ParsePriority(n.Priority)
    if err != nil {
        return fmt.Errorf("error parsing syslog priority: %v", err)
    }

    *c = ConfigSyslog{
        Enabled:  n.Enabled,
        Network:  n.Network,
        Address:  n.Address,
        Priority: priority,
        Tag:      n.Tag,
    }

    return nil
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

type nativeConfigFile struct {
    Enabled bool
    Mode string
    Path string
    Rotate bool
    RotateCompress bool
    RotateMaxSize string
    RotateKeepAge string
    RotateKeepNumber int
}

func (c ConfigFile) MarshalJSON() ([]byte, error) {
    h := humanizex.NewHumanizer(language.English)

    return json.Marshal(nativeConfigFile{
        Enabled:          c.Enabled,
        Mode:             fmt.Sprintf("%#o", c.Mode),
        Path:             c.Path,
        Rotate:           c.Rotate,
        RotateCompress:   c.RotateCompress,
        RotateMaxSize:    h.FormatBytesIEC(int64(c.RotateMaxSize)),
        RotateKeepAge:    h.FormatDuration(c.RotateKeepAge),
        RotateKeepNumber: c.RotateKeepNumber,
    })
}

func (c *ConfigFile) UnmarshalJSON(data []byte) error {
    var n nativeConfigFile
    err := json.Unmarshal(data, &n)
    if err != nil { return err }

    h := humanizex.NewHumanizer(language.English)

    rotateMaxSize, err := h.ParseBytesIEC(n.RotateMaxSize)
    if err != nil { return err }

    rotateKeepAge, err := h.ParseDuration(n.RotateKeepAge)
    if err != nil { return err }

    var mode int64
    if strings.HasPrefix(n.Mode, "0") {
        var err error
        mode, err = strconv.ParseInt(n.Mode, 8, 32)
        if err != nil { return fmt.Errorf("invalid octal mode %s\n") }
    } else {
        mode, err = strconv.ParseInt(n.Mode, 10, 32)
        if err != nil { return fmt.Errorf("invalid mode %s\n") }
    }

    *c = ConfigFile{
        Enabled:          n.Enabled,
        Mode:             os.FileMode(mode),
        Path:             n.Path,
        Rotate:           n.Rotate,
        RotateCompress:   n.RotateCompress,
        RotateMaxSize:    int(rotateMaxSize),
        RotateKeepAge:    rotateKeepAge,
        RotateKeepNumber: n.RotateKeepNumber,
    }

    return nil
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
