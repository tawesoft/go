package email

import (
    "crypto/rand"
    "fmt"
    "io"
    "mime"
    "net/mail"
    "strings"
    "time"
)

// msgid generates a unique message ID for a sender
func msgid(from mail.Address) string {
    var rnd = make([]byte, 16)
    rand.Read(rnd)
    var tnow = time.Now().UTC().Format("20060102150405")
    var domain = "localhost"
    
    var atpos = strings.LastIndexByte(from.Address, '@')
    if atpos > 0 {
        domain = from.Address[atpos:] // includes leading @
    }
    
    return fmt.Sprintf("<%s.%x%s>", tnow, rnd, domain)
}

// strall returns true iff f(x) is true for each rune in string xs
func strall(
    xs string,
    f func(c rune) bool,
) bool {
    for _, x := range xs {
        // if c > unicode.MaxASCII {
        if !f(x) { return false }
    }

    return true
}

// optionalQEncode returns the original string if it is printable ASCII, or a UTF-8 encoded version otherwise.
func optionalQEncode(x string) string {
    var onlyPrintableAscii = func(c rune) bool {
        return (c >= 0x20) && (c <= 0x7e)
    }
    
    if strall(x, onlyPrintableAscii) {
        return x
    } else {
        return mime.QEncoding.Encode("utf-8", x)
    }
}

// lineBreaker wraps a writer forcing a line break every 76 characters (RFC 2045)
type lineBreaker struct {
    column int
    writer io.Writer
}

// Write writes p to the wrapped writer, forcing a line break every 76 characters (RFC 2045) across
// multiple calls to Write
func (lb lineBreaker) Write(p []byte) (n int, err error) {
    
    const LIMIT = 76 // RFC 2045
    var offset int
    var written int
    lb.column += len(p)
    
    for offset = 0; lb.column >= 76; lb.column -= LIMIT {
        n, err := lb.writer.Write(p[offset:offset + LIMIT])
        written += n
        if err != nil { return written, err }
        
        offset += LIMIT
        
        if offset != len(p) {
            n, err := io.WriteString(lb.writer, "\r\n")
            written += n
            if err != nil { return written, err }
        }
    }
    
    if offset != len(p) {
        n, err := lb.writer.Write(p[offset:])
        written += n
        if err != nil { return written, err }
    }
    
    return written, nil
}
