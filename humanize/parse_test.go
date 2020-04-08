package humanize

import (
    "testing"
)

func TestAcceptInteger(t *testing.T) {
    type test struct {
        text string
        expectedLength int
        expectedValue int64
    }
    
    var tests = []test{
        {"",          0,    0},
        {"foo",       0,    0},
        {"1",         1,    1},
        {"2",         1,    2},
        {"12",        2,   12},
        {"123",       3,  123},
        {"321",       3,  321},
        {"1,000,000", 9,  1e6},
        {"123foo",    3,  123},
        {"321foo",    3,  321},
        {"-1",        2,   -1},
        {"-",         0,    0},
        {"--",        0,    0},
        {"--1",       3,    1},
        {"-foo",      0,    0},
        {"-123",      4, -123},
        {"1-23",      1,    1},
    }
    
    for _, i := range tests {
        var result, length = AcceptInteger(nil, i.text)
        if result != i.expectedValue || length != i.expectedLength {
            t.Errorf("AcceptInteger(%s): got (%d, %d) but expected (%d, %d)", i.text, result, length, i.expectedValue, i.expectedLength)
        }
    }
}

func TestAcceptIntegerUnitPrefix(t *testing.T) {
    type test struct {
        text string
        expectedLength int
        expectedValue int64
    }
    
    var tests = []test{
        {"",        0, 1},
        {"M",       1, 1e6},
        {"Mi",      2, 1024*1024},
        {"_foo",    0, 1},
        {"M_foo",   1, 1e6},
        {"Mi_foo",  2, 1024*1024},
    }
    
    for _, i := range tests {
        var result, length = AcceptIntegerUnitPrefix(i.text)
        if result != i.expectedValue || length != i.expectedLength {
            t.Errorf("AcceptIntegerPrefix(%s): got (%d, %d) but expected (%d, %d)", i.text, result, length, i.expectedValue, i.expectedLength)
        }
    }
}

func TestParseBytes(t *testing.T) {
    
    type test struct {
        text string
        expected int64
        err bool
    }
    
    var tests = []test{
        {"B",                    0, false},
        {"-B",                   0, true},
        {"1",                    1, false},
        {"1B",                   1, false},
        {"1 B",                  1, false},
        {"128B",               128, false},
        {"128 B",              128, false},
        {"128k",           128_000, false},
        {"128kB",          128_000, false},
        {"128 kB",         128_000, false},
        {"128Ki",         128*1024, false},
        {"128 Ki",        128*1024, false},
        {"128KiB",        128*1024, false},
        {"128 KiB",       128*1024, false},
        {"1,000,000 B",        1e6, false},
        {"1GiB",    1024*1024*1024, false},
        {"-1GiB",  -1024*1024*1024, false},
    }
    
    for _, i := range tests {
        var result, err = ParseBytes(nil, i.text)
        if err != nil && i.err == false {
            t.Errorf("ParseBytes(%s): expected %d but got error: %v", i.text, i.expected, err)
        } else if err != nil && i.err == true {
            // pass
        } else if result != i.expected {
            t.Errorf("ParseBytes(%s): got %d but expected %d", i.text, result, i.expected)
        }
    }
}
