package humanizex

import (
    "testing"

    "golang.org/x/text/language"
)

func TestFormatParts(t *testing.T) {

    factors := Factors{
        Factors:    []Factor{
            { 0.10, Unit{"small",  "small(ascii)"},  FactorModeUnitPrefix},
            { 1.00, Unit{"", ""}, FactorModeIdentity},
            {10.00, Unit{"big",    "big(ascii)"},    FactorModeUnitPrefix},
        },
    }

    factors2 := Factors{
        Components: 2,
        Factors:    []Factor{
            { 0.10, Unit{"small",  "small(ascii)"},  FactorModeUnitPrefix},
            { 1.00, Unit{"", ""}, FactorModeIdentity},
            {10.00, Unit{"big",    "big(ascii)"},    FactorModeUnitPrefix},
        },
    }

    type test struct {
        factors Factors
        value float64
        unit Unit
        expectedParts []Part
    }

    tests := []test{
        {factors, 0.01, Unit{"unit", "unit(ascii)"}, []Part{{0.1, Unit{"smallunit", "small(ascii)unit(ascii)"}}}},
        {factors, 0.15, Unit{"unit", "unit(ascii)"}, []Part{{1.5, Unit{"smallunit", "small(ascii)unit(ascii)"}}}},
        {factors, 1.00, Unit{"unit", "unit(ascii)"}, []Part{{1.0, Unit{"unit", "unit(ascii)"}}}},
        {factors, 1.50, Unit{"unit", "unit(ascii)"}, []Part{{1.5, Unit{"unit", "unit(ascii)"}}}},
        {factors, 15.0, Unit{"unit", "unit(ascii)"}, []Part{{1.5, Unit{"bigunit", "big(ascii)unit(ascii)"}}}},

        {factors2, 15.0, Unit{"unit", "unit(ascii)"}, []Part{
            {1.0, Unit{"big", "big(ascii)"}},
            {5.0, Unit{"unit", "unit(ascii)"}},
        }},

        {CommonFactors.Time, 0, Unit{"s", "s"}, []Part{
            {0.0, Unit{"s", "s"}},
        }},

        {CommonFactors.Time, 60 * 2.5, Unit{"s", "s"}, []Part{
            { 2.0, Unit{"min", "min"}},
            {30.0, Unit{"s", "s"}},
        }},

        // skips zero units
        {CommonFactors.Time, 60 * 2.0, Unit{"s", "s"}, []Part{
            { 2.0, Unit{"min", "min"}},
        }},
        {CommonFactors.Time, 1 + (60 * 60 * 2.0), Unit{"s", "s"}, []Part{
            { 2.0, Unit{"h", "h"}},
            { 1.0, Unit{"s", "s"}},
        }},
    }

    const epsilon = 0.01

    for _, test := range tests {
        parts := FormatParts(test.value, test.unit, test.factors)
        if !partsEqual(parts, test.expectedParts, epsilon) {
            t.Errorf("FormatParts(%f, %q, factors): got %v but expected %v",
                test.value, test.unit, parts, test.expectedParts)
        }
    }
}

func TestFormat(t *testing.T) {
    english := NewHumanizer(language.English)
    danish  := NewHumanizer(language.Danish)

    type test struct {
        humanizer Humanizer
        humanizerName string
        factors Factors
        value float64
        unit Unit
        expected string
    }

    tests := []test{
        {english, "english", CommonFactors.Distance, 1500 * 1000, Unit{"m", "m"}, "1,500 km"},
        {danish,  "danish", CommonFactors.Distance, 1500 * 1000, Unit{"m", "m"}, "1.500 km"},
    }

    for _, test := range tests {
        str := test.humanizer.Format(test.value, test.unit, test.factors)
        if str.Utf8 != test.expected {
            t.Errorf("humanizer<%s>.FormatParts(%f, %q, factors): got %v but expected %v",
                test.humanizerName, test.value, test.unit, str, test.expected)
        }
    }
}
