package humanizex

import (
    "testing"
)

func TestFactorBracket(t *testing.T) {
    factors := Factors{
        Factors:    []Factor{
            { 0.10, Unit{"small", "small"},   0},
            { 1.00, Unit{"normal", "normal"}, 0},
            {10.00, Unit{"big",    "big"},    0},
        },
    }

    type test struct {
        value float64
        expectedFactorIndex int
    }

    tests := []test{
        { 0.01,  0},
        { 0.10,  0},
        { 0.50,  0},
        { 1.00,  1},
        { 1.50,  1},
        {10.00,  2},
        {11.00,  2},
    }

    for _, test := range tests {
        idx := factors.bracket(test.value)
        if idx != test.expectedFactorIndex {
            t.Errorf("factors.bracket(%f): got idx %d but expected %d",
                test.value, idx, test.expectedFactorIndex)
        }
    }
}
