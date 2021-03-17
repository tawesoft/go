package humanizex

import (
    "fmt"
    "strings"
    "unicode"
    "unicode/utf8"
)

// Catagory Zs in http://www.unicode.org/Public/UNIDATA/UnicodeData.txt
var whitespaceRunes = string([]rune{
    '\u0020',
    '\u00A0',
    '\u1680',
    '\u2000',
    '\u2001',
    '\u2002',
    '\u2003',
    '\u2004',
    '\u2005',
    '\u2006',
    '\u2007',
    '\u2008',
    '\u2009',
    '\u200A',
    '\u202F',
    '\u205F',
    '\u3000',
})

func skipSpace(str string, bytesRead int) int {
    for _, r := range str[bytesRead:] {
        if !utf8.ValidRune(r) { break }
        if !unicode.IsSpace(r) { break }
        bytesRead += utf8.RuneLen(r)
    }

    return bytesRead
}

func (h *humanizer) Parse(str string, unit Unit, factors Factors) (float64, error) {
    v, bytesRead, err := h.Accept(str, unit, factors)
    if err != nil { return 0, err }

    if len(str) != bytesRead {
        return 0, fmt.Errorf("error parsing %q: unexpected trailing content after byte %d", str, bytesRead)
    }

    return v, err
}

func (h *humanizer) Accept(str string, unit Unit, factors Factors) (float64, int, error) {
    var v float64
    var bytesRead int
    var lastFactor Factor

    components := factors.Components
    if components < 1 { components = 1 }

    for i := 0; i < components; i++ {
        c, f, r, err := h.acceptOne(str[bytesRead:], factors)
        if (err != nil) && (i > 0) { break }
        if err != nil { return 0, 0, err }

        v += c
        bytesRead += r
        lastFactor = f
        if bytesRead == len(str) { break }
    }

    if len(str) == bytesRead { return v, bytesRead, nil }
    if lastFactor.Mode & FactorModeReplace == FactorModeReplace { return v, bytesRead, nil }

    // optionally accept the final unit
    bytesRead = skipSpace(str, bytesRead)
    if len(str) == bytesRead { return v, bytesRead, nil }
    remaining := str[bytesRead:]

    if strings.HasPrefix(remaining, unit.Utf8) {
        bytesRead += len(unit.Utf8)
    } else if strings.HasPrefix(remaining, unit.Ascii) {
        bytesRead += len(unit.Ascii)
    }

    // optionally accept trailing space
    if len(str) != bytesRead {
        bytesRead = skipSpace(str, bytesRead)
    }

    return v, bytesRead, nil
}

func (h *humanizer) acceptOne(str string, factors Factors) (float64, Factor, int, error) {
    v, bytesRead, err := h.NF.AcceptFloat(str)

    if err != nil {
        // only ever strconv.ErrRange
        return 0, Factor{}, 0, fmt.Errorf("error parsing number component of %q: %v", str, err)
    } else if bytesRead == 0 {
        return 0, Factor{}, 0, nil
    }

    bytesRead = skipSpace(str, bytesRead)
    if len(str) == bytesRead { return v, Factor{}, bytesRead, nil }

    remaining := str[bytesRead:]
    multiplier := 1.0
    unitLen := 0 // logical length in runes
    unitLenBytes := 0 // actual length in bytes
    factor := Factor{}
    for _, f := range factors.Factors {
    // don't break early, use the longest prefix just in case
    // we have a hypothetical unit prefix with a common "unit
    // prefix prefix" e.g. "x" and "xx"
        if strings.HasPrefix(remaining, f.Unit.Utf8) {
            if unitLen < len(f.Unit.Ascii) {
                multiplier = f.Magnitude
                unitLen = len(f.Unit.Ascii)
                unitLenBytes = len(f.Unit.Utf8)
                factor = f
            }
        } else if strings.HasPrefix(remaining, f.Unit.Ascii) {
            if unitLen < len(f.Unit.Ascii) {
                multiplier = f.Magnitude
                unitLen = len(f.Unit.Ascii)
                unitLenBytes = len(f.Unit.Ascii)
                factor = f
            }
        }
    }

    return v * multiplier, factor, bytesRead + unitLenBytes, nil
}
