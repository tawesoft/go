package operator

// Code generated by (tawesoft.co.uk/go/operator) template-numbers.py: DO NOT EDIT.



type uintUnary struct {
    Identity        func(uint) uint
    Not             func(uint) uint
    Zero            func(uint) bool
    NonZero         func(uint) bool
    Positive        func(uint) bool
    Negative        func(uint) bool
}

type uintUnaryChecked struct {
}

type uintBinary struct {
    Add             func(uint, uint) uint
    Sub             func(uint, uint) uint
    Mul             func(uint, uint) uint
    Div             func(uint, uint) uint
    Mod             func(uint, uint) uint
    
    Eq              func(uint, uint) bool
    Neq             func(uint, uint) bool
    Lt              func(uint, uint) bool
    Lte             func(uint, uint) bool
    Gt              func(uint, uint) bool
    Gte             func(uint, uint) bool
    
    And             func(uint, uint) uint
    Or              func(uint, uint) uint
    Xor             func(uint, uint) uint
    AndNot          func(uint, uint) uint
    
    Shl             func(uint, uint) uint
    Shr             func(uint, uint) uint
}

type uintBinaryChecked struct {
    Add             func(uint, uint) (uint, error)
    Sub             func(uint, uint) (uint, error)
    Mul             func(uint, uint) (uint, error)
    
    Shl             func(uint, uint) (uint, error)
    Shr             func(uint, uint) (uint, error)
}

type uintNary struct {
    Add             func(... uint) uint
    Sub             func(... uint) uint
    Mul             func(... uint) uint
}

type uintNaryChecked struct {
    Add             func(... uint) (uint, error)
    Sub             func(... uint) (uint, error)
    Mul             func(... uint) (uint, error)
}

var Uint = struct {
    Unary           uintUnary
    Binary          uintBinary
    Nary            uintNary
}{
    Unary:          uintUnary{
        Identity:   func(a uint) uint { return a },
        Not:        func(a uint) uint { return ^a },
        Zero:       func(a uint) bool { return a == 0 },
        NonZero:    func(a uint) bool { return a != 0 },
        Positive:   uintUnaryPositive,
        Negative:   uintUnaryNegative,
    },
    
    Binary:          uintBinary{
        Add:        func(a uint, b uint) uint { return a + b },
        Sub:        func(a uint, b uint) uint { return a - b },
        Mul:        func(a uint, b uint) uint { return a * b },
        Div:        func(a uint, b uint) uint { return a / b },
        
        Eq:         func(a uint, b uint) bool { return a == b },
        Neq:        func(a uint, b uint) bool { return a != b },
        Lt:         func(a uint, b uint) bool { return a <  b },
        Lte:        func(a uint, b uint) bool { return a <= b },
        Gt:         func(a uint, b uint) bool { return a >  b },
        Gte:        func(a uint, b uint) bool { return a >= b },
        
        And:        func(a uint, b uint) uint { return a & b },
        Or:         func(a uint, b uint) uint { return a | b },
        Xor:        func(a uint, b uint) uint { return a ^ b },
        AndNot:     func(a uint, b uint) uint { return a &^ b },
        Mod:        func(a uint, b uint) uint { return a % b },
        
        Shl:        func(a uint, b uint) uint { return a << b },
        Shr:        func(a uint, b uint) uint { return a >> b },
    },
    
    Nary:           uintNary{
        Add:        uintNaryAdd,
        Mul:        uintNaryMul,
    },
}

var UintChecked = struct {
    Unary           uintUnaryChecked
    Binary          uintBinaryChecked
    Nary            uintNaryChecked
}{
    Unary:          uintUnaryChecked{
    },
    
    Binary:         uintBinaryChecked{
        Add:        uintBinaryCheckedAdd,
        Sub:        uintBinaryCheckedSub,
        Mul:        uintBinaryCheckedMul,
        Shl:        uintBinaryCheckedShl,
    },
    
    Nary:           uintNaryChecked{
        Add:        uintNaryCheckedAdd,
        Mul:        uintNaryCheckedMul,
    },
}

func uintUnaryPositive(a uint) bool {
    return a > 0
}

func uintUnaryNegative(a uint) bool {
    return a < 0
}




func uintBinaryCheckedAdd(a uint, b uint) (v uint, err error) {
    if (b > 0) && (a > (maxUint - b)) { return v, ErrorOverflow }
    if (b < 0) && (a < (minUint - b)) { return v, ErrorOverflow }
    return a + b, nil
}

func uintBinaryCheckedSub(a uint, b uint) (v uint, err error) {
    if (b < 0) && (a > (maxUint + b)) { return v, ErrorOverflow }
    if (b > 0) && (a < (minUint + b)) { return v, ErrorOverflow }
    return a - b, nil
}

func uintBinaryCheckedMul(a uint, b uint) (v uint, err error) {
    if (a > (maxUint / b)) { return v, ErrorOverflow }
    if (a < (minUint / b)) { return v, ErrorOverflow }
    
    return a * b, nil
}

func uintBinaryCheckedShl(a uint, b uint) (v uint, err error) {
    if b > uint(uintMostSignificantBit(maxUint)) { return v, ErrorOverflow }
    return v, err
}

func uintMostSignificantBit(a uint) (result int) {
  for a > 0 {
      a >>= 1
      result++
  }
  return result;
}

func uintNaryAdd(xs ... uint) (result uint) {
    for i := 0; i < len(xs); i++ {
        result += xs[i]
    }
    return result
}

func uintNaryCheckedAdd(xs ... uint) (result uint, err error) {
    for i := 0; i < len(xs); i++ {
        result, err = uintBinaryCheckedAdd(result, xs[i])
        if err != nil { return result, err }
    }
    return result, nil
}

func uintNaryMul(xs ... uint) (result uint) {
    result = 1
    for i := 0; i < len(xs); i++ {
        result *= xs[i]
    }
    return result
}

func uintNaryCheckedMul(xs ... uint) (result uint, err error) {
    result = 1
    for i := 0; i < len(xs); i++ {
        result, err = uintBinaryCheckedMul(result, xs[i])
        if err != nil { return result, err }
    }
    return result, nil
}

