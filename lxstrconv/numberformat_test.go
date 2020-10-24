package lxstrconv

import (
    "reflect"
    "testing"
    
    "golang.org/x/text/language"
)

// TestDecimalNumberFormatLocale tests that the decimalFormat correctly detects the
// correct group and decimal separators and digits for a given locale.
//
// It is possible that these tests may become stale and fail with newer
// versions of Unicode or newer versions of golang.org/x/text. In that case,
// the test is incorrect, not the code.
//
// It's also possible that these tests are incorrect due to my own
// misunderstanding of certain locales. Corrections are welcome!
func TestDecimalNumberFormatLocale(t *testing.T) {
    type test struct {
        lang language.Tag
        expectedGroupSeparator rune
        expectedPoint rune
        expectedDigits []rune
    }
    
    digitsLatin   := []rune("0123456789")
    digitsArabic  := []rune("٠١٢٣٤٥٦٧٨٩")
    digitsBengali := []rune("০১২৩৪৫৬৭৮৯")
    
    var tests = []test{
        {language.BritishEnglish,   ',',        '.',      digitsLatin},
        {language.French,           '\u00a0',   ',',      digitsLatin},
        {language.Dutch,            '.',        ',',      digitsLatin},
        {language.Arabic,           '\u066c',   '\u066b', digitsArabic},
        {language.Malayalam,        ',',        '.',      digitsLatin},
        {language.Bengali,          ',',        '.',      digitsBengali},
    }
    
    for _, test := range tests {
        f := NewDecimalFormat(test.lang).(decimalFormat)
        
        if f.GroupSeparator != test.expectedGroupSeparator {
            t.Errorf("expected group separator %c (0x%x) for %s, but got %c (0x%x)",
                test.expectedGroupSeparator, test.expectedGroupSeparator,
                test.lang,
                f.GroupSeparator, f.GroupSeparator)
        }
        
        if f.Point != test.expectedPoint {
            t.Errorf("expected point %c (0x%x) for %s, but got %c (0x%x)",
                test.expectedPoint, test.expectedPoint,
                test.lang,
                f.Point, f.Point)
        }
        
        if !reflect.DeepEqual(f.Digits[:], test.expectedDigits) {
            t.Errorf("expected digits %s (%+v) for %s, but got %s (%+v)",
                string(test.expectedDigits), test.expectedDigits,
                test.lang,
                string(f.Digits[:]), f.Digits)
        }
    }
}


// TestDecimalNumberFormatParse tests parsing of decimal numbers in different
// locales.
func TestDecimalNumberFormatParse(t *testing.T) {
    type test struct {
        lang language.Tag
        in string
        
        expectedIntValue int64
        expectedIntLen   int
        
        expectedFloatValue float64
        expectedFloatLen int
    }
    
    var tests = []test{
        {language.BritishEnglish, "1,234.56", 1234, 5, 1234.56, 8},
        {language.French,         "1 234,56", 1234, 5, 1234.56, 8},
        {language.Arabic,         "١\u066c٢٣٤\u066b٥٦", 1234, 10, 1234.56, 16},
    }
    
    for _, test := range tests {
        f := NewDecimalFormat(test.lang)
        
        {
            value, length, err := f.AcceptInt(test.in)
            if err != nil {
                t.Errorf("unexpected error for %s %v: %v", test.lang, test.in, err)
                continue
            }
            
            if length != test.expectedIntLen {
                t.Errorf("expected int len %d for %s %v but got %d",
                test.expectedIntLen, test.lang, test.in, length)
            }
            
            if value != test.expectedIntValue {
                t.Errorf("expected int value %d for %s %v but got %d",
                test.expectedIntValue, test.lang, test.in, value)
            }
        }
        
        {
            value, length, err := f.AcceptFloat(test.in)
            if err != nil {
                t.Errorf("unexpected error for %s %v: %v", test.lang, test.in, err)
                continue
            }
            
            if length != test.expectedFloatLen {
                t.Errorf("expected float len %d for %s %v but got %d",
                test.expectedFloatLen, test.lang, test.in, length)
            }
            
            if value != test.expectedFloatValue {
                t.Errorf("expected float value %f for %s %v but got %f",
                test.expectedFloatValue, test.lang, test.in, value)
            }
        }
    }
}
