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
    LOG_SYSLOG // should only ever be used internally by syslogd
    LOG_LPR
    LOG_NEWS
    LOG_UUCP
    LOG_CRON
    LOG_AUTHPRIV
    LOG_FTP
    _ // unused
    _ // unused
    _ // unused
    _ // unused
    LOG_LOCAL0
    LOG_LOCAL1
    LOG_LOCAL2
    LOG_LOCAL3
    LOG_LOCAL4
    LOG_LOCAL5
    LOG_LOCAL6
    LOG_LOCAL7
)

/*
func (p Priority) MarshalJSON() ([]byte, error) {
    return []byte(p.String()), nil
}

func (p *Priority) UnmarshalJSON(data []byte) error {
    number, err := ParsePriority(string(data))
    if err != nil { return err }

    *p = number
    return nil
}
*/

func (p Priority) String() string {
    var r []string

    const severityMask = 0b0111
    const facilityMask = 0b1111111111111000

    if p & severityMask  == LOG_EMERG    { r = append(r, "EMERG") }
    if p & severityMask  == LOG_ALERT    { r = append(r, "ALERT") }
    if p & severityMask  == LOG_CRIT     { r = append(r, "CRIT") }
    if p & severityMask  == LOG_ERR      { r = append(r, "ERR") }
    if p & severityMask  == LOG_WARNING  { r = append(r, "WARNING") }
    if p & severityMask  == LOG_NOTICE   { r = append(r, "NOTICE") }
    if p & severityMask  == LOG_INFO     { r = append(r, "INFO") }
    if p & severityMask  == LOG_DEBUG    { r = append(r, "DEBUG") }

    if p & facilityMask  == LOG_KERN     { r = append(r, "KERN") }
    if p & facilityMask  == LOG_USER     { r = append(r, "USER") }
    if p & facilityMask  == LOG_MAIL     { r = append(r, "MAIL") }
    if p & facilityMask  == LOG_DAEMON   { r = append(r, "DAEMON") }
    if p & facilityMask  == LOG_AUTH     { r = append(r, "AUTH") }
    if p & facilityMask  == LOG_LPR      { r = append(r, "LPR") }
    if p & facilityMask  == LOG_NEWS     { r = append(r, "NEWS") }
    if p & facilityMask  == LOG_UUCP     { r = append(r, "UUCP") }
    if p & facilityMask  == LOG_CRON     { r = append(r, "CRON") }
    if p & facilityMask  == LOG_AUTHPRIV { r = append(r, "AUTHPRIV") }
    if p & facilityMask  == LOG_FTP      { r = append(r, "FTP") }

    if p & severityMask  == LOG_LOCAL0   { r = append(r, "LOCAL0") }
    if p & severityMask  == LOG_LOCAL1   { r = append(r, "LOCAL1") }
    if p & severityMask  == LOG_LOCAL2   { r = append(r, "LOCAL2") }
    if p & severityMask  == LOG_LOCAL3   { r = append(r, "LOCAL3") }
    if p & severityMask  == LOG_LOCAL4   { r = append(r, "LOCAL4") }
    if p & severityMask  == LOG_LOCAL5   { r = append(r, "LOCAL5") }
    if p & severityMask  == LOG_LOCAL6   { r = append(r, "LOCAL6") }
    if p & severityMask  == LOG_LOCAL7   { r = append(r, "LOCAL7") }

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
