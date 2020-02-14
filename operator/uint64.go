package operator

// Code generated by (tawesoft.co.uk/go/operator) template-numbers.py: DO NOT EDIT.


// Some overflow checks with reference to stackoverflow.com/a/1514309/5654201

type uint64Unary struct {
    Identity        func(uint64) uint64
    Not             func(uint64) uint64
    Zero            func(uint64) bool
    NonZero         func(uint64) bool
    Positive        func(uint64) bool
    Negative        func(uint64) bool
}

type uint64UnaryChecked struct {
}

type uint64Binary struct {
    Add             func(uint64, uint64) uint64
    Sub             func(uint64, uint64) uint64
    Mul             func(uint64, uint64) uint64
    Div             func(uint64, uint64) uint64
    Mod             func(uint64, uint64) uint64
    
    Eq              func(uint64, uint64) bool
    Neq             func(uint64, uint64) bool
    Lt              func(uint64, uint64) bool
    Lte             func(uint64, uint64) bool
    Gt              func(uint64, uint64) bool
    Gte             func(uint64, uint64) bool
    
    And             func(uint64, uint64) uint64
    Or              func(uint64, uint64) uint64
    Xor             func(uint64, uint64) uint64
    AndNot          func(uint64, uint64) uint64
    
    Shl             func(uint64, uint) uint64
    Shr             func(uint64, uint) uint64
}

type uint64BinaryChecked struct {
    Add             func(uint64, uint64) (uint64, error)
    Sub             func(uint64, uint64) (uint64, error)
    Mul             func(uint64, uint64) (uint64, error)
    Div             func(uint64, uint64) (uint64, error)
    
    Shl             func(uint64, uint) (uint64, error)
    Shr             func(uint64, uint) (uint64, error)
}

type uint64Nary struct {
    Add             func(... uint64) uint64
    Sub             func(... uint64) uint64
    Mul             func(... uint64) uint64
}

type uint64NaryChecked struct {
    Add             func(... uint64) (uint64, error)
    Sub             func(... uint64) (uint64, error)
    Mul             func(... uint64) (uint64, error)
}

// Uint64 implements operations on one (unary), two (binary), or many (nary) arguments of type uint64.
var Uint64 = struct {
    Unary           uint64Unary
    Binary          uint64Binary
    Nary            uint64Nary
    Reduce          func(operatorIdentity uint64, operator func(uint64, uint64) uint64, elements ... uint64) uint64
}{
    Unary:          uint64Unary{
        Identity:   func(a uint64) uint64 { return a },
        Not:        func(a uint64) uint64 { return ^a },
        Zero:       func(a uint64) bool { return a == 0 },
        NonZero:    func(a uint64) bool { return a != 0 },
        Positive:   uint64UnaryPositive,
        Negative:   uint64UnaryNegative,
    },
    
    Binary:          uint64Binary{
        Add:        func(a uint64, b uint64) uint64 { return a + b },
        Sub:        func(a uint64, b uint64) uint64 { return a - b },
        Mul:        func(a uint64, b uint64) uint64 { return a * b },
        Div:        func(a uint64, b uint64) uint64 { return a / b },
        
        Eq:         func(a uint64, b uint64) bool { return a == b },
        Neq:        func(a uint64, b uint64) bool { return a != b },
        Lt:         func(a uint64, b uint64) bool { return a <  b },
        Lte:        func(a uint64, b uint64) bool { return a <= b },
        Gt:         func(a uint64, b uint64) bool { return a >  b },
        Gte:        func(a uint64, b uint64) bool { return a >= b },
        
        And:        func(a uint64, b uint64) uint64 { return a & b },
        Or:         func(a uint64, b uint64) uint64 { return a | b },
        Xor:        func(a uint64, b uint64) uint64 { return a ^ b },
        AndNot:     func(a uint64, b uint64) uint64 { return a &^ b },
        Mod:        func(a uint64, b uint64) uint64 { return a % b },
        
        Shl:        func(a uint64, b uint) uint64 { return a << b },
        Shr:        func(a uint64, b uint) uint64 { return a >> b },
    },
    
    Nary:           uint64Nary{
        Add:        uint64NaryAdd,
        Mul:        uint64NaryMul,
    },
    
    Reduce:         uint64Reduce,
}

// Uint64Checked implements operations on one (unary), two (binary), or many (nary) arguments of type uint64, returning an
// error in cases such as overflow or an undefined operation.
var Uint64Checked = struct {
    Unary           uint64UnaryChecked
    Binary          uint64BinaryChecked
    Nary            uint64NaryChecked
    Reduce          func(operatorIdentity uint64, operator func(uint64, uint64) (uint64, error), elements ... uint64) (uint64, error)
}{
    Unary:          uint64UnaryChecked{
    },
    
    Binary:         uint64BinaryChecked{
        Add:        uint64BinaryCheckedAdd,
        Sub:        uint64BinaryCheckedSub,
        Mul:        uint64BinaryCheckedMul,
        Div:        uint64BinaryCheckedDiv,
        Shl:        uint64BinaryCheckedShl,
    },
    
    Nary:           uint64NaryChecked{
        Add:        uint64NaryCheckedAdd,
        Mul:        uint64NaryCheckedMul,
    },
    
    Reduce:         uint64CheckedReduce,
}

func uint64UnaryPositive(a uint64) bool {
    return a > 0
}

func uint64UnaryNegative(a uint64) bool {
    return a < 0
}




func uint64BinaryCheckedAdd(a uint64, b uint64) (v uint64, err error) {
    if (b > 0) && (a > (maxUint64 - b)) { return v, ErrorOverflow }
    if (b < 0) && (a < (minUint64 - b)) { return v, ErrorOverflow }
    return a + b, nil
}

func uint64BinaryCheckedSub(a uint64, b uint64) (v uint64, err error) {
    if (b < 0) && (a > (maxUint64 + b)) { return v, ErrorOverflow }
    if (b > 0) && (a < (minUint64 + b)) { return v, ErrorOverflow }
    return a - b, nil
}

func uint64BinaryCheckedMul(a uint64, b uint64) (v uint64, err error) {
    if (a > (maxUint64 / b)) { return v, ErrorOverflow }
    if (a < (minUint64 / b)) { return v, ErrorOverflow }
    
    return a * b, nil
}

func uint64BinaryCheckedDiv(a uint64, b uint64) (v uint64, err error) {
    if (b == 0) { return v, ErrorUndefined }
    
    return a / b, nil
}

func uint64BinaryCheckedShl(a uint64, b uint) (v uint64, err error) {
    if b > uint(uint64MostSignificantBit(maxUint64)) { return v, ErrorOverflow }
    return v, err
}

func uint64MostSignificantBit(a uint64) (result int) {
  for a > 0 {
      a >>= 1
      result++
  }
  return result;
}

func uint64NaryAdd(xs ... uint64) (result uint64) {
    for i := 0; i < len(xs); i++ {
        result += xs[i]
    }
    return result
}

func uint64NaryCheckedAdd(xs ... uint64) (result uint64, err error) {
    for i := 0; i < len(xs); i++ {
        result, err = uint64BinaryCheckedAdd(result, xs[i])
        if err != nil { return result, err }
    }
    return result, nil
}

func uint64NaryMul(xs ... uint64) (result uint64) {
    result = 1
    for i := 0; i < len(xs); i++ {
        result *= xs[i]
    }
    return result
}

func uint64NaryCheckedMul(xs ... uint64) (result uint64, err error) {
    result = 1
    for i := 0; i < len(xs); i++ {
        result, err = uint64BinaryCheckedMul(result, xs[i])
        if err != nil { return result, err }
    }
    return result, nil
}

func uint64Reduce(operatorIdentity uint64, operator func(uint64, uint64) uint64, elements ... uint64) (result uint64) {
    result = operatorIdentity
    for i := 0; i < len(elements); i++ {
        result = operator(result, elements[i])
    }
    return result
}

func uint64CheckedReduce(operatorIdentity uint64, operator func(uint64, uint64) (uint64, error), elements ... uint64) (result uint64, err error) {
    result = operatorIdentity
    for i := 0; i < len(elements); i++ {
        result, err = operator(result, elements[i])
        if err != nil { return result, err }
    }
    return result, err
}

