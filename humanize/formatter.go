package humanize

// https://physics.nist.gov/cuu/Units/checklist.html

type Format struct {
    // GroupSeparator separates groups of digits e.g. thousands
    GroupSeparator     rune // e.g. `1000000 => 1<-GroupSeparator->000<-GroupSeparator->000 => "1,000,000"`
    
    // DecimalSeparator separates integer and fractional parts of a decimal number
    DecimalSeparator   rune // e.g. `1.23 => 1<--DecimalSeparator-->23 => "1.23"`
    
    // Group digits is how many digits to group by e.g. 3 for thousands
    GroupDigits int // e.g. for 3, `1000000 => "1,000,000"`
    
    // GroupMinDigits is how many digits are required to start grouping
    GroupMinDigits int // e.g. for 5, `1234 => "1234"` but `12345 => "12,345"`
    
    // Grouper is an optional function for more complicated numbering systems that returns true to insert a group
    // separator at a given index for a string of a certain length. See the example for the IndianFormatter. Length is
    // the total number of digits on that side of the decimal point. Index tends left towards positive infinity and
    // tends right towards negative infinity around the decimal point.
    //
    // If Grouper is specified, GroupDigits and GroupMinDigits are ignored
    Grouper func(index int, len int) bool
}

var SimpleFormat = Format{
    GroupSeparator:   ',',
    DecimalSeparator: '.',
    GroupDigits:      3,
    GroupMinDigits:   5,
}

var UnicodeFormat = Format{
    GroupSeparator:   '\u2009', // thin Space
    DecimalSeparator: '.',
    GroupDigits:      3,
    GroupMinDigits:   5,
}

var IndianFormat = Format{
    GroupSeparator:   ',',
    DecimalSeparator: '.',
    Grouper: func(index int, len int) bool {
        // len unused in this implementation
        
        // e.g.  1,000,00,00,000,00,00,000.123
        // index 1 10F ED CB A98 76 54 321|123
        //                        positive|negative
        // len   100000000000000000 = 18
        // len   123                =  3
        
        // AFAIK(?) symmetrical around decimal point
        if index < 0 { index = 0 }
        
        // insert a separator at index 3, 5, 7, 10, 12, 14, 18, 20, 22, 25 ...
        //                              +2 +2 +3  +2  +2  +3  +2  +2  +3 ...
        // mod 7                       3, 5, 0, 3,  5,  0,   3, 5,  0,
        var remainder = index % 7
        
        return (remainder == 0) || (remainder == 3) || (remainder == 5)
    },
}
