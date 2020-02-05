package humanize

import (
    "fmt"
)

// AcceptInteger parses as much of an integer number as possible. All separators, such as comma, period, underscore,
// space, and Unicode thin space, are ignored. It returns a 2 tuple: the value of the parsed integer, and the length
// of the characters successfully parsed. For example, the string "1,000 M" returns (1000, 6) and the string "foo"
// returns (0, 0).
func AcceptInteger(s string) (value int64, length int) {
    var accu int64
    
    if len(s) == 0 { return 0, 0 }
    
    if s[0] == '-' {
        var v, l = AcceptInteger(s[1:])
        if l > 0 {
            return v * -1, l + 1
        } else {
            return 0, 0
        }
    }
    
    for i, c := range s {
        switch c {
            case '\u2009': // Unicode thin space
            case ' ':
            case ',':
            case '.':
            case '_':
            default:
                if c >= '0' && c <= '9' {
                    accu *= 10
                    accu += int64(c - '0')
                } else {
                    return accu, i
                }
        }
    }
    
    return accu, len(s)
}

// AcceptIntegerUnitPrefix looks for any SI or IEC prefix (like "M" for "Mega") greater in magnitude than 1 (that is, it
// ignores prefixes like "μ" for "micro"). If found, it returns a 2 tuple: the quantity represented by the prefix (in
// the case of "M", this is 1,000,000), and the length of the characters parsed (in the case of "M", this is 1). If not
// found, length is zero and the quantity is 1.
func AcceptIntegerUnitPrefix(s string) (quantity int64, length int) {
    if (len(s) >= 2) && (s[1] == 'i') {
        switch s[0] {
            case 'K': return PrefixKibi, 2
            case 'M': return PrefixMebi, 2
            case 'G': return PrefixGibi, 2
            case 'T': return PrefixTebi, 2
            case 'P': return PrefixPebi, 2
            case 'E': return PrefixExbi, 2
        }
    }
    
    if len(s) >= 1 {
        switch s[0] {
            case 'k': return PrefixKilo,  1
            case 'M': return PrefixMega,  1
            case 'G': return PrefixGiga,  1
            case 'T': return PrefixTera,  1
            case 'P': return PrefixPeta,  1
            case 'E': return PrefixExa,   1
        }
    }
    
    return 1, 0
}

// ParseIntegerPartial parses any integer value (such as "1,000"), optionally followed by an SI or IEC unit prefix
// (such as "k"). For example the number "1k" is parsed as the integer 1,000. If found, it returns a 2 tuple: the
// normalised quantity, and the length of the characters parsed (e.g. in the case of "1k" this is 2). If not found,
// length is zero and the quantity is zero.
func ParseIntegerPartial(text string) (int64, int) {
    
    // e.g. 1,000
    var value, ilength = AcceptInteger(text)
    if ilength == 0 { return 0, 0 }
    
    // e.g. K, Ki
    var prefix, plength = AcceptIntegerUnitPrefix(text[ilength:])
    
    return prefix * value, ilength + plength
}

// ParseInteger works like ParseIntegerPartial, except the entire string must be successfully parsed: otherwise an
// error is returned.
func ParseInteger(text string) (int64, error) {
    
    var value, length = ParseIntegerPartial(text)
    
    if length == len(text) {
        return value, nil
    } else {
        return 0, &ParseError{text, fmt.Errorf("unrecognised trailing content while parsing integer: expected end of string at character %d", length)}
    }
}

// FormatInteger returns a value such as "1.2M"
func FormatInteger(formatter *Formatter, value int64) string {
    return FormatFloat(formatter, float64(value))
}

// FormatIntegerIEC returns a value such as "1.2Mi"
func FormatIntegerIEC(formatter *Formatter, value int64) string {
    return FormatFloatIEC(formatter, float64(value))
}
