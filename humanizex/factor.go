package humanizex

// Factors describes a way to format a quantity with units.
type Factors struct{
    // Factors is a list of Factor entries in ascending order of size.
    Factors []Factor

    // Components controls how the formatting is broken up -
    //
    // - if zero (default) or 1 (interchangeable), formatting has a single
    // component e.g. "1.5 M".
    //
    // - if 2 or more, formatting is broken up into previous factors e.g.
    // "1 h 50 min" (2 components) or "1 h 50 min 25 s" (3 components)
    Components int
}

// Factor defines one entry in an ordered list of Factors.
type Factor struct{

    // Magnitude defines the absolute size of the factor e.g. 1000000 for the
    // SI unit prefix "M".
    Magnitude float64

    // Label describes the magnitude, usually as a unit prefix (like SI "M")
    // or as a replacement (like "min"), controlled by Mode.
    Unit Unit

    // Mode controls the formatting of this factor
    Mode FactorMode
}

// FactorMode controls the formatting of a factor.
type FactorMode int

const (
    // FactorModeIdentity indicates that the given factor label represents the
    // unit with no changes.
    FactorModeIdentity   = FactorMode(0)

    // FactorModeUnitPrefix indicates the given factor label is a unit prefix
    // e.g. "Ki" is a byte prefix giving "KiB".
    FactorModeUnitPrefix = FactorMode(1)

    // FactorModeReplace indicates the given factor label replaces the current
    // unit e.g. the duration of time 100 s becomes 1 min 40 s, not 1 hs
    // (hectosecond)!
    FactorModeReplace    = FactorMode(2)

    // FactorModeInputCompat indicates that the given factor label should only
    // be considered on input. This may be combined with any other FactorMode
    // by a bitwise OR operation.
    FactorModeInputCompat = FactorMode(4)
)

// bracket returns the index of the first Factor greater or equal to n except
// if n is smaller than the first Factor, returns the first Factor (zero).
// Excludes factors with mode FactorModeInputCompat.
func (f Factors) bracket(n float64) int {
    if len(f.Factors) == 0 {
        panic("operation not defined on an empty list of factors")
    }

    if n < f.Factors[0].Magnitude { return 0 }

    for i, factor := range f.Factors[1:] {
        if factor.Mode & FactorModeInputCompat == FactorModeInputCompat {
            continue // skip
        }

        if n < factor.Magnitude {
            return i
        }
    }

    return len(f.Factors) - 1
}
