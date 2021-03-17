package humanizex

import (
    "testing"

    "golang.org/x/text/language"
)

func TestParseAcceptComponent(t *testing.T) {
    english := NewHumanizer(language.English).(*humanizer)
    danish  := NewHumanizer(language.Danish).(*humanizer)

    type test struct {
        humanizer *humanizer
        humanizerName string
        factors Factors
        value string
        unit Unit
        expected float64
        expectedRead int
        expectedError error
    }

    tests := []test{
        {english, "english", CommonFactors.Distance, "1.500 km", Unit{"m", "m"}, 1500, 7, nil},
        {danish,  "danish",  CommonFactors.Distance, "1,500 km", Unit{"m", "m"}, 1500, 7, nil},
    }

    for _, test := range tests {
        v, _, bytesRead, err := test.humanizer.acceptOne(test.value, test.factors)
        if (err != nil) || (v != test.expected) || (bytesRead != test.expectedRead) {
            t.Errorf("humanizer<%s>.acceptOne(%s, %q, factors): got %v, %d, %v but expected %v, %d, %v",
                test.humanizerName, test.value, test.unit,
                v, bytesRead, err,
                test.expected, test.expectedRead, test.expectedError)
        }
    }
}

func TestParseAcceptAllComponents(t *testing.T) {
    english := NewHumanizer(language.English).(*humanizer)
    danish  := NewHumanizer(language.Danish).(*humanizer)

    type test struct {
        humanizer *humanizer
        humanizerName string
        factors Factors
        value string
        unit Unit
        expected float64
        expectedRead int
        expectedError error
    }

    tests := []test{
        {english, "english", CommonFactors.Distance, "1.500 km", Unit{"m", "m"}, 1500, 8, nil},
        {danish,  "danish",  CommonFactors.Distance, "1,500 km", Unit{"m", "m"}, 1500, 8, nil},

        // unit m is optional
        {english, "english", CommonFactors.Distance, "1.500 k", Unit{"m", "m"}, 1500, 7, nil},
        {english, "english", CommonFactors.Distance, "1.500 km Trailing", Unit{"m", "m"}, 1500, 9, nil},

        {danish,  "english", CommonFactors.Time,     "30 min 1 s", CommonUnits.Second, 1 + (30 * 60), 10, nil},
    }

    for _, test := range tests {
        v, bytesRead, err := test.humanizer.Accept(test.value, test.unit, test.factors)
        if (err != nil) || (v != test.expected) || (bytesRead != test.expectedRead) {
            t.Errorf("humanizer<%s>.accept(%q, %q, factors): got %v, %d, %v but expected %v, %d, %v",
                test.humanizerName, test.value, test.unit,
                v, bytesRead, err,
                test.expected, test.expectedRead, test.expectedError)
        }
    }
}
