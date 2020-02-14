package operator
// IGNORE
// This file is templated by template-numbers.py // IGNORE
type _t int // IGNORE
const max_T =  1000 // IGNORE
const min_T = -1000 // IGNORE

// Overflow checks with reference to https://stackoverflow.com/a/1514309/5654201

type _tBinary struct {
    Add             func(_t, _t) _t
    Sub             func(_t, _t) _t
    Mul             func(_t, _t) _t
    Div             func(_t, _t) _t
    Mod             func(_t, _t) _t
    
    And             func(_t, _t) _t
    Or              func(_t, _t) _t
    Xor             func(_t, _t) _t
    AndNot          func(_t, _t) _t
    
    Shl             func(_t, uint) _t
    Shr             func(_t, uint) _t
}

type _tBinaryChecked struct {
    Add             func(_t, _t) (_t, error)
    Sub             func(_t, _t) (_t, error)
    Mul             func(_t, _t) (_t, error)
    
    Shl             func(_t, uint) (_t, error)
    Shr             func(_t, uint) (_t, error)
}

var _T = struct {
    Binary          _tBinary
}{
    Binary:          _tBinary{
        Add:        func(a _t, b _t) _t { return a + b }, // CLASS integers, floats, complex
        Sub:        func(a _t, b _t) _t { return a - b }, // CLASS integers, floats, complex
        Mul:        func(a _t, b _t) _t { return a * b }, // CLASS integers, floats, complex
        Div:        func(a _t, b _t) _t { return a / b }, // CLASS integers, floats, complex
        
        And:        func(a _t, b _t) _t { return a & b },  // CLASS integers
        Or:         func(a _t, b _t) _t { return a | b },  // CLASS integers
        Xor:        func(a _t, b _t) _t { return a ^ b },  // CLASS integers
        AndNot:     func(a _t, b _t) _t { return a &^ b }, // CLASS integers
        Mod:        func(a _t, b _t) _t { return a % b },  // CLASS integers
        
        Shl:        func(a _t, b uint) _t { return a << b }, // CLASS integers
        Shr:        func(a _t, b uint) _t { return a >> b }, // CLASS integers
    },
}

var _TChecked = struct {
    Binary          _tBinaryChecked
}{
    Binary:         _tBinaryChecked{
        Add:        _tBinaryCheckedAdd, // CLASS integers, floats, complex
        Sub:        _tBinaryCheckedSub, // CLASS integers, floats, complex
        Mul:        _tBinaryCheckedMul, // CLASS integers, floats, complex
        Shl:        _tBinaryCheckedShl, // CLASS integers
    },
}

func _tBinaryCheckedAdd(a _t, b _t) (v _t, err error) {
    if (b > 0) && (a > (max_T - b)) { return v, ErrorOverflow }
    if (b < 0) && (a < (min_T - b)) { return v, ErrorOverflow }
    return a + b, nil
}

func _tBinaryCheckedSub(a _t, b _t) (v _t, err error) {
    if (b < 0) && (a > (max_T + b)) { return v, ErrorOverflow }
    if (b > 0) && (a < (min_T + b)) { return v, ErrorOverflow }
    return a - b, nil
}

func _tBinaryCheckedMul(a _t, b _t) (v _t, err error) {
    if (a == -1) && (b == min_T) { return v, ErrorOverflow } // CLASS integers; signed
    if (b == -1) && (a == min_T) { return v, ErrorOverflow } // CLASS integers; signed
    if (a > (max_T / b)) { return v, ErrorOverflow }
    if (a < (min_T / b)) { return v, ErrorOverflow }
    
    return a * b, nil
}

func _tBinaryCheckedShl(a _t, b uint) (v _t, err error) { // CLASS integers
    if a < 0 { return v, ErrorUndefined } // CLASS integers; signed
    if b > uint(_tMostSignificantBit(max_T)) { return v, ErrorOverflow } // CLASS integers
    return v, err // CLASS integers
} // CLASS integers

func _tMostSignificantBit(a _t) (result int) { // CLASS integers
  for a > 0 { // CLASS integers
      a >>= 1 // CLASS integers
      result++ // CLASS integers
  } // CLASS integers
  return result; // CLASS integers
} // CLASS integers
