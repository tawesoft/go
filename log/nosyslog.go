// +build windows

package log

func dial(network string, addr string, priority Priority, tag string) (Syslog, error) {
    return nil, fmt.Errorf("syslog is not available on this platform")
}
