package operator

import (
    "math"
    "testing"
)

func TestIntMSB(t *testing.T) {
    type test struct {
        n int
        msb int
    }
    
    var tests = []test {
        //         8421
        { 0, 0}, // 1100
        { 1, 1}, // 1100
        { 8, 4}, // 1000
        {12, 4}, // 1100
        {int(maxInt32), 31},
    }
    
    for idx, i := range tests {
        var result = intMostSignificantBit(i.n)
        if result != i.msb {
            t.Errorf("test %d: got %d, but expected %d", idx, result, i.msb)
        }
    }
}

func TestIntBinary(t *testing.T) {
    type test struct {
        a int
        b int
        f func(int, int) int
        expected int
    }
    
    var tests = []test {
        {5, 7, Int.Binary.Add, 12},
        {7, 5, Int.Binary.Sub,  2},
        {5, 7, Int.Binary.Sub, -2},
        {5, 7, Int.Binary.Mul, 35},
        {6, 2, Int.Binary.Div,  3},
        {1+2+4, 2+8, Int.Binary.And,     2},
        {1+2+4, 2+8, Int.Binary.Or,     15},
        {1+2+4, 2+8, Int.Binary.Xor,    13},
        {1+2+4, 2+8, Int.Binary.AndNot,  5}, // 0111 &^ 1010 = 0111 & 0101 = 0101 = 5
    }

    type testShift struct {
        a int
        b uint
        f func(int, uint) int
        expected int
    }
    
    var testsShift = []testShift {
        {3, 2, Int.Binary.Shl, 12}, // 0011 << 2 = 1100 = 12
        {12, 2, Int.Binary.Shr, 3}, // 1100 >> 2 = 0011 = 3
    }
    
    for idx, i := range tests {
        var result = i.f(i.a, i.b)
        if result != i.expected {
            t.Errorf("test %d: got %d, but expected %d", idx, result, i.expected)
        }
    }
    
    for idx, i := range testsShift {
        var result = i.f(i.a, i.b)
        if result != i.expected {
            t.Errorf("shift test %d: got %d, but expected %d", idx, result, i.expected)
        }
    }
}


func TestIntCheckedShl(t *testing.T) {
    type test struct {
        a int32
        b uint
        expectedValue int32
        expectedError error
    }
    
    var tests = []test {
        {1, 1, 2, nil},
        {3, 1, 6, nil},
        {-1, 1, 0, ErrorUndefined},
        {maxInt32, 1, 0, ErrorOverflow},
        {int32(math.Exp2(31) - 1), 1, 0, ErrorOverflow},
        {int32(math.Exp2(30)), 2, 0, ErrorOverflow},
        {int32(math.Exp2(29)), 1, int32(math.Exp(30)), nil},
        {int32(math.Exp2(30)),  1, 0, ErrorOverflow},
    }
    
    for idx, i := range tests {
        var result, err = int32BinaryCheckedShl(i.a, i.b)
        if err == i.expectedError {
            // pass
        } else if err != nil {
            t.Errorf("test %d: unexpected error %v (expected %v)", idx, err, i.expectedError)
        } else if result != i.expectedValue {
            t.Errorf("test %d: got %d, but expected %d", idx, result, i.expectedValue)
        }
    }
}

func TestUintCheckedShl(t *testing.T) {
    type test struct {
        a uint32
        b uint
        expectedValue uint32
        expectedError error
    }
    
    var tests = []test {
        {1, 1, 2, nil},
        {3, 1, 6, nil},
        {maxInt32, 1, 0, ErrorOverflow},
        {uint32(math.Exp2(31)), 1, 0, ErrorOverflow},
        {uint32(math.Exp2(32) - 1), 1, 0, ErrorOverflow},
        {uint32(math.Exp2(30)), 1, uint32(math.Exp(31)), nil},
        {uint32(math.Exp2(31)), 1, uint32(math.Exp(32)), nil},
        {uint32(math.Exp2(30)), 2, uint32(math.Exp(32)), nil},
    }
    
    for idx, i := range tests {
        var result, err = uint32BinaryCheckedShl(i.a, i.b)
        if err == i.expectedError {
            // pass
        } else if err != nil {
            t.Errorf("test %d: unexpected error %v (expected %v)", idx, err, i.expectedError)
        } else if result != i.expectedValue {
            t.Errorf("test %d: got %d, but expected %d", idx, result, i.expectedValue)
        }
    }
}

func TestIntCheckedAdd(t *testing.T) {
    type test struct {
        a int8
        b int8
        expectedValue int8
        expectedError error
    }
    
    var tests = []test {
        {  -1,   -1,   -2, nil},
        {   1,   -1,    0, nil},
        {   0,    0,    0, nil},
        {   8,    4,   12, nil},
        {  64,   63,  127, nil},
        {  64,   63,  127, nil},
        { -64,  -64, -128, nil},
        {  64,   64,    0, ErrorOverflow},
        { -64,  -65,    0, ErrorOverflow},
    }
    
    for idx, i := range tests {
        var result, err = int8BinaryCheckedAdd(i.a, i.b)
        if err == i.expectedError {
            // pass
        } else if err != nil {
            t.Errorf("test %d: unexpected error %v (expected %v)", idx, err, i.expectedError)
        } else if result != i.expectedValue {
            t.Errorf("test %d: got %d, but expected %d", idx, result, i.expectedValue)
        }
    }
}

func TestUintCheckedAdd(t *testing.T) {
    type test struct {
        a uint8
        b uint8
        expectedValue uint8
        expectedError error
    }
    
    var tests = []test {
        {   0,    0,    0, nil},
        {   8,    4,   12, nil},
        { 128,  127,  255, nil},
        { 128,  128,    0, ErrorOverflow},
    }
    
    for idx, i := range tests {
        var result, err = uint8BinaryCheckedAdd(i.a, i.b)
        if err == i.expectedError {
            // pass
        } else if err != nil {
            t.Errorf("test %d: unexpected error %v (expected %v)", idx, err, i.expectedError)
        } else if result != i.expectedValue {
            t.Errorf("test %d: got %d, but expected %d", idx, result, i.expectedValue)
        }
    }
}

func TestIntNary(t *testing.T) {
    type test struct {
        xs []int
        f func(...int) int
        expected int
    }
    
    var tests = []test {
        {[]int{1, 2, 3, 4, 5, 6}, Int.Nary.Add, 21},
        {[]int{1, 2, 3, 4, 5, 6}, Int.Nary.Mul, 720},
    }
    
    for idx, i := range tests {
        var result = i.f(i.xs...)
        if result != i.expected {
            t.Errorf("test %d: got %d, but expected %d", idx, result, i.expected)
        }
    }
}

func TestIntNaryChecked(t *testing.T) {
    type test struct {
        xs []int8
        f func(...int8) (int8, error)
        expectedValue int8
        expectedError error
    }
    
    var tests = []test {
        {[]int8{1, 2, 3, 4, 5, 6},          Int8Checked.Nary.Add, 21, nil},
        {[]int8{1, 2, 3, 4},                Int8Checked.Nary.Mul, 24, nil},
        {[]int8{120, -120, -120, -8, -1},   Int8Checked.Nary.Add,  0, ErrorOverflow},
        {[]int8{120, 4, 4},                 Int8Checked.Nary.Add,  0, ErrorOverflow},
        {[]int8{32, 2, 2},                  Int8Checked.Nary.Mul,  0, ErrorOverflow},
    }
    
    for idx, i := range tests {
        var result, err = i.f(i.xs...)
        if err == i.expectedError {
            // pass
        } else if err != nil {
            t.Errorf("test %d: unexpected error %v (expected %v)", idx, err, i.expectedError)
        } else if result != i.expectedValue {
            t.Errorf("test %d: got %d, but expected %d", idx, result, i.expectedValue)
        }
    }
}

func TestIntReduce(t *testing.T) {
    type test struct {
        operatorIdentity int
        xs []int
        f func(int, int) int
        expected int
    }
    
    var tests = []test {
        {0, []int{1, 2, 3, 4, 5, 6}, Int.Binary.Add, 21},
        {1, []int{1, 2, 3, 4, 5, 6}, Int.Binary.Mul, 720},
    }
    
    for idx, i := range tests {
        var result = Int.Reduce(i.operatorIdentity, i.f, i.xs...)
        if result != i.expected {
            t.Errorf("test %d: got %d, but expected %d", idx, result, i.expected)
        }
    }
}

func TestIntCheckedReduce(t *testing.T) {
    type test struct {
        operatorIdentity int8
        xs []int8
        f func(int8, int8) (int8, error)
        expectedValue int8
        expectedError error
    }
    
    var tests = []test {
        {0, []int8{1, 2, 3, 4, 5, 6},          Int8Checked.Binary.Add, 21, nil},
        {1, []int8{1, 2, 3, 4},                Int8Checked.Binary.Mul, 24, nil},
        {0, []int8{120, -120, -120, -8, -1},   Int8Checked.Binary.Add,  0, ErrorOverflow},
        {0, []int8{120, 4, 4},                 Int8Checked.Binary.Add,  0, ErrorOverflow},
        {1, []int8{32, 2, 2},                  Int8Checked.Binary.Mul,  0, ErrorOverflow},
    }
    
    for idx, i := range tests {
        var result, err = Int8Checked.Reduce(i.operatorIdentity, i.f, i.xs...)
        if err == i.expectedError {
            // pass
        } else if err != nil {
            t.Errorf("test %d: unexpected error %v (expected %v)", idx, err, i.expectedError)
        } else if result != i.expectedValue {
            t.Errorf("test %d: got %d, but expected %d", idx, result, i.expectedValue)
        }
    }
}
