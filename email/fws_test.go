package email

import (
    "testing"
)

// TestFWS tests folding whitespace encoding of long headers
func TestFWS(t *testing.T) {
    var tests = [][5]interface{}{
        // --- typical cases, no split ---
        // 012345678901234567890123456789
        // Subject: test subject
        // Subject: test
        //  subject2
        // Subject: test subject
        //  subject test subject
        {"Subject", "test",         "test", 21, nil},
        {"Subject", "test subject", "test subject", 21, nil},
        {"Subject", "test subject", "test subject", 20, nil},
        {"Subject", "test subject2", "test\r\n subject2", 20, nil},
        {"Subject", "test subject subject test subject", "test subject\r\n subject test subject", 20, nil},

        // --- typical cases, one line split on next token ---
        // 012345678901234567890123456789
        // Subject: test subject
        {"Subject", "test subject", "test\r\n subject", 19, nil},

        // --- typical cases, multi line split ---
        // 012345678901234567890123456789
        // Subject: test subject
        //  with a second line
        {"Subject", "test subject with a second line", "test subject\r\n with a second line", 20, nil},

        // 012345678901234567890123456789
        // Subject: test subject
        //  with a second line
        //  and also a third
        //  line
        {"Subject", "test subject with a second line and also a third line", "test subject\r\n with a second line\r\n and also a third\r\n line", 20, nil},

        // --- whitespace tests ---
        {"Subject", "", "", 100, nil},
        {"Subject", "   ", "   ", 100, nil},

        // 012345678901234567890123456789
        // Subject: ****
        // **
        {"Subject", "      ", "    \r\n  ", 12, nil},

        // 012345678901234567890123456789
        // Subject: test   subject
        {"Subject", "test   subject", "test  \r\n subject", 21, nil},

        // 012345678901234567890123456789
        // Subject:    test subject
        {"Subject", "   test subject", "   test\r\n subject", 22, nil},


        // --- error cases ---
        // 012345678901234567890123456789
        // Subject: 0123456789012
        // Subject: test
        //  01234567890123456789
        {"Subject", "01234567890",                      "01234567890", 20, nil},
        {"Subject", "01234567890 foo",                  "01234567890\r\n foo", 20, nil},
        {"Subject", "0123456789012",                    "", 20, fwsWrapErr},
        {"Subject", "0123456789012 foo",                "", 20, fwsWrapErr},
        {"Subject", "test 0123456789012345678",         "test\r\n 0123456789012345678", 20, nil},
        {"Subject", "test 01234567890123456789",        "", 20, fwsWrapErr},
        {"Subject", "01234567890 0123456789012345678",  "01234567890\r\n 0123456789012345678", 20, nil},

    }

    for index, test := range tests {
        key, value, expected, maxLine :=
            test[0].(string), test[1].(string), test[2].(string), test[3].(int)

        var expectedErr error = nil
        if test[4] != nil { expectedErr = test[4].(error) }

        result, err := fwsWrap(value, len(key), maxLine)
        if err != expectedErr {
            t.Errorf("Test %d failed: fwsWrap(%q, %q, %d), got result %q and error %v but expected error %v\n",
                index, key, value, maxLine, result, err, expectedErr)
        } else if result != expected {
            t.Errorf("Test %d failed: fwsWrap(%q, %q, %d), got %q but wanted %q\n",
                index, key, value, maxLine, result, expected)
        }
    }
}
