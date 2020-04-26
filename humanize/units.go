package humanize

/*
SI Metric Prefixes
deci  d 10^−1
centi c 10^-2
milli m 10^-3
micro μ 10^−6
nano  n 10^−9
pico  p 10^−12
femto f 10^−15
*/

const (
    PrefixDeci  float64 = 1.0/1e1
    PrefixCenti float64 = 1.0/1e2
    PrefixMilli float64 = 1.0/1e3
    PrefixMicro float64 = 1.0/1e6
    PrefixNano  float64 = 1.0/1e9
    PrefixPico  float64 = 1.0/1e12
    PrefixFemto float64 = 1.0/1e15
)

/*
SI Metric Prefixes
1000^1 k kilo
1000^2 M mega
1000^3 G giga
1000^4 T tera
1000^5 P peta
1000^6 E exa
1000^7 Z zetta
1000^8 Y yotta
*/

const (
    PrefixKilo  int64 = 1e3
    PrefixMega  int64 = 1e6
    PrefixGiga  int64 = 1e9
    PrefixTera  int64 = 1e12
    PrefixPeta  int64 = 1e15
    PrefixExa   int64 = 1e18
)

/*
IEC Binary prefixes
1024^1 Ki kibi
1024^2 Mi mebi
1024^3 Gi gibi
1024^4 Ti tebi
1024^5 Pi pebi
1024^6 Ei exbi
1024^7 Zi zebi
1024^8 Yi yobi
*/

const (
    PrefixKibi int64 = 1024
    PrefixMebi = PrefixKibi * 1024
    PrefixGibi = PrefixMebi * 1024
    PrefixTebi = PrefixGibi * 1024
    PrefixPebi = PrefixTebi * 1024
    PrefixExbi = PrefixPebi * 1024
)

// Optionally reads a space separating a number from a unit (e.g. "1 KiB") and if found, returns 1 to advance the
// input stream by one rune, or 0 to keep it in the same place.
func acceptUnitSeparator(format *NumberFormat, text string) int {
    // https://physics.nist.gov/cuu/Units/checklist.html
    if (len(text) > 0) && (text[0] == ' ') { return 1; } else { return 0; }
}
