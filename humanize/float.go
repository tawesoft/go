package humanize

import (
    "fmt"
    "math"
)

func formatPureFloat(formatter *Formatter, suffix string, value float64) string {
    var i, f = math.Modf(value); f = math.Abs(f)
    var rounder = math.Copysign(0.5, value)
    
    if value == 0.0 {
        return fmt.Sprintf("0%s%s", formatter.UnitSeparator, suffix)
    }
    
    if (formatter.DecimalPlaces < 1) || f == 0.0 {
        // || (math.Abs(f) <= math.Pow(0.1, float64(formatter.SignificantFigures - 1))) {
        return fmt.Sprintf("%d%s%s", int64(i+rounder), formatter.UnitSeparator, suffix)
    }
    
    return fmt.Sprintf("%d%s%0*d%s%s",
        int64(i+rounder),
        formatter.DecimalSeparator,
        formatter.DecimalPlaces,
        int64((f*math.Pow10(formatter.DecimalPlaces))+0.5),
        formatter.UnitSeparator,
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

func formatFloat(formatter *Formatter, mappings []mappingSI, value float64) string {
    if formatter == nil { formatter = &DefaultFormatter }
    
    if math.IsInf(value, 1) {
        return "Infinity"
    } else if math.IsInf(value, -1) {
        return "-Infinity"
    } else if math.IsNaN(value) {
        return "NaN"
    } else if math.Abs(value) < 1.0 {
        return formatPureFloat(formatter, "", value) // TODO small floats
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
        
        return formatPureFloat(formatter, unit, value / divisor)
    }
}

func FormatFloat(formatter *Formatter, value float64) string {
    return formatFloat(formatter, mappingsSI, value)
}

func FormatFloatIEC(formatter *Formatter, value float64) string {
    return formatFloat(formatter, mappingsIEC, value)
}
