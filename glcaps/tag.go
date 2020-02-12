package glcaps

import (
    "fmt"
    "strconv"
    "strings"
)

func operationBooleanAnd(a bool, b bool) bool { return a && b }
func operationBooleanOr (a bool, b bool) bool { return a || b }
func operationIntEq     (a int,  b int)  bool { return a == b }
func operationIntNeq    (a int,  b int)  bool { return a != b }
func operationIntLt     (a int,  b int)  bool { return a <  b }
func operationIntLte    (a int,  b int)  bool { return a <= b }
func operationIntGt     (a int,  b int)  bool { return a >  b }
func operationIntGte    (a int,  b int)  bool { return a >= b }
func operationFloat32Eq (a float32,  b float32)  bool { return a == b }
func operationFloat32Neq(a float32,  b float32)  bool { return a != b }
func operationFloat32Lt (a float32,  b float32)  bool { return a <  b }
func operationFloat32Lte(a float32,  b float32)  bool { return a <= b }
func operationFloat32Gt (a float32,  b float32)  bool { return a >  b }
func operationFloat32Gte(a float32,  b float32)  bool { return a >= b }

type command interface{
    evalBool (b *Binding, extensions Extensions) bool
    evalInt  (b *Binding, extensions Extensions) int
    evalFloat(b *Binding, extensions Extensions) float32
    hasBoolRepresentation() bool
    hasIntRepresentation() bool
    hasFloatRepresentation() bool
}

type requirement interface{
    evalBool (field string, result bool)    error
    evalInt  (field string, result int)     error
    evalFloat(field string, result float32) error
}

type tag struct {
    command command
    requirements []requirement
}

// ===[ requirementRequired ]================================================================[ requirementRequired ]===

type requirementRequired struct {}

func (r requirementRequired) evalBool(field string, result bool) error {
    if result { return nil }
    return fmt.Errorf("%s is required", field)
}

func (r requirementRequired) evalInt(field string, result int) error {
    panic("not an int")
}

func (r requirementRequired) evalFloat(field string, result float32) error {
    panic("not a float")
}

// ===[ requirementComparison ]============================================================[ requirementComparison ]===

type requirementComparison struct {
    constant string
    symbol string
    operationi func(int, int)         bool
    operationf func(float32, float32) bool
}

func (r requirementComparison) evalBool(field string, result bool) error {
    panic("not a bool")
}

func (r requirementComparison) evalInt(field string, result int) error {
    var i, err = strconv.ParseInt(r.constant, 10, 32)
    if err != nil { panic("not an integer constant") }
    
    if r.operationi(result, int(i)) { return nil }
    return fmt.Errorf("%s is %d but must be %s %s", field, result, r.symbol, r.constant)
}

func (r requirementComparison) evalFloat(field string, result float32) error {
    var f, err = strconv.ParseFloat(r.constant, 32)
    if err != nil { panic("not a float constant") }
    
    if r.operationf(result, float32(f)) { return nil }
    return fmt.Errorf("%s is %.2f but must be %s %s", field, result, r.symbol, r.constant)
}

// ===[ commandValue ]==============================================================================[ commandValue ]===

type commandValue struct {
    value string
}

func (c commandValue) evalBool(_ *Binding, _ Extensions) bool {
    switch c.value {
        case "true":  return true
        case "false": return false
        default: panic(fmt.Sprintf("not a boolean: '%s'", c.value))
    }
}

func (c commandValue) evalInt(_ *Binding, _ Extensions) int {
    var result, err = strconv.ParseInt(c.value, 10, 32)
    if err != nil { panic(fmt.Sprintf("not an integer: '%s'", c.value)) }
    return int(result)
}

func (c commandValue) evalFloat(_ *Binding, _ Extensions) float32 {
    var result, err = strconv.ParseFloat(c.value, 32)
    if err != nil { panic(fmt.Sprintf("not a float: '%s'", c.value)) }
    return float32(result)
}

func (c commandValue) hasBoolRepresentation() bool {
    switch c.value {
        case "true":  return true
        case "false": return true
        default:      return false
    }
}

func (c commandValue) hasIntRepresentation() bool {
    var _, err = strconv.ParseInt(c.value, 10, 32)
    return err == nil
}

func (c commandValue) hasFloatRepresentation() bool {
    var hasRadix = (strings.IndexByte(c.value, '.') >= 0)
    var _, err = strconv.ParseFloat(c.value, 32)
    return err == nil && hasRadix
}

// ===[ commandBinaryBoolean ]==============================================================[ commandBinaryBoolean ]===

type commandBinaryBoolean struct {
    a command
    b command
    operation func(bool, bool) bool
}

func (c commandBinaryBoolean) evalBool(b *Binding, e Extensions) bool {
    return c.operation(c.a.evalBool(b, e), c.b.evalBool(b, e))
}

func (c commandBinaryBoolean) evalInt(b *Binding, _ Extensions) int {
    panic("not an integer")
}

func (c commandBinaryBoolean) evalFloat(b *Binding, _ Extensions) float32 {
    panic("not a float")
}

func (c commandBinaryBoolean) hasBoolRepresentation() bool {
    return true
}

func (c commandBinaryBoolean) hasIntRepresentation() bool {
    return false
}

func (c commandBinaryBoolean) hasFloatRepresentation() bool {
    return false
}

// ===[ commandNot ]==================================================================================[ commandNot ]===

type commandNot struct {
    inner command
}

func (c commandNot) evalBool(b *Binding, e Extensions) bool {
    return !c.inner.evalBool(b, e)
}

func (c commandNot) evalInt(b *Binding, _ Extensions) int {
    panic("not an integer")
}

func (c commandNot) evalFloat(b *Binding, _ Extensions) float32 {
    panic("not a float")
}

func (c commandNot) hasBoolRepresentation() bool {
    return true
}

func (c commandNot) hasIntRepresentation() bool {
    return false
}

func (c commandNot) hasFloatRepresentation() bool {
    return false
}

// ===[ commandCompare ]==========================================================================[ commandCompare ]===

type commandCompare struct {
    a command
    b command
    operationi func(int, int)         bool
    operationf func(float32, float32) bool
}

func (c commandCompare) evalBool(b *Binding, e Extensions) bool {
    if c.a.hasFloatRepresentation() && c.b.hasFloatRepresentation() {
        return c.operationf(c.a.evalFloat(b, e), c.b.evalFloat(b, e))
    } else if c.a.hasIntRepresentation() && c.b.hasIntRepresentation() {
        return c.operationi(c.a.evalInt(b, e), c.b.evalInt(b, e))
    } else {
        panic(fmt.Sprintf("cannot compare mismatched types (%+v and %+v)", c.a, c.b))
    }
}

func (c commandCompare) evalInt(b *Binding, e Extensions) int {
    panic("not an integer")
}

func (c commandCompare) evalFloat(b *Binding, e Extensions) float32 {
    panic("not a float")
}

func (c commandCompare) hasBoolRepresentation() bool {
    return false
}

func (c commandCompare) hasIntRepresentation() bool {
    return c.a.hasIntRepresentation() && c.b.hasIntRepresentation()
}

func (c commandCompare) hasFloatRepresentation() bool {
    return c.a.hasFloatRepresentation() && c.b.hasFloatRepresentation()
}

// ===[ commandExt ]==================================================================================[ commandExt ]===

type commandExt struct {
    name string
}

func (c commandExt) evalBool(b *Binding, e Extensions) bool {
    return e.Contains(c.name)
}

func (c commandExt) evalInt(b *Binding, e Extensions) int {
    panic("not an integer")
}

func (c commandExt) evalFloat(b *Binding, e Extensions) float32 {
    panic("not a float")
}

func (c commandExt) hasBoolRepresentation() bool {
    return true
}

func (c commandExt) hasIntRepresentation() bool {
    return false
}

func (c commandExt) hasFloatRepresentation() bool {
    return false
}

// ===[ commandGetIntegerv ]==================================================================[ commandGetIntegerV ]===

type commandGetIntegerv struct {
    name string
}

func (c commandGetIntegerv) evalBool(b *Binding, e Extensions) bool {
    panic("not a bool")
}

func (c commandGetIntegerv) evalInt(b *Binding, e Extensions) int {
    var result int32
    var id, exists = glconstants[c.name]
    if !exists { return 0 }
    b.GetIntegerv(id, &result)
    return int(result)
}

func (c commandGetIntegerv) evalFloat(b *Binding, e Extensions) float32 {
    panic("not a float")
}

func (c commandGetIntegerv) hasBoolRepresentation() bool {
    return false
}

func (c commandGetIntegerv) hasIntRepresentation() bool {
    return true
}

func (c commandGetIntegerv) hasFloatRepresentation() bool {
    return false
}

// ===[ commandGetFloatv ]======================================================================[ commandGetFloatV ]===

type commandGetFloatv struct {
    name string
}

func (c commandGetFloatv) evalBool(b *Binding, e Extensions) bool {
    panic("not a bool")
}

func (c commandGetFloatv) evalInt(b *Binding, e Extensions) int {
    panic("not an int")
}

func (c commandGetFloatv) evalFloat(b *Binding, e Extensions) float32 {
    var result float32
    var id, exists = glconstants[c.name]
    if !exists { return 0.0 }
    b.GetFloatv(id, &result)
    return result
}

func (c commandGetFloatv) hasBoolRepresentation() bool {
    return false
}

func (c commandGetFloatv) hasIntRepresentation() bool {
    return false
}

func (c commandGetFloatv) hasFloatRepresentation() bool {
    return true
}

// ===[ commandIf ]====================================================================================[ commandIf ]===

type commandIf struct {
    clause command
    implication command
    otherwise command
}

func (c commandIf) evalBool(b *Binding, e Extensions) bool {
    if !c.hasBoolRepresentation() { panic(fmt.Sprintf("both clauses of %+v must have a bool representation", c)) }
    if c.clause.evalBool(b, e) {
        return c.implication.evalBool(b, e)
    } else {
        return c.otherwise.evalBool(b, e)
    }
}

func (c commandIf) evalInt(b *Binding, e Extensions) int {
    if !c.hasIntRepresentation() {
        panic(fmt.Sprintf("both clauses of %+v must have an int representation", c))
    }
    if c.clause.evalBool(b, e) {
        return c.implication.evalInt(b, e)
    } else {
        return c.otherwise.evalInt(b, e)
    }
}

func (c commandIf) evalFloat(b *Binding, e Extensions) float32 {
    if !c.hasFloatRepresentation() {
        var suffix = " but neither clause does"
        if c.implication.hasFloatRepresentation() && !c.otherwise.hasFloatRepresentation() {
            suffix = " but the second clause doesn't"
        } else if !c.implication.hasFloatRepresentation() && c.otherwise.hasFloatRepresentation() {
            suffix = " but the first clause doesn't"
        }
        panic(fmt.Sprintf("both clauses of %+v must have a float representation%s", c, suffix))
    }
    if c.clause.evalBool(b, e) {
        return c.implication.evalFloat(b, e)
    } else {
        return c.otherwise.evalFloat(b, e)
    }
}

func (c commandIf) hasBoolRepresentation() bool {
    return c.implication.hasBoolRepresentation() && c.otherwise.hasBoolRepresentation()
}

func (c commandIf) hasIntRepresentation() bool {
    return c.implication.hasIntRepresentation() && c.otherwise.hasIntRepresentation()
}

func (c commandIf) hasFloatRepresentation() bool {
    return c.implication.hasFloatRepresentation() && c.otherwise.hasFloatRepresentation()
}
