package humanizex

var CommonUnits = struct{
    None           Unit
    Second         Unit
    Meter          Unit
    Byte           Unit
    Bit            Unit
    BitsPerSecond  Unit
}{
    None:          Unit{"", ""},
    Second:        Unit{"s", "s"},
    Meter:         Unit{"m", "m"},
    Byte:          Unit{"B", "B"},
    Bit:           Unit{"b", "b"},
    BitsPerSecond: Unit{"bps", "bps"},
}

var CommonFactors = struct{
    // Time is time units in seconds, minutes, hours, days and years as min, h,
    // d, and y. These are non-SI units but generally accepted in context.
    // For times smaller than a second (e.g. nanoseconds), use SI instead.
    // The expected unit is a second (Unit{"s", "s"} or CommonUnits.Second)
    Time Factors

    // Distance are SI units that stop at kilo (because nobody uses
    // megametres or gigametres!) but includes centi. The expected unit is the
    // SI unit for distance, the metre (Unit{"m", "m"} or CommonUnits.Meter)
    Distance Factors

    // IEC are the "ibi" unit prefixes for bytes e.g. Ki, Mi, Gi with a
    // factor of 1024.
    IEC Factors

    // JEDEC are the old unit prefixes for bytes: K, M, G (only) with a factor
    // of 1024.
    JEDEC Factors

    // SIBytes are the SI unit prefixes for bytes e.g. k, M, G with a
    // factor of 1000. Unlike the normal SI Factors, it is assumed based on
    // context that when a "K" is input this is intended to mean the "k" SI
    // unit prefix instead of Kelvin - I've never heard of a Kelvin-Byte!
    SIBytes Factors

    // SIUncommon are the SI unit prefixes including deci, deca, and hecto
    SIUncommon Factors

    // SI are the SI unit prefixes except centi, deci, deca, and hecto
    SI Factors
}{
    Time: Factors{
        Factors: []Factor{
            {1,                                 Unit{"s", "s"},     FactorModeReplace},
            {60,                                Unit{"min", "min"}, FactorModeReplace},
            {60 * 60,                           Unit{"h", "h"},     FactorModeReplace},
            {24 * 60 * 60,                      Unit{"d", "d"},     FactorModeReplace},
            {365.2422 * 24 * 60 * 60,           Unit{"y", "y"},     FactorModeReplace},
        },
        Components: 2,
    },
    Distance: Factors{
        Factors: []Factor{
            {1E-9,                              Unit{"n", "n"},     FactorModeUnitPrefix}, // nano
            {1E-6,                              Unit{"μ", "u"},     FactorModeUnitPrefix}, // micro
            {1E-3,                              Unit{"m", "m"},     FactorModeUnitPrefix}, // milli
            {1E-2,                              Unit{"c", "c"},     FactorModeUnitPrefix}, // centi
            {1,                                 Unit{ "",  ""},     FactorModeIdentity},
            {1000,                              Unit{"k", "k"},     FactorModeUnitPrefix}, // kilo
        },
    },
    IEC: Factors{
        Factors: []Factor{
            {1,                                 Unit{ "",  ""},     FactorModeUnitPrefix},
            {1024,                              Unit{"Ki", "Ki"},   FactorModeUnitPrefix},
            {1024 * 1024,                       Unit{"Mi", "Mi"},   FactorModeUnitPrefix},
            {1024 * 1024 * 1024,                Unit{"Gi", "Gi"},   FactorModeUnitPrefix},
            {1024 * 1024 * 1024 * 1024,         Unit{"Ti", "Ti"},   FactorModeUnitPrefix},
        },
    },
    JEDEC: Factors{
        Factors: []Factor{
            {1,                                 Unit{ "",  ""},     FactorModeIdentity},
            {1024,                              Unit{"K", "K"},     FactorModeUnitPrefix},
            {1024 * 1024,                       Unit{"M", "M"},     FactorModeUnitPrefix},
            {1024 * 1024 * 1024,                Unit{"G", "G"},     FactorModeUnitPrefix},
        },
    },
    SIBytes: Factors{
        Factors: []Factor{
            {1,                                 Unit{ "",  ""},     FactorModeIdentity},
            { 1E3,                              Unit{"k", "k"},     FactorModeUnitPrefix},
            { 1E3,                              Unit{"K", "K"},     FactorModeUnitPrefix | FactorModeInputCompat}, // Kelvin-Bytes(!)
            { 1E6,                              Unit{"M", "M"},     FactorModeUnitPrefix},
            { 1E9,                              Unit{"G", "G"},     FactorModeUnitPrefix},
            {1E12,                              Unit{"T", "T"},     FactorModeUnitPrefix},
        },
    },
    SIUncommon: Factors{
        Factors: []Factor{
            {1E-9,                              Unit{"n", "n"},     FactorModeUnitPrefix}, // nano
            {1E-6,                              Unit{"μ", "u"},     FactorModeUnitPrefix}, // micro
            {1E-3,                              Unit{"m", "m"},     FactorModeUnitPrefix}, // milli
            {1E-2,                              Unit{"c", "c"},     FactorModeUnitPrefix}, // centi
            {1E-1,                              Unit{"d", "d"},     FactorModeUnitPrefix}, // deci
            {1,                                 Unit{ "",  ""},     FactorModeIdentity},
            { 1E1,                              Unit{"da", "da"},   FactorModeUnitPrefix}, // deca
            { 1E2,                              Unit{"h", "h"},     FactorModeUnitPrefix}, // hecto
            { 1E3,                              Unit{"k", "k"},     FactorModeUnitPrefix}, // kilo
            { 1E6,                              Unit{"M", "M"},     FactorModeUnitPrefix},
            { 1E9,                              Unit{"G", "G"},     FactorModeUnitPrefix},
            {1E12,                              Unit{"T", "T"},     FactorModeUnitPrefix},
        },
    },
    SI: Factors{
        Factors: []Factor{
            {1E-9,                              Unit{"n", "n"},     FactorModeUnitPrefix}, // nano
            {1E-6,                              Unit{"μ", "u"},     FactorModeUnitPrefix}, // micro
            {1E-3,                              Unit{"m", "m"},     FactorModeUnitPrefix}, // milli
            {1,                                 Unit{ "",  ""},     FactorModeIdentity},
            { 1E3,                              Unit{"k", "k"},     FactorModeUnitPrefix}, // kilo
            { 1E6,                              Unit{"M", "M"},     FactorModeUnitPrefix},
            { 1E9,                              Unit{"G", "G"},     FactorModeUnitPrefix},
            {1E12,                              Unit{"T", "T"},     FactorModeUnitPrefix},
        },
    },
}
