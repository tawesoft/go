package operator

import (
    "testing"
)

const T = true
const F = false

func TestBooleanUnary(t *testing.T) {
    type test struct {
        p bool
        f func(bool) bool
        expected bool
    }
    
    var tests = []test {
        {T, Bool.Unary.True,     T}, // test 0
        {F, Bool.Unary.True,     T}, // test 1
        {T, Bool.Unary.False,    F}, // test 2
        {F, Bool.Unary.False,    F}, // test 3
        {T, Bool.Unary.Not,      F}, // test 4
        {F, Bool.Unary.Not,      T}, // test 5
        {T, Bool.Unary.Identity, T}, // test 6
        {F, Bool.Unary.Identity, F}, // test 7
    }
    
    for idx, i := range tests {
        var result = i.f(i.p)
        if result != i.expected {
            t.Errorf("test %d: got %t, but expected %t", idx, result, i.expected)
        }
    }
}

func TestBooleanBinaru(t *testing.T) {
    var f = Bool.Binary
    
    var testfns = []func(bool, bool)bool{
        f.True,          // test  0
        f.False,         // test  1
        f.P,             // test  2
        f.Q,             // test  3
        f.NotP,          // test  4
        f.NotQ,          // test  5
        f.Eq,            // test  6
        f.Neq,           // test  7
        f.And,           // test  8
        f.Nand,          // test  9
        f.Or,            // test 10
        f.Nor,           // test 11
        f.Xor,           // test 12
        f.Xnor,          // test 13
        f.Implies,       // test 14
        f.NonImplies,    // test 15
        f.ConImplies,    // test 16
        f.ConNonImplies, // test 17
    }
    
    var truthTable = []bool{
        // 0     1      2  3  4     5     6   7    8    9     10  11   12   13    14  15    16  17
        // True, False, P, Q, NotP, NotQ, Eq, Neq, And, Nand, Or, Nor, Xor, Xnor, =>, =/=>, <=, <=/=
        
        // p: T, q: T (case 0)
           T,    F,     T, T, F,   F,     T,  F,   T,   F,    T,  F,   F,   T,    T,  F,    T,  F,
           
        // p: T, q: F (case 1)
           T,    F,     T, F, F,   T,     F,  T,   F,   T,    T,  F,   T,   F,    F,  T,    T,  F,
           
        // p: F, q: T (case 2)
           T,    F,     F, T, T,   F,     F,  T,   F,   T,    T,  F,   T,   F,    T,  F,    F,  T,
           
        // p: F, q: F (case 3)
           T,    F,     F, F, T,   T,     T,  F,   F,   T,    F,  T,   F,   T,    T,  F,    T,  F,
    }

    for idx, i := range testfns {
        var results = []bool{i(T, T), i(T, F), i(F, T), i(F, F)}
        var expected = []bool{
            truthTable[0 * len(testfns) + idx],
            truthTable[1 * len(testfns) + idx],
            truthTable[2 * len(testfns) + idx],
            truthTable[3 * len(testfns) + idx],
        }
        
        for n := 0; n < 4; n++ {
            if results[n] == expected[n] { continue }
            t.Errorf("test %d case %d: got %t but expected %t", idx, n, results[n], expected[n])
        }
    }
}

func TestBooleanNary(t *testing.T) {
    type test struct {
        f func(...bool) bool
        p []bool
        expected bool
    }
    
    var tests = []test {
        {Bool.Nary.True, []bool{},  T},             // test  0
        {Bool.Nary.True, []bool{T}, T},             // test  1
        {Bool.Nary.True, []bool{F}, T},             // test  2
        {Bool.Nary.True, []bool{T, F, T, F},  T},   // test  3
        
        {Bool.Nary.False, []bool{},  F},            // test  4
        {Bool.Nary.False, []bool{T}, F},            // test  5
        {Bool.Nary.False, []bool{F}, F},            // test  6
        {Bool.Nary.False, []bool{T, F, T, F}, F},   // test  7
        
        {Bool.Nary.All,   []bool{},  T},            // test  8
        {Bool.Nary.All,   []bool{T}, T},            // test  9
        {Bool.Nary.All,   []bool{F}, F},            // test 10
        {Bool.Nary.All,   []bool{T, T, T, T}, T},   // test 11
        {Bool.Nary.All,   []bool{T, T, F, T}, F},   // test 12
        
        {Bool.Nary.Any,   []bool{},  F},            // test 13
        {Bool.Nary.Any,   []bool{T}, T},            // test 14
        {Bool.Nary.Any,   []bool{F}, F},            // test 15
        {Bool.Nary.Any,   []bool{T, T, T, T}, T},   // test 16
        {Bool.Nary.Any,   []bool{T, T, F, T}, T},   // test 17
        {Bool.Nary.Any,   []bool{F, F, F, F}, F},   // test 18
        
        {Bool.Nary.None,  []bool{},  T},            // test 19
        {Bool.Nary.None,  []bool{T}, F},            // test 20
        {Bool.Nary.None,  []bool{F}, T},            // test 21
        {Bool.Nary.None,  []bool{T, T, T, T}, F},   // test 22
        {Bool.Nary.None,  []bool{T, T, F, T}, F},   // test 23
        {Bool.Nary.None,  []bool{F, F, F, F}, T},   // test 24
        
        {boolNaryAll1,    []bool{},  T},            // test 25
        {boolNaryAll1,    []bool{T}, T},            // test 26
        {boolNaryAll1,    []bool{F}, F},            // test 27
        {boolNaryAll1,    []bool{T, T, T, T}, T},   // test 28
        {boolNaryAll1,    []bool{T, T, F, T}, F},   // test 29
        
        {boolNaryAll2,    []bool{},  T},            // test 30
        {boolNaryAll2,    []bool{T}, T},            // test 31
        {boolNaryAll2,    []bool{F}, F},            // test 32
        {boolNaryAll2,    []bool{T, T, T, T}, T},   // test 33
        {boolNaryAll2,    []bool{T, T, F, T}, F},   // test 34
        
        {boolNaryAny1,    []bool{},  F},            // test 35
        {boolNaryAny1,    []bool{T}, T},            // test 36
        {boolNaryAny1,    []bool{F}, F},            // test 37
        {boolNaryAny1,    []bool{T, T, T, T}, T},   // test 38
        {boolNaryAny1,    []bool{T, T, F, T}, T},   // test 39
        {boolNaryAny1,    []bool{F, F, F, F}, F},   // test 40
        
        {boolNaryAny2,    []bool{},  F},            // test 41
        {boolNaryAny2,    []bool{T}, T},            // test 42
        {boolNaryAny2,    []bool{F}, F},            // test 43
        {boolNaryAny2,    []bool{T, T, T, T}, T},   // test 44
        {boolNaryAny2,    []bool{T, T, F, T}, T},   // test 45
        {boolNaryAny2,    []bool{F, F, F, F}, F},   // test 46
        
        {boolNaryNone1,   []bool{},  T},            // test 47
        {boolNaryNone1,   []bool{T}, F},            // test 48
        {boolNaryNone1,   []bool{F}, T},            // test 49
        {boolNaryNone1,   []bool{T, T, T, T}, F},   // test 50
        {boolNaryNone1,   []bool{T, T, F, T}, F},   // test 51
        {boolNaryNone1,   []bool{F, F, F, F}, T},   // test 52
        
        {boolNaryNone2,   []bool{},  T},            // test 53
        {boolNaryNone2,   []bool{T}, F},            // test 54
        {boolNaryNone2,   []bool{F}, T},            // test 55
        {boolNaryNone2,   []bool{T, T, T, T}, F},   // test 56
        {boolNaryNone2,   []bool{T, T, F, T}, F},   // test 57
        {boolNaryNone2,   []bool{F, F, F, F}, T},   // test 58
    }
    
    for idx, i := range tests {
        var result = i.f(i.p...)
        if result != i.expected {
            t.Errorf("test %d: got %t, but expected %t", idx, result, i.expected)
        }
    }
}
