package humanizex

import (
    "math"
    "time"

    "golang.org/x/text/language"
    "golang.org/x/text/message"
    "tawesoft.co.uk/go/lxstrconv"
)

// String holds an Utf8-encoded and an Ascii-compatible encoding of a string.
type String struct {
    // Utf8 is the native Utf8-encoded Unicode representation
    Utf8   string

    // Ascii is an alternative version accepted for non-Unicode inputs (such
    // as when a user does not know how to enter µ on their keyboard (on mine,
    // it's Right-Alt+m)) or for non-Unicode output (such as legacy systems).
    Ascii  string
}

// Unit describes some quantity e.g. "m" for length in metres, "k" for the SI
// unit prefix kilo, "km" for a kilometre.
type Unit String

// Cat concatenates two units (u + v) and returns the result.
func (u Unit) Cat(v Unit) Unit {
    return Unit{
        u.Utf8 + v.Utf8,
        u.Ascii + v.Ascii,
    }
}

// Part describes some component of a formatting result e.g. 1.5 km or 1 hour
type Part struct {
    Magnitude float64
    Unit      Unit
}

func partEqual(a Part, b Part, epsilon float64) bool {
    if a.Unit != b.Unit { return false }
    return math.Abs(a.Magnitude - b.Magnitude) < epsilon
}

func partsEqual(a []Part, b []Part, epsilon float64) bool {
    if len(a) != len(b) { return false }

    for i := 0; i < len(a); i++ {
        if !partEqual(a[i], b[i], epsilon) { return false }
    }

    return true
}

// Humanizer implements a locale-aware way to parse and format humanized
// quantities.
type Humanizer interface {
    // Format is a general purpose locale-aware way to format any quantity
    // with a defined set of factors. The unit argument is the base unit
    // e.g. s for seconds, m for meters, B for bytes.
    Format(value float64, unit Unit, factors Factors) String

    FormatNumber(number float64) String             // e.g. 12 k
    FormatDistance(meters float64) String           // e.g. 10 µm, 10 km
    FormatDuration(duration time.Duration) string   // e.g. 1 h 50 min
    FormatSeconds(seconds float64) string           // e.g. 1 h 50 min
    FormatBytesJEDEC(bytes int64) string            // e.g. 12 KB, 5 MB
    FormatBytesIEC(bytes int64) string              // e.g. 12 kB, 5 MB
    FormatBytesSI(bytes int64) string               // e.g. 12 KiB, 5 MiB

    // Accept is a general purpose locale-aware way to parse any quantity
    // with a defined set of factors from the start of the string str. The
    // provided unit is optional and is accepted if it appears in str.
    //
    // Accept returns the value, the number of bytes successfully parsed (which
    // may be zero), or an error.
    Accept(str string, unit Unit, factors Factors) (float64, int, error)

    // Parse is a general purpose locale-aware way to parse any quantity
    // with a defined set of factors.  The provided unit is optional and is
    // accepted if it appears in str.
    Parse(str string, unit Unit, factors Factors) (float64, error)

    ParseDuration(str string) (time.Duration, error)
    ParseBytesJEDEC(str string) (int64, error)
    ParseBytesIEC(str string) (int64, error)
    ParseBytesSI(str string) (int64, error)
}

type humanizer struct {
    Tag language.Tag
    NF lxstrconv.NumberFormat
    Printer *message.Printer
}

func (h *humanizer) FormatDistance(meters float64) String {
    return h.Format(meters, CommonUnits.Meter, CommonFactors.Distance)
}

func (h *humanizer) FormatBytesJEDEC(bytes int64) string {
    return h.Format(float64(bytes), CommonUnits.Byte, CommonFactors.JEDEC).Utf8
}

func (h *humanizer) FormatBytesIEC(bytes int64) string {
    return h.Format(float64(bytes), CommonUnits.Byte, CommonFactors.IEC).Utf8
}

func (h *humanizer) FormatBytesSI(bytes int64) string {
    return h.Format(float64(bytes), CommonUnits.Byte, CommonFactors.SI).Utf8
}

func (h *humanizer) FormatDuration(duration time.Duration) string {
    return h.Format(duration.Seconds(), CommonUnits.Second, CommonFactors.Time).Utf8
}

func (h *humanizer) FormatSeconds(seconds float64) string {
    return h.Format(seconds, CommonUnits.Second, CommonFactors.Time).Utf8
}

func (h *humanizer) FormatNumber(number float64) String {
    return h.Format(number, CommonUnits.None, CommonFactors.SI)
}

func (h *humanizer) ParseDuration(str string) (time.Duration, error) {
    v, err := h.Parse(str, CommonUnits.Second, CommonFactors.Time)
    return time.Second * time.Duration(v), err
}

func (h *humanizer) ParseBytesJEDEC(str string) (int64, error) {
    v, err := h.Parse(str, CommonUnits.Byte, CommonFactors.JEDEC)
    return int64(v), err
}

func (h *humanizer) ParseBytesIEC(str string) (int64, error) {
    v, err := h.Parse(str, CommonUnits.Byte, CommonFactors.IEC)
    return int64(v), err
}

func (h *humanizer) ParseBytesSI(str string) (int64, error) {
    v, err := h.Parse(str, CommonUnits.Byte, CommonFactors.SI)
    return int64(v), err
}

// NewHumanizer initialises a human language number encoder/decoder for the
// given tag (representing a specific language or locale).
//
// The language.Tag is usually a named language from golang.org/x/text/language
// e.g. language.English and controls how numbers are written e.g. comma
// placement, decimal point, digits.
func NewHumanizer(tag language.Tag, options ... interface{}) Humanizer {
    return &humanizer{
        Tag:     tag,
        NF:      lxstrconv.NewDecimalFormat(tag),
        Printer: message.NewPrinter(tag),
    }
}
