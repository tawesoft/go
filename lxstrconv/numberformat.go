package lxstrconv

import (
    "sort"
)

// Unicode standard guarantees that this codepoint is never assigned
const runeNone = rune(0xFFFF)

// NumberFormat defines an interface for parsing numbers in a specific format
// (such as a decimal number in a specific locale, with support for a digit
// separator such as commas and a decimal point). Numbers are assumed to be
// in the normal base (e.g. base 10 for decimal) for that locale.
//
// Errors are either nil, strconv.ErrSyntax or strconv.ErrRange
type NumberFormat interface {
    ParseInt(string)   (int64, error)
    ParseFloat(string) (float64, error)
    
    // AcceptInt parses as much of an integer as possible. The second return
    // value is the number of bytes (not runes) successfully parsed. The error
    // value is always either nil or strconv.ErrRange.
    AcceptInt(string)   (int64, int, error)
    
    // AcceptFloat parses as much of a float as possible. The second return
    // value is the number of bytes (not runes) successfully parsed. The error
    // value is always either nil or strconv.ErrRange.
    AcceptFloat(string) (float64, int, error)
}

// repeatingRune returns any rune that appears more than once in a given string
func repeatingRune(s string) rune {
    // easy efficient algorithm: sort the string, then walk it and see if the
    // current rune repeats
    sl := []rune(s)
    sort.Slice(sl, func(i int, j int) bool { return sl[i] < sl[j] })
    
    current := runeNone
    
    for _, c := range sl {
        if c == current {
            return c
        }
        current = c
    }
    
    return runeNone
}



