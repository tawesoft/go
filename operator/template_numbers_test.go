package operator

import "math" // CLASS floats
// IGNORE
// This file is templated by template-numbers.py // IGNORE
// The token _t is replaced by a snakeCase type // IGNORE
// The token _T is replaced by a CamelCase type // IGNORE
// And lines are conditionally filtered by a command or command; options comment // IGNORE
type _t int // IGNORE
const max_T =  1000 // IGNORE
const min_T = -1000 // IGNORE

// Some overflow checks with reference to stackoverflow.com/a/1514309/5654201

type _tUnary struct {
    Identity        func(_t) _t
    Abs             func(_t) _t // CLASS integers, floats; signed
    Negation        func(_t) _t // CLASS integers, floats; signed
    Not             func(_t) _t // CLASS integers; unsigned
    Zero            func(_t) bool
    NonZero         func(_t) bool
    Positive        func(_t) bool
    Negative        func(_t) bool
}

type _tUnaryChecked struct {
    Abs             func(_t) (_t, error) // CLASS integers, floats; signed
    Negation        func(_t) (_t, error) // CLASS integers, floats; signed
}

type _tBinary struct {
    Add             func(_t, _t) _t // CLASS integers, floats, complex
    Sub             func(_t, _t) _t // CLASS integers, floats, complex
    Mul             func(_t, _t) _t // CLASS integers, floats, complex
    Div             func(_t, _t) _t // CLASS integers, floats, complex
    Mod             func(_t, _t) _t // CLASS integers, floats, complex
    
    Eq              func(_t, _t) bool // CLASS integers, floats, complex
    Neq             func(_t, _t) bool // CLASS integers, floats, complex
    Lt              func(_t, _t) bool // CLASS integers, floats, complex
    Lte             func(_t, _t) bool // CLASS integers, floats, complex
    Gt              func(_t, _t) bool // CLASS integers, floats, complex
    Gte             func(_t, _t) bool // CLASS integers, floats, complex
    
    And             func(_t, _t) _t // CLASS integers
    Or              func(_t, _t) _t // CLASS integers
    Xor             func(_t, _t) _t // CLASS integers
    AndNot          func(_t, _t) _t // CLASS integers
    
    Shl             func(_t, uint) _t // CLASS integers
    Shr             func(_t, uint) _t // CLASS integers
}

type _tBinaryChecked struct {
    Add             func(_t, _t) (_t, error) // CLASS integers, floats, complex
    Sub             func(_t, _t) (_t, error) // CLASS integers, floats, complex
    Mul             func(_t, _t) (_t, error) // CLASS integers, floats, complex
    Div             func(_t, _t) (_t, error) // CLASS integers, floats, complex
    
    Shl             func(_t, uint) (_t, error) // CLASS integers
    Shr             func(_t, uint) (_t, error) // CLASS integers
}

type _tNary struct {
    Add             func(... _t) _t // CLASS integers, floats, complex
    Sub             func(... _t) _t // CLASS integers, floats, complex
    Mul             func(... _t) _t // CLASS integers, floats, complex
}

type _tNaryChecked struct {
    Add             func(... _t) (_t, error) // CLASS integers, floats, complex
    Sub             func(... _t) (_t, error) // CLASS integers, floats, complex
    Mul             func(... _t) (_t, error) // CLASS integers, floats, complex
}

// _T implements operations on one (unary), two (binary), or many (nary) arguments of type _t.
var _T = struct {
    Unary           _tUnary
    Binary          _tBinary
    Nary            _tNary
    Reduce          func(operatorIdentity _t, operator func(_t, _t) _t, elements ... _t) _t
}{
    Unary:          _tUnary{
        Identity:   func(a _t) _t { return a },
        Abs:        _tUnaryAbs,      // CLASS integers, floats; signed
        Negation:   func(a _t) _t { return -a }, // CLASS integers, floats; signed
        Not:        func(a _t) _t { return ^a }, // CLASS integers; unsigned
        Zero:       func(a _t) bool { return a == 0 },
        NonZero:    func(a _t) bool { return a != 0 },
        Positive:   _tUnaryPositive,
        Negative:   _tUnaryNegative,
    },
    
    Binary:          _tBinary{
        Add:        func(a _t, b _t) _t { return a + b }, // CLASS integers, floats, complex
        Sub:        func(a _t, b _t) _t { return a - b }, // CLASS integers, floats, complex
        Mul:        func(a _t, b _t) _t { return a * b }, // CLASS integers, floats, complex
        Div:        func(a _t, b _t) _t { return a / b }, // CLASS integers, floats, complex
        
        Eq:         func(a _t, b _t) bool { return a == b }, // CLASS integers, floats, complex
        Neq:        func(a _t, b _t) bool { return a != b }, // CLASS integers, floats, complex
        Lt:         func(a _t, b _t) bool { return a <  b }, // CLASS integers, floats, complex
        Lte:        func(a _t, b _t) bool { return a <= b }, // CLASS integers, floats, complex
        Gt:         func(a _t, b _t) bool { return a >  b }, // CLASS integers, floats, complex
        Gte:        func(a _t, b _t) bool { return a >= b }, // CLASS integers, floats, complex
        
        And:        func(a _t, b _t) _t { return a & b },  // CLASS integers
        Or:         func(a _t, b _t) _t { return a | b },  // CLASS integers
        Xor:        func(a _t, b _t) _t { return a ^ b },  // CLASS integers
        AndNot:     func(a _t, b _t) _t { return a &^ b }, // CLASS integers
        Mod:        func(a _t, b _t) _t { return a % b },  // CLASS integers
        
        Shl:        func(a _t, b uint) _t { return a << b }, // CLASS integers
        Shr:        func(a _t, b uint) _t { return a >> b }, // CLASS integers
    },
    
    Nary:           _tNary{
        Add:        _tNaryAdd, // CLASS integers, floats, complex
        Mul:        _tNaryMul, // CLASS integers, floats, complex
    },
    
    Reduce:         _tReduce,
}

// _TChecked implements operations on one (unary), two (binary), or many (nary) arguments of type _t, returning an
// error in cases such as overflow or an undefined operation.
var _TChecked = struct {
    Unary           _tUnaryChecked
    Binary          _tBinaryChecked
    Nary            _tNaryChecked
    Reduce          func(operatorIdentity _t, operator func(_t, _t) (_t, error), elements ... _t) (_t, error)
}{
    Unary:          _tUnaryChecked{
        Abs:        _tUnaryCheckedAbs,      // CLASS integers, floats; signed
        Negation:   _tUnaryCheckedNegation, // CLASS integers, floats; signed
    },
    
    Binary:         _tBinaryChecked{
        Add:        _tBinaryCheckedAdd, // CLASS integers, floats, complex
        Sub:        _tBinaryCheckedSub, // CLASS integers, floats, complex
        Mul:        _tBinaryCheckedMul, // CLASS integers, floats, complex
        Div:        _tBinaryCheckedDiv, // CLASS integers, floats, complex
        Shl:        _tBinaryCheckedShl, // CLASS integers
    },
    
    Nary:           _tNaryChecked{
        Add:        _tNaryCheckedAdd, // CLASS integers, floats, complex
        Mul:        _tNaryCheckedMul, // CLASS integers, floats, complex
    },
    
    Reduce:         _tCheckedReduce,
}

func _tUnaryPositive(a _t) bool {
    return math.Signbit(float64(a)) == false // CLASS floats
    return a > 0 // CLASS integers
}

func _tUnaryNegative(a _t) bool {
    return math.Signbit(float64(a)) == true // CLASS floats
    return a < 0 // CLASS integers
}

func _tUnaryAbs(a _t) _t { // CLASS integers, floats; signed
    return _t(math.Abs(float64(a))) // CLASS floats
    if a < 0 { return -a } // CLASS integers; signed
    return a // CLASS integers; signed
} // CLASS integers, floats; signed

// note abs(+/- Inf) = +Inf // CLASS floats
func _tUnaryCheckedAbs(a _t) (v _t, err error) { // CLASS integers, floats; signed
    if a == min_T { return v, ErrorOverflow } // CLASS integers; signed
    if math.IsNaN(float64(a)) { return v, ErrorNaN } // CLASS floats
    return _t(math.Abs(float64(a))), nil // CLASS floats
    if a < 0 { return -a, nil } // CLASS integers; signed
    return a, nil // CLASS integers; signed
} // CLASS integers, floats; signed

func _tUnaryCheckedNegation(a _t) (v _t, err error) { // CLASS integers, floats; signed
    if (a == min_T) { return v, ErrorOverflow } // CLASS integers; signed
    if math.IsNaN(float64(a)) { return v, ErrorNaN } // CLASS floats
    return -a, nil // CLASS integers, floats; signed
} // CLASS integers, floats; signed

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

func _tBinaryCheckedDiv(a _t, b _t) (v _t, err error) {
    if math.IsNaN(float64(a)) { return v, ErrorNaN } // CLASS floats
    if (b == -1) && (a == min_T) { return v, ErrorOverflow } // CLASS integers; signed
    if (b == 0) { return v, ErrorUndefined }
    
    return a / b, nil
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

func _tNaryAdd(xs ... _t) (result _t) {
    for i := 0; i < len(xs); i++ {
        result += xs[i]
    }
    return result
}

func _tNaryCheckedAdd(xs ... _t) (result _t, err error) {
    for i := 0; i < len(xs); i++ {
        result, err = _tBinaryCheckedAdd(result, xs[i])
        if err != nil { return result, err }
    }
    return result, nil
}

func _tNaryMul(xs ... _t) (result _t) {
    result = 1
    for i := 0; i < len(xs); i++ {
        result *= xs[i]
    }
    return result
}

func _tNaryCheckedMul(xs ... _t) (result _t, err error) {
    result = 1
    for i := 0; i < len(xs); i++ {
        result, err = _tBinaryCheckedMul(result, xs[i])
        if err != nil { return result, err }
    }
    return result, nil
}

func _tReduce(operatorIdentity _t, operator func(_t, _t) _t, elements ... _t) (result _t) {
    result = operatorIdentity
    for i := 0; i < len(elements); i++ {
        result = operator(result, elements[i])
    }
    return result
}

func _tCheckedReduce(operatorIdentity _t, operator func(_t, _t) (_t, error), elements ... _t) (result _t, err error) {
    result = operatorIdentity
    for i := 0; i < len(elements); i++ {
        result, err = operator(result, elements[i])
        if err != nil { return result, err }
    }
    return result, err
}

