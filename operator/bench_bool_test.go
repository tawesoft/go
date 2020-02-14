package operator

// Conclusion: the v1 version of each function wins in every case except BenchmarkBoolNaryAll2Short

import (
    "testing"
)

const benchloops = 1000000

// ===[ BoolNaryAll1 ]==============================================================================[ BoolNaryAll1 ]===

func BenchmarkBoolNaryAll1Short1(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAll1(T, T, T, T) != T { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAll1Short2(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAll1(F, T, T, T) != F { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAll1Short3(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAll1(F, T, T, T) != F { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAll1Long1(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAll1(T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T) != T { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAll1Long2(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAll1(T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, F) != F { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAll1Long3(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAll1(F, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T) != F { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAll1Mixed(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAll1(T, T, T, T) != T { panic("assertion failed") }
        if boolNaryAll1(F, T, T, T) != F { panic("assertion failed") }
        if boolNaryAll1(F, T, T, T) != F { panic("assertion failed") }
        if boolNaryAll1(T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T) != T { panic("assertion failed") }
        if boolNaryAll1(T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, F) != F { panic("assertion failed") }
        if boolNaryAll1(F, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T) != F { panic("assertion failed") }
    }
}

// ===[ BoolNaryAll2 ]==============================================================================[ BoolNaryAll2 ]===

func BenchmarkBoolNaryAll2Short1(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAll2(T, T, T, T) != T { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAll2Short2(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAll2(F, T, T, T) != F { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAll2Short3(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAll2(F, T, T, T) != F { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAll2Long1(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAll2(T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T) != T { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAll2Long2(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAll2(T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, F) != F { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAll2Long3(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAll2(F, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T) != F { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAll2Mixed(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAll2(T, T, T, T) != T { panic("assertion failed") }
        if boolNaryAll2(F, T, T, T) != F { panic("assertion failed") }
        if boolNaryAll2(F, T, T, T) != F { panic("assertion failed") }
        if boolNaryAll2(T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T) != T { panic("assertion failed") }
        if boolNaryAll2(T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, F) != F { panic("assertion failed") }
        if boolNaryAll2(F, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T) != F { panic("assertion failed") }
    }
}

// ===[ BoolNaryAny1 ]==============================================================================[ BoolNaryAny1 ]===

func BenchmarkBoolNaryAny1Short1(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAny1(F, F, F, F) != F { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAny1Short2(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAny1(T, F, F, F) != T { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAny1Short3(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAny1(F, F, F, T) != T { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAny1Long1(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAny1(F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F) != F { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAny1Long2(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAny1(T, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F) != T { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAny1Long3(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAny1(F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, T) != T { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAny1Mixed(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAny1(F, F, F, F) != F { panic("assertion failed") }
        if boolNaryAny1(T, F, F, F) != T { panic("assertion failed") }
        if boolNaryAny1(F, F, F, T) != T { panic("assertion failed") }
        if boolNaryAny1(F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F) != F { panic("assertion failed") }
        if boolNaryAny1(T, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F) != T { panic("assertion failed") }
        if boolNaryAny1(F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, T) != T { panic("assertion failed") }
    }
}

// ===[ BoolNaryAny2 ]==============================================================================[ BoolNaryAny2 ]===

func BenchmarkBoolNaryAny2Short1(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAny2(F, F, F, F) != F { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAny2Short2(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAny2(T, F, F, F) != T { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAny2Short3(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAny2(F, F, F, T) != T { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAny2Long1(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAny2(F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F) != F { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAny2Long2(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAny2(T, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F) != T { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAny2Long3(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAny2(F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, T) != T { panic("assertion failed") }
    }
}

func BenchmarkBoolNaryAny2Mixed(b *testing.B) {
    for i := 0; i < benchloops; i++ {
        if boolNaryAny2(F, F, F, F) != F { panic("assertion failed") }
        if boolNaryAny2(T, F, F, F) != T { panic("assertion failed") }
        if boolNaryAny2(F, F, F, T) != T { panic("assertion failed") }
        if boolNaryAny2(F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F) != F { panic("assertion failed") }
        if boolNaryAny2(T, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F) != T { panic("assertion failed") }
        if boolNaryAny2(F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, T) != T { panic("assertion failed") }
    }
}

