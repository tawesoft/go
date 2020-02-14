package operator

type boolUnary struct {
    True            func(bool) bool
    False           func(bool) bool
    Identity        func(bool) bool
    Not             func(bool) bool
}

type boolBinary struct{
    True            func(bool, bool) bool // always returns true
    False           func(bool, bool) bool // always returns false
    P               func(bool, bool) bool // for (p, q) returns p
    Q               func(bool, bool) bool // for (p, q) returns q
    NotP            func(bool, bool) bool // for (p, q) returns ¬ p
    NotQ            func(bool, bool) bool // for (p, q) returns ¬ q
    Eq              func(bool, bool) bool // same as XNOR
    Neq             func(bool, bool) bool // same as XOR
    And             func(bool, bool) bool
    Nand            func(bool, bool) bool
    Or              func(bool, bool) bool
    Nor             func(bool, bool) bool
    Xor             func(bool, bool) bool // same as Neq
    Xnor            func(bool, bool) bool // same as Eq
    Implies         func(bool, bool) bool // for (p, q) returns p -> q
    NonImplies      func(bool, bool) bool // non implication; for (p, q) returns ¬ (p -> q)
    ConImplies      func(bool, bool) bool // converse implies; for (p, q) returns (q -> p)
    ConNonImplies   func(bool, bool) bool // converse not implies; for (p, q) returns ¬ (q -> p)
}

type boolNary struct{
    True            func(...bool) bool // always returns true
    False           func(...bool) bool // always returns false
    
    All             func(...bool) bool // returns true if all are true a.k.a. a AND b AND c ...
    Any             func(...bool) bool // returns true if any are true a.k.a. a OR b OR c
    None            func(...bool) bool // returns true if none are true (all are false)
}

// Bool implements operations on one (unary), two (binary), or many (nary) arguments of type bool.
var Bool = struct {
    Unary  boolUnary
    Binary boolBinary
    Nary   boolNary
}{
    Unary: boolUnary{
        True:       func(_ bool) bool { return true},
        False:      func(_ bool) bool { return false},
        Identity:   func(p bool) bool { return p },
        Not:        func(p bool) bool { return !p },
    },
    
    Binary: boolBinary{
        True:       func(_ bool, _ bool) bool { return true },
        False:      func(_ bool, _ bool) bool { return false },
        P:          func(p bool, _ bool) bool { return p },
        Q:          func(_ bool, q bool) bool { return q },
        NotP:       func(p bool, _ bool) bool { return !p },
        NotQ:       func(_ bool, q bool) bool { return !q },
        Eq:         func(p bool, q bool) bool { return p == q },
        Neq:        func(p bool, q bool) bool { return p != q },
        And:        func(p bool, q bool) bool { return p && q },
        Nand:       func(p bool, q bool) bool { return !(p && q) },
        Or:         func(p bool, q bool) bool { return p || q },
        Nor:        func(p bool, q bool) bool { return !(p || q) },
        Xor:        func(p bool, q bool) bool { return p != q },
        Xnor:       func(p bool, q bool) bool { return p == q },
        Implies:    func(p bool, q bool) bool { return (!p) || q },
        NonImplies: func(p bool, q bool) bool { return p && (!q) },
        ConImplies: func(p bool, q bool) bool { return (!q) || p },
        ConNonImplies: func(p bool, q bool) bool { return q && (!p) },
    },
    
    Nary: boolNary{
        True:       func(_...bool) bool { return true },
        False:      func(_...bool) bool { return false },
        All:        boolNaryAll1,
        Any:        boolNaryAny1,
        None:       boolNaryNone1,
    },
}

func boolNaryAll1(p...bool) bool {
    for i := 0; i < len(p); i++ {
        if !p[i] { return false }
    }
    return true
}

func boolNaryAll2(p...bool) (result bool) {
    result = true
    for i := 0; i < len(p); i++ {
        result = result && p[i]
    }
    return result
}

func boolNaryAny1(p...bool) bool {
    for i := 0; i < len(p); i++ {
        if p[i] { return true }
    }
    return false
}

func boolNaryAny2(p...bool) (result bool) {
    for i := 0; i < len(p); i++ {
        result = result || p[i]
    }
    return result
}

func boolNaryNone1(p...bool) bool {
    for i := 0; i < len(p); i++ {
        if p[i] { return false }
    }
    return true
}

func boolNaryNone2(p...bool) (result bool) {
    result = false
    for i := 0; i < len(p); i++ {
        result = result || p[i]
    }
    return !result
}

