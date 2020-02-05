package humanize

import (
    "fmt"
)

// ParseBytesPartial parses any bytes value (such as "50"), optionally followed by an SI or IEC unit prefix
// (such as "M" or "Mi"), optionally followed by "B" for bytes. For example the value "1MiB" is parsed as the integer
// 1,024. If found, it returns a 2 tuple: the normalised quantity, and the length of the characters parsed (e.g. in
// the case of "1MiB" this is 4). If not found, length is zero and the quantity is zero. Negative byte values are not
// allowed: a value and length of zero will be returned.
func ParseBytesPartial(text string) (int64, int) {
    
    var value, length = ParseIntegerPartial(text)
    
    if text[length:] == "B" {
        return value, length + 1
    } else {
        return value, length
    }
}

// ParseBytes works like ParseBytesPartial, except the entire string must be successfully parsed: otherwise an
// error is returned.
func ParseBytes(text string) (int64, error) {
    
    var value, length = ParseBytesPartial(text)
    
    if length == len(text) {
        return value, nil
    } else {
        return 0, &ParseError{text, fmt.Errorf("unrecognised trailing content while parsing bytes: expected end of string at character %d", length)}
    }
}

// FormatBytesIEC returns a string such as "1.2MiB"
func FormatBytesIEC(formatter *Formatter, value int64) string {
    return FormatIntegerIEC(formatter, value) + "B"
}

// FormatBytesSI returns a string such as "1.2MB"
func FormatBytesSI(formatter *Formatter, value int64) string {
    return FormatInteger(formatter, value) + "B"
}
