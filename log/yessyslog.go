// +build linux

package log

import (
    "log/syslog"
)

type sysl struct {
    w *syslog.Writer
}

func (s *sysl) Close() error {
    return s.w.Close()
}

func (s *sysl) Debug(m string) error {
    return s.w.Debug(m)
}

func (s *sysl) Info(m string) error {
    return s.w.Info(m)
}

func (s *sysl) Warning(m string) error {
    return s.w.Warning(m)
}

func (s *sysl) Err(m string) error {
    return s.w.Err(m)
}

func (s *sysl) Emerg(m string) error {
    return s.w.Emerg(m)
}

func (s *sysl) Crit(m string) error {
    return s.w.Crit(m)
}

func (s *sysl) Write(data []byte) (int, error) {
    return s.Write(data)
}

// Connect to a local or remote syslog, if available.
// See https://golang.org/pkg/log/syslog/#Dial
func dial(network string, addr string, priority Priority, tag string) (Syslog, error) {
    w, err := syslog.Dial(network, addr, syslog.Priority(priority), tag)
    if err != nil { return nil, err }
    return &sysl{w}, nil
}
