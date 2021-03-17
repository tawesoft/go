package humanizex

import (
    "math"
    "strings"
)

// FormatParts is a general purpose locale-aware way to format any quantity
// with a defined set of factors into a list of parts. The unit argument is the
// base unit e.g. s for seconds, m for meters, B for bytes. In the simple case
// a list of only one Part is returned, e.g. a Part representing 1.5 km. In
// other cases, there may be multiple parts e.g. the two parts "1 h" and
// "30 min" making up the time "1 h 30 min". The number of parts returned is
// never more than that defined by the factors argument's Components field
// (but may be fewer).
func FormatParts(n float64, unit Unit, factors Factors) []Part {

    numComponents := factors.Components
    if numComponents == 0 { numComponents = 1 }

    parts := make([]Part, 0, numComponents)

    for i := 0; i < numComponents; i++ {
        factorIdx := factors.bracket(n)
        factor := factors.Factors[factorIdx]

        part := Part{
            Magnitude: n / factor.Magnitude,
            Unit:      factor.Unit,
        }

        // skip if not the only non-zero component, but zero magnitude
        const epsilon = 0.01
        if (part.Magnitude < epsilon) && (len(parts) > 0) {
            continue
        }

        if i == numComponents - 1 { // last?
            if factor.Mode & FactorModeUnitPrefix == FactorModeUnitPrefix {
                part.Unit = part.Unit.Cat(unit)
            } else if factor.Mode & FactorModeReplace == FactorModeReplace {
                part.Unit = factor.Unit
            } else {
                part.Unit = unit
            }
        } else {
            int, frac := math.Modf(part.Magnitude)
            part.Magnitude = int
            n = frac * factor.Magnitude
        }

        parts = append(parts, part)
    }

    return parts
}

func (h *humanizer) Format(n float64, unit Unit, factors Factors) String {
    resultUtf8  := make([]string, 0, factors.Components)
    resultAscii := make([]string, 0, factors.Components)
    parts := FormatParts(n, unit, factors)

    for _, part := range parts {

        places := 0

        if part.Magnitude < 10.0 {
            places = 1
        }

        if part.Magnitude < 1.0 {
            places = 2
        }

        _, frac := math.Modf(part.Magnitude)
        if math.Abs(frac) < 0.01 {
            places = 0
        }

        str := h.Printer.Sprintf("%.*f", places, part.Magnitude)
        resultUtf8 = append(resultUtf8, str, part.Unit.Utf8)
        resultAscii = append(resultAscii, str, part.Unit.Ascii)
    }

    return String{
        strings.Join(resultUtf8, " "),
        strings.Join(resultAscii, " "),
    }
}

