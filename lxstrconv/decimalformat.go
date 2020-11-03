package lxstrconv

import (
    "math"
    "strconv"
    "unicode"
    "unicode/utf8"
    
    "golang.org/x/text/language"
    "golang.org/x/text/message"
    "golang.org/x/text/number"
)

// acceptRune returns the length of r in bytes if r is the first rune in s,
// otherwise returns zero.
func acceptRune(r rune, s string) int {
    if f, ok := firstRune(s); ok && (f == r) {
        return utf8.RuneLen(r)
    } else {
        return 0
    }
}

// firstRune returns the first rune in a string and true, or (_, false).
func firstRune(s string) (rune, bool) {
    for _, c := range s {
        return c, true
    }
    return runeNone, false
}

// guessDecimalGroupSeparator guesses, for a printer in a given locale,
// the group separator rune in a decimal number system e.g. comma for British.
func guessDecimalGroupSeparator(p *message.Printer) rune {
    // heuristic: any rune that appears at least twice is probably a comma
    s := p.Sprint(number.Decimal(1234567890))
    return repeatingRune(s)
}

// guessDecimalPointSeparator guesses, for a printer in a given locale,
// the decimal point rune in a decimal number system, e.g. period for British.
func guessDecimalPoint(p *message.Printer) rune {
    // heuristic: any rune that is common to both these strings is probably a
    // decimal point. Concat the strings and find any repeated rune.
    s1 := p.Sprint(number.Decimal(1.23))
    s2 := p.Sprint(number.Decimal(4.56))
    s := s1 + s2
    return repeatingRune(s)
}

// guessDecimalDigits guesses, for a printer in a given locale, the digits
// representing the values 0 to 9.
func guessDecimalDigits(p *message.Printer, out *[10]rune) {
    for i := 0; i < 10; i++ {
        s := []rune(p.Sprint(number.Decimal(i)))
        if len(s) == 1 {
            out[i] = s[0]
        } else {
            out[i] = runeNone
        }
    }
}

// decimalFormat defines how a decimal (base-10) number should be parsed for a
// given locale. Note that the behaviour is undefined for locales that have
// non-base-10 number systems.
//
// This structure is currently internal until we have more confidence it is
// correct for all languages with decimal number systems.
type decimalFormat struct {
    // GroupSeparator is a digits separator such as commas for thousands. In
    // addition to any separator defined here, a parser will ignore whitespace.
    GroupSeparator rune
    
    // Point is separator between the integer and fractional part of
    // a decimal number.
    Point rune
    
    // Digits are an ascending list of digit runes
    Digits [10]rune
}

func (f decimalFormat) ParseInt(s string) (int64, error) {
    if len(s) == 0 { return 0, strconv.ErrSyntax }
    
    value, length, err := f.AcceptInt(s)
    
    if err != nil { return 0, err }
    if len(s) != length { return 0, strconv.ErrSyntax }
    
    return value, nil
}

func (f decimalFormat) ParseFloat(s string) (float64, error) {
    if len(s) == 0 { return 0, strconv.ErrSyntax }
    
    value, length, err := f.AcceptFloat(s)
    
    if err != nil { return 0, err }
    if len(s) != length { return 0, strconv.ErrSyntax }
    
    return value, nil
}

// NewDecimalFormat constructs, for a given locale, a NumberFormat that
// defines how a decimal (base-10) number should be parsed. Note that the
// behaviour is undefined for locales that have non-base-10 number systems.
func NewDecimalFormat(tag language.Tag) NumberFormat {
    
    // Unfortunately, I couldn't find any exported symbols in /x/text that
    // gives this information directly (as would be ideal). Therefore this
    // function works by printing numbers in the current locale and using
    // heuristics to guess the correct separators.
    
    p := message.NewPrinter(tag)
    
    format := decimalFormat{
        GroupSeparator: guessDecimalGroupSeparator(p),
        Point:          guessDecimalPoint(p),
    }
    
    guessDecimalDigits(p, &format.Digits)
    
    return format
}

// returns (0-9, true) for a decimal digit in any language, or (_, false)
func decimalRuneToInt(d rune, digits *[10]rune) (int, bool) {
    for i := 0; i < 10; i++ {
        if d == digits[i] { return i, true }
    }
    return 0, false
}

// AcceptInteger parses as much of an integer number as possible. It returns a
// 2 tuple: the value of the parsed integer, and the length of the characters
// successfully parsed. For example, for some locales, the string "1,000X"
// returns (1000, 5) and the string "foo" returns (0, 0).
//
// Err is always nil, strconv.ErrRange or strconv.ErrSyntax
func (f decimalFormat) AcceptInt(s string) (value int64, length int, err error) {

    if len(s) == 0 { return 0, 0, nil }
    
    if s[0] == '-' {
        // TODO better negative check e.g. "(1)" for "-1"
        v, l, _ := f.AcceptUint(s[1:])
        // TODO bounds check
        if l > 0 {
            return int64(v) * -1, l + 1, nil
        } else {
            return 0, 0, nil
        }
    }
    
    // TODO bounds check
    v, l, err := f.AcceptUint(s)
    return int64(v), l, nil
}

// AcceptUint: see AcceptInt
func (f decimalFormat) AcceptUint(s string) (value uint64, length int, err error) {
    var accu uint64
    
    for i, c := range s {
        if c == f.GroupSeparator {
            // pass
        } else if unicode.IsSpace(c) {
            // pass
        } else if d, ok := decimalRuneToInt(c, &f.Digits); ok {
            accu *= 10
            accu += uint64(d)
            // TODO bounds check
        } else {
            // TODO this count is runes but should be bytes!
            return accu, i, nil
        }
    }
    
    return accu, len(s), nil
}

// AcceptFloat parses as much of a floating point number as possible. It returns
// a 2 tuple: the value of the parsed float, and the length of the characters
// successfully parsed. For example, for some locales, the string "1.23X"
// returns (1.23, 4) and the string "foo" returns (0.0, 0).
//
// Err is always nil, strconv.ErrRange or strconv.ErrSyntax
func (f decimalFormat) AcceptFloat(s string) (value float64, length int, err error) {
    var left, right int64
    var leftLen, rightLen, pointLen int
    var fLeft, fRight float64
    
    // accept leading decimal point
    if first, ok := firstRune(s); ok && first != f.Point {
        left, leftLen, err = f.AcceptInt(s)
        // TODO check err (Currently always nil)
        if leftLen == 0 { return 0, 0, nil }
        fLeft = float64(left)
    }
    
    pointLen = acceptRune(f.Point, s[leftLen:])
    
    if pointLen > 0 && (s[leftLen +pointLen] != '-') {
        right, rightLen, err = f.AcceptInt(s[leftLen +pointLen:])
        // TODO check err (currently always nil)
    }
    
    if right > 0.0 {
        fRight = float64(right)
        places := float64(1.0 + math.Floor(math.Log10(fRight)))
        fRight *= math.Pow(0.1, places)
        fRight = math.Copysign(fRight, fLeft)
    }
    
    value = fLeft + fRight
    length = leftLen + pointLen + rightLen
    
    return value, length, nil
}
