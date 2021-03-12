package log

import (
    "fmt"
    "strings"
)

// Because syslog cannot be built on Windows, we want to avoid importing
// the package (making it optional) and instead duplicate the constants
// from https://pkg.go.dev/log/syslog#Priority

// Priority is a combination of the syslog facility and severity.
// See https://pkg.go.dev/log/syslog#Priority
type Priority int

const (
    LOG_EMERG Priority = iota
    LOG_ALERT
    LOG_CRIT
    LOG_ERR
    LOG_WARNING
    LOG_NOTICE
    LOG_INFO
    LOG_DEBUG
)

const (
    LOG_KERN Priority = iota << 3
    LOG_USER
    LOG_MAIL
    LOG_DAEMON
    LOG_AUTH
    // LOG_SYSLOG // omitted - should only ever be used internally by syslogd
    LOG_LPR
    LOG_NEWS
    LOG_UUCP
    LOG_CRON
    LOG_AUTHPRIV
    LOG_FTP

    LOG_LOCAL0
    LOG_LOCAL1
    LOG_LOCAL2
    LOG_LOCAL3
    LOG_LOCAL4
    LOG_LOCAL5
    LOG_LOCAL6
    LOG_LOCAL7
)

func (p *Priority) MarshalJSON() ([]byte, error) {
    return []byte(p.String()), nil
}

func (p *Priority) UnmarshalJSON(data []byte) error {
    number, err := ParsePriority(string(data))
    if err != nil { return err }

    *p = number
    return nil
}

func (p Priority) String() string {
    var r []string

    if p & LOG_EMERG    == LOG_EMERG    { r = append(r, "EMERG") }
    if p & LOG_ALERT    == LOG_ALERT    { r = append(r, "ALERT") }
    if p & LOG_CRIT     == LOG_CRIT     { r = append(r, "CRIT") }
    if p & LOG_ERR      == LOG_ERR      { r = append(r, "ERR") }
    if p & LOG_WARNING  == LOG_WARNING  { r = append(r, "WARNING") }
    if p & LOG_NOTICE   == LOG_NOTICE   { r = append(r, "NOTICE") }
    if p & LOG_INFO     == LOG_INFO     { r = append(r, "INFO") }
    if p & LOG_DEBUG    == LOG_DEBUG    { r = append(r, "DEBUG") }

    if p & LOG_KERN     == LOG_KERN     { r = append(r, "KERN") }
    if p & LOG_USER     == LOG_USER     { r = append(r, "USER") }
    if p & LOG_MAIL     == LOG_MAIL     { r = append(r, "MAIL") }
    if p & LOG_DAEMON   == LOG_DAEMON   { r = append(r, "DAEMON") }
    if p & LOG_AUTH     == LOG_AUTH     { r = append(r, "AUTH") }
    if p & LOG_LPR      == LOG_LPR      { r = append(r, "LPR") }
    if p & LOG_NEWS     == LOG_NEWS     { r = append(r, "NEWS") }
    if p & LOG_UUCP     == LOG_UUCP     { r = append(r, "UUCP") }
    if p & LOG_CRON     == LOG_CRON     { r = append(r, "CRON") }
    if p & LOG_AUTHPRIV == LOG_AUTHPRIV { r = append(r, "AUTHPRIV") }
    if p & LOG_FTP      == LOG_FTP      { r = append(r, "FTP") }

    if p & LOG_LOCAL0   == LOG_LOCAL0   { r = append(r, "LOCAL0") }
    if p & LOG_LOCAL1   == LOG_LOCAL1   { r = append(r, "LOCAL1") }
    if p & LOG_LOCAL2   == LOG_LOCAL2   { r = append(r, "LOCAL2") }
    if p & LOG_LOCAL3   == LOG_LOCAL3   { r = append(r, "LOCAL3") }
    if p & LOG_LOCAL4   == LOG_LOCAL4   { r = append(r, "LOCAL4") }
    if p & LOG_LOCAL5   == LOG_LOCAL5   { r = append(r, "LOCAL5") }
    if p & LOG_LOCAL6   == LOG_LOCAL6   { r = append(r, "LOCAL6") }
    if p & LOG_LOCAL7   == LOG_LOCAL7   { r = append(r, "LOCAL7") }

    if len(r) > 0 {
        return strings.Join(r, "|")
    } else {
        return ""
    }
}

func ParsePriority(str string) (Priority, error) {
    // check against longest possible sensible string
    if len(str) > len("LOG_WARNING | LOG_AUTHPRIV") {
        return 0, fmt.Errorf("syslog priority string too long")
    }

    // defaults from https://golang.org/pkg/log/syslog/#Priority
    facility := LOG_KERN
    level := LOG_EMERG

    parts := strings.SplitN(str, "|", 2)

    for _, p := range parts {
        p = strings.TrimSpace(p)

        switch p {
            case "EMERG":    level = LOG_EMERG
            case "ALERT":    level = LOG_ALERT
            case "CRIT":     level = LOG_CRIT
            case "ERR":      level = LOG_ERR
            case "WARNING":  level = LOG_WARNING
            case "NOTICE":   level = LOG_NOTICE
            case "INFO":     level = LOG_INFO
            case "DEBUG":    level = LOG_DEBUG

            case "KERN":     facility = LOG_KERN
            case "USER":     facility = LOG_USER
            case "MAIL":     facility = LOG_MAIL
            case "DAEMON":   facility = LOG_DAEMON
            case "AUTH":     facility = LOG_AUTH
            case "LPR":      facility = LOG_LPR
            case "NEWS":     facility = LOG_NEWS
            case "UUCP":     facility = LOG_UUCP
            case "CRON":     facility = LOG_CRON
            case "AUTHPRIV": facility = LOG_AUTHPRIV
            case "FTP":      facility = LOG_FTP

            case "LOCAL0":   facility = LOG_LOCAL0
            case "LOCAL1":   facility = LOG_LOCAL1
            case "LOCAL2":   facility = LOG_LOCAL2
            case "LOCAL3":   facility = LOG_LOCAL3
            case "LOCAL4":   facility = LOG_LOCAL4
            case "LOCAL5":   facility = LOG_LOCAL5
            case "LOCAL6":   facility = LOG_LOCAL6
            case "LOCAL7":   facility = LOG_LOCAL7

            default:
                return 0, fmt.Errorf("invalid syslog severity/facility string %q", str)
        }
    }

    return facility | level, nil
}
