package humanize

type Formatter struct {
    DecimalPlaces int
    DecimalSeparator   string // e.g. 1<--DecimalSeparator-->23
    UnitSeparator      string // e.g. 1.2<--UnitSeparator-->M
}

var DefaultDotFormatter       = Formatter{2, ".", " "}
var DefaultCommaFormatter     = Formatter{2, ",", " "}
var DefaultUnicodeFormatter   = Formatter{2, "⎖", " "}
var DefaultMiddleDotFormatter = Formatter{2, "·", " "}
var DefaultBraileFormatter    = Formatter{2, "⠨", " "}

var DefaultFormatter = DefaultDotFormatter

