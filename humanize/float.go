package humanize

import (
    "fmt"
    "math"
)

func formatPureFloat(format *Format, sigfigs int, suffix string, value float64) string {
    var i, f = math.Modf(value); f = math.Abs(f)
    var rounder = math.Copysign(0.5, value)
    
    if value == 0.0 { return fmt.Sprintf("0 %s", suffix) }
    
    if (sigfigs < 1) || f == 0.0 {
        // || (math.Abs(f) <= math.Pow(0.1, float64(formatter.SignificantFigures - 1))) {
        return fmt.Sprintf("%d %s", int64(i+rounder), suffix)
    }
    
    return fmt.Sprintf("%d%c%0*d %s",
        int64(i+rounder),
        format.DecimalSeparator,
        sigfigs,
        int64((f*math.Pow10(sigfigs))+0.5),
        suffix)
}

type mappingSI struct {
    places int
    unit string
    divisor int64
}

var mappingsSI = []mappingSI{
    {19, "E", PrefixExa},
    {16, "P", PrefixPeta},
    {13, "T", PrefixTera},
    {10, "G", PrefixGiga},
    { 7, "M", PrefixMega},
    { 4, "k", PrefixKilo},
}

var mappingsIEC = []mappingSI{
    {19, "Ei", PrefixExbi},
    {16, "Pi", PrefixPebi},
    {13, "Ti", PrefixTebi},
    {10, "Gi", PrefixGibi},
    { 7, "Mi", PrefixMebi},
    { 4, "Ki", PrefixKibi},
}

func formatFloat(format *Format, mappings []mappingSI, sigfigs int, value float64) string {
    if format == nil { format = &SimpleFormat }
    
    if math.IsInf(value, 1) {
        return "Infinity"
    } else if math.IsInf(value, -1) {
        return "-Infinity"
    } else if math.IsNaN(value) {
        return "NaN"
    } else if math.Abs(value) < 1.0 {
        return formatPureFloat(format, sigfigs, "", value)
    } else {
        var places = int(math.Log10(math.Abs(value))) + 1 // e.g. for 1500, 4
        
        var unit string
        var divisor float64 = 1.0
        
        for _, i := range mappings {
            if places >= i.places {
                unit = i.unit
                divisor = float64(i.divisor)
                break
            }
        }
        
        return formatPureFloat(format, sigfigs, unit, value / divisor)
    }
}

func FormatFloatSI(formatter *Format, sigfigs int, value float64) string {
    return formatFloat(formatter, mappingsSI, sigfigs, value)
}

func FormatFloatIEC(formatter *Format, sigfigs int, value float64) string {
    return formatFloat(formatter, mappingsIEC, sigfigs, value)
}
