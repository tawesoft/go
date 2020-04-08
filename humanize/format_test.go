package humanize

import (
    "reflect"
    "runtime"
    "testing"
)

func GetFunctionName(i interface{}) string {
    return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func TestFormatBytes(t *testing.T) {
    
    type test struct {
        value int64
        expected string
        fn func(format *Format, sigfigs int, value int64) string
    }
    
    var tests = []test{
        {    0,                "0 B", FormatBytesSI},
        { 1000,               "1 kB", FormatBytesSI},
        {-1000,              "-1 kB", FormatBytesSI},
        {-1001,           "-1.00 kB", FormatBytesSI},
        {-1050,           "-1.05 kB", FormatBytesSI},
        { 1024,              "1 KiB", FormatBytesIEC},
        {-1024,             "-1 KiB", FormatBytesIEC},
        {-1025,          "-1.00 KiB", FormatBytesIEC},
        {-1076,          "-1.05 KiB", FormatBytesIEC},
        { 1024*1024,         "1 MiB", FormatBytesIEC},
        {-1024*1024,        "-1 MiB", FormatBytesIEC},
        { 1e6,                "1 MB", FormatBytesSI},
        {-1e6,               "-1 MB", FormatBytesSI},
        { 1536,           "1.50 KiB", FormatBytesIEC},
        { 1500,            "1.50 kB", FormatBytesSI},
        { 15000,             "15 kB", FormatBytesSI},
        { 15*1024,           "15 KiB", FormatBytesIEC},
        { 15500,          "15.50 kB", FormatBytesSI},
        { 15*1024 + 512, "15.50 KiB", FormatBytesIEC},
    }
    
    for _, i := range tests {
        var result = i.fn(nil, 2, i.value)
        if result != i.expected {
            t.Errorf("FormatBytes(%d, %s): got %s but expected %s", i.value, GetFunctionName(i.fn), result, i.expected)
        }
    }
}
