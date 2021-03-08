package email

import (
    "crypto/rand"
    "fmt"
    "io"
    "mime"
    "strings"
    "time"
)

// NewMessageID generates a cryptographically unique RFC 2822 3.6.4 Message ID
// (not including angle brackets).
func NewMessageID(host string) string {
    var rnd = make([]byte, 16)
    rand.Read(rnd)
    var tnow = time.Now().UTC().Format("20060102150405")
    return fmt.Sprintf("%s.%x@%s", tnow, rnd, host)
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

// folding whitespace for a key: value header
var fwsWrapErr = fmt.Errorf("header value component too long")
var fwsNoneErr = fmt.Errorf("header value too long")

func fwsWrap(line string, keyLen int, maxLine int) (string, error) {
    //fmt.Printf("\n\n=== wrap (%q, %d, %d) ===\n", line, keyLen, maxLine)

    // thorny but tested code. Complicated by the fact that the first line has
    // a header "Key: " prefix, and subsequent lines have a single space indent.

    keyLen += 2 // ": "
    if len(line) + keyLen <= maxLine { return line, nil }

    result := make([]string, 0)

    offset := 0
    start, idx := 0, 0
    remainingBudget := maxLine - keyLen // special case for first line
    remainingBudget++                   // trailing space that we break on
    //lineBudget := remainingBudget
    lineBudgetIncludingIndent := remainingBudget // normal for first line

    for {
        // look at next space.
        // we could support break on "<" and ">" but it won't get us that much
        idx = strings.IndexByte(line[offset:], ' ')
        //fmt.Printf("consider %q... (space at %d)\n", line[offset:], idx)

        if (idx >= 0) && (idx <= remainingBudget) {
            //fmt.Printf("in budget, continue... (budget %d => %d)\n", remainingBudget, remainingBudget- (idx + 1))
            remainingBudget -= idx + 1
            offset += idx + 1
            continue
        } else if (idx < 0) && (len(line[start:]) <= lineBudgetIncludingIndent) {
            //fmt.Printf("finished (no further spaces, but line of length %d in budget %d): %q\n", len(line[start:]), lineBudget, line[start:])
            segment := line[start:]
            if len(segment) > lineBudgetIncludingIndent { return "", fwsWrapErr }
            result = append(result, segment)
            break
        } else if (idx < 0) {
            //fmt.Printf("finished (no further spaces, but line %q of length %d not in budget %d)\n", line[start:], len(line[start:]), lineBudget)

            if offset < 1 {
                // first element too long to wrap
                return "", fwsWrapErr
            }

            a := line[start:offset-1]
            b := line[offset:]
            if len(a) > lineBudgetIncludingIndent { return "", fwsWrapErr }
            if len(b) > maxLine - 1 { return "", fwsWrapErr }
            result = append(result, line[start:offset-1], line[offset:])
            break
        } else {
            if offset < 1 {
                // first element too long to wrap
                return "", fwsWrapErr
            }

            //fmt.Printf("not finished, but out of budget %d...\n", remainingBudget)
            segment :=  line[start:offset-1]
            if len(segment) > lineBudgetIncludingIndent { return "", fwsWrapErr }
            result = append(result, segment)
            remainingBudget = maxLine - 1 + 1 // -1 for leading whitespace, +1 for trailing space that we break on
            //lineBudget = remainingBudget
            lineBudgetIncludingIndent = remainingBudget
            start = offset
            continue
        }
    }

    //fmt.Printf("\n\n\n")
    return strings.Join(result, "\r\n "), nil
}

func fwsNone(line string, keyLen int, maxLine int) (string, error) {
    keyLen += 2 // ": "
    if len(line) + keyLen <= maxLine { return line, nil }
    return "", fwsNoneErr
}
