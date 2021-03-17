package zerolog

import (
    "fmt"
    "io"
    "os"
    "strings"

    "github.com/rs/zerolog"
    "gopkg.in/natefinch/lumberjack.v2"
    "tawesoft.co.uk/go/log"
)

func New(cfg log.Config) (logger zerolog.Logger, closer func() error, err error) {

    closers := make([]func() error, 0)
    writers := make([]io.Writer, 0)

    closer = func() error {
        var errs []string
        for _, closefn := range closers {
            err := closefn()
            if err != nil {
                errs = append(errs, err.Error())
            }
        }

        if len(errs) > 0 {
            return fmt.Errorf("%d logger close errors: %s",
                len(errs), strings.Join(errs, ". Also, "))
        }

        return nil
    }

    if cfg.Stderr.Enabled {
        c := cfg.Stderr

        writers = append(writers, zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
            w.NoColor = !c.ShouldColorize(os.Stderr)
            w.Out = os.Stderr
        }))
    }

    if cfg.Syslog.Enabled {
        c := cfg.Syslog
        syslog, err := c.Dial()
        if err != nil { return zerolog.Logger{}, nil, err }

        writers = append(writers, zerolog.SyslogLevelWriter(syslog))
        closers = append(closers, func() error { return syslog.Close() })
    }

    if cfg.File.Enabled && cfg.File.Rotate {
        // touch the file first to set correct permissions
        c := cfg.File
        mode := os.FileMode(0644)
        if c.Mode != 0 { mode = c.Mode }

        f, err := os.OpenFile(c.Path, os.O_WRONLY|os.O_CREATE, mode)
        if err != nil {
            closer()
            return zerolog.Logger{}, nil, err
        }

        // but then close it for rotator to rotate if necessary
        f.Close()

        rotator := &lumberjack.Logger{
            Filename:   c.Path,
            MaxSize:    c.RotateMaxSize / (1024 * 1024),
            MaxAge:     int(c.RotateKeepAge.Hours() / 24),
            MaxBackups: c.RotateKeepNumber,
            Compress:   c.RotateCompress,
        }

        writers = append(writers, rotator)
        closers = append(closers, func() error { return rotator.Close() })
    }

    if cfg.File.Enabled && !cfg.File.Rotate {
        c := cfg.File
        mode := os.FileMode(0644)
        if c.Mode != 0 { mode = c.Mode }

        f, err := os.OpenFile(c.Path, os.O_WRONLY|os.O_CREATE, mode)
        if err != nil {
            closer()
            return zerolog.Logger{}, nil, err
        }

        writers = append(writers, f)
        closers = append(closers, func() error { return f.Close() })
    }

    mw := zerolog.MultiLevelWriter(writers...)
    return zerolog.New(mw).With().Timestamp().Logger(), closer, nil
}
