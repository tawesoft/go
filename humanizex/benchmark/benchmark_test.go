package main

import (
    "testing"
    "golang.org/x/text/language"
    thumanize "tawesoft.co.uk/go/humanizex"
    dhumanize "github.com/dustin/go-humanize"
)

func BenchmarkTawesoftFormatBytes(b *testing.B) {
    h := thumanize.NewHumanizer(language.English)

    for i := 0; i < b.N; i++ {
        value := h.FormatBytesIEC(8.5*1024*1024*1024)
        if value != "8.5 GiB" { panic("assertion failed (got "+value+")") }
    }
}

func BenchmarkDustinFormatBytes(b *testing.B) {
    for i := 0; i < b.N; i++ {
        value := dhumanize.IBytes(8.5*1024*1024*1024)
        if value != "8.5 GiB" { panic("assertion failed (got "+value+")") }
    }
}

func BenchmarkTawesoftFormatFloatSI(b *testing.B) {
    h := thumanize.NewHumanizer(language.English)

    for i := 0; i < b.N; i++ {
        value := h.Format(8.5*1000*1000*1000,
            thumanize.CommonUnits.None, thumanize.CommonFactors.SI).Utf8
        if value != "8.5 G" { panic("assertion failed (got "+value+")") }
    }
}

func BenchmarkDustinFormatFloatSI(b *testing.B) {
    for i := 0; i < b.N; i++ {
        value := dhumanize.SI(8.5*1000*1000*1000, "")
        if value != "8.5 G" { panic("assertion failed (got "+value+")") }
    }
}

func BenchmarkTawesoftParseBytes(b *testing.B) {
    h := thumanize.NewHumanizer(language.English)

    for i := 0; i < b.N; i++ {
        value, err := h.ParseBytesIEC("8.5 GiB")
        if err != nil { b.Errorf("unexpected error %v", err) }
        if value != 8.5 * 1024 * 1024 * 1024 { b.Errorf("assertion failed (got %d)", value) }
    }
}


func BenchmarkDustinParseBytes(b *testing.B) {
    for i := 0; i < b.N; i++ {
        value, err := dhumanize.ParseBytes("8.5 GiB")
        if err != nil { b.Errorf("unexpected error %v", err) }
        if value != 8.5 * 1024 * 1024 * 1024 { b.Errorf("assertion failed (got %d)", value) }
    }
}
