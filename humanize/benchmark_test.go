package humanize

import (
    "testing"
    dhumanize "github.com/dustin/go-humanize"
)

func BenchmarkTawesoftFormatBytes(b *testing.B) {
    for i := 0; i < b.N; i++ {
        value := FormatBytesIEC(nil, 42*1024*1024)
        if value != "42 MiB" { panic("assertion failed") }
    }
}

func BenchmarkDustinFormatBytes(b *testing.B) {
    for i := 0; i < b.N; i++ {
        value := dhumanize.Bytes(82*1e6)
        if value != "82 MB" { panic("assertion failed") }
    }
}

func BenchmarkTawesoftParseBytes(b *testing.B) {
    for i := 0; i < b.N; i++ {
        value, err := ParseBytes("42 MiB")
        if value != 44040192 && err != nil { panic("assertion failed") }
    }
}

func BenchmarkDustinParseBytes(b *testing.B) {
    for i := 0; i < b.N; i++ {
        value, err := dhumanize.ParseBytes("42 MiB")
        if value != 44040192 && err != nil { panic("assertion failed") }
    }
}
