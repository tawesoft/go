package glcaps

import (
    "fmt"
    "strconv"
    "strings"
)

func operationStringEq  (a string,   b string) bool { return a == b }
func operationStringNeq (a string,   b string) bool { return a != b }

type command interface{
    evalBool  (b *Binding, extensions Extensions) bool
    evalInt   (b *Binding, extensions Extensions) int
    evalFloat (b *Binding, extensions Extensions) float32
    evalString(b *Binding, extensions Extensions) string
    hasBoolRepresentation() bool
    hasIntRepresentation() bool
    hasFloatRepresentation() bool
    hasStringRepresentation() bool
}

type requirement interface{
    evalBool  (field string, result bool)    error
    evalInt   (field string, result int)     error
    evalFloat (field string, result float32) error
    evalString(field string, result string) error
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

func (r requirementRequired) evalString(field string, result string) error {
    if len(result) > 0 { return nil }
    return fmt.Errorf("%s is required", field)
}

// ===[ requirementComparison ]============================================================[ requirementComparison ]===

type requirementComparison struct {
    constant string
    symbol string
    operationi func(int, int)         bool
    operationf func(float32, float32) bool
    operations func(string, string)   bool
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

func (r requirementComparison) evalString(field string, result string) error {
    if r.operations(result, r.constant) { return nil }
    return fmt.Errorf("%s is %s but must be %s %s", field, result, r.symbol, r.constant)
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

func (c commandValue) evalString(_ *Binding, _ Extensions) string {
    return c.value
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

func (c commandValue) hasStringRepresentation() bool {
    return true
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

func (c commandBinaryBoolean) evalString(_ *Binding, _ Extensions) string {
    panic("not a string")
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

func (c commandBinaryBoolean) hasStringRepresentation() bool {
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

func (c commandNot) evalString(b *Binding, _ Extensions) string {
    panic("not a string")
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

func (c commandNot) hasStringRepresentation() bool {
    return false
}

// ===[ commandCompare ]==========================================================================[ commandCompare ]===

type commandCompare struct {
    a command
    b command
    operationi func(int, int)         bool
    operationf func(float32, float32) bool
    operations func(string,  string)  bool
}

func (c commandCompare) evalBool(b *Binding, e Extensions) bool {
    if c.a.hasFloatRepresentation() && c.b.hasFloatRepresentation() {
        return c.operationf(c.a.evalFloat(b, e), c.b.evalFloat(b, e))
    } else if c.a.hasIntRepresentation() && c.b.hasIntRepresentation() {
        return c.operationi(c.a.evalInt(b, e), c.b.evalInt(b, e))
    } else if c.a.hasStringRepresentation() && c.b.hasStringRepresentation() {
        if c.operations != nil {
            return c.operations(c.a.evalString(b, e), c.b.evalString(b, e))
        } else {
            panic(fmt.Sprintf("string operation not defined for this comparison"))
        }
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

func (c commandCompare) evalString(b *Binding, e Extensions) string {
    panic("not a string")
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

func (c commandCompare) hasStringRepresentation() bool {
    return c.a.hasStringRepresentation() && c.b.hasStringRepresentation()
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

func (c commandExt) evalString(b *Binding, e Extensions) string {
    panic("not a string")
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

func (c commandExt) hasStringRepresentation() bool {
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

func (c commandGetIntegerv) evalString(b *Binding, e Extensions) string {
    panic("not a string")
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

func (c commandGetIntegerv) hasStringRepresentation() bool {
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

func (c commandGetFloatv) evalString(b *Binding, e Extensions) string {
    panic("not a string")
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

func (c commandGetFloatv) hasStringRepresentation() bool {
    return false
}

// ===[ commandGetString ]======================================================================[ commandGetString ]===

type commandGetString struct {
    name string
}

func (c commandGetString) evalBool(b *Binding, e Extensions) bool {
    panic("not a bool")
}

func (c commandGetString) evalInt(b *Binding, e Extensions) int {
    panic("not an int")
}

func (c commandGetString) evalFloat(b *Binding, e Extensions) float32 {
    panic("not a float")
}

func (c commandGetString) evalString(b *Binding, e Extensions) string {
    var id, exists = glconstants[c.name]
    if !exists { return "" }
    return b.GetString(id)
}

func (c commandGetString) hasBoolRepresentation() bool {
    return false
}

func (c commandGetString) hasIntRepresentation() bool {
    return false
}

func (c commandGetString) hasFloatRepresentation() bool {
    return false
}

func (c commandGetString) hasStringRepresentation() bool {
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
    if !c.hasIntRepresentation() { panic(fmt.Sprintf("both clauses of %+v must have an int representation", c)) }
    if c.clause.evalBool(b, e) {
        return c.implication.evalInt(b, e)
    } else {
        return c.otherwise.evalInt(b, e)
    }
}

func (c commandIf) evalFloat(b *Binding, e Extensions) float32 {
    if !c.hasFloatRepresentation() { panic(fmt.Sprintf("both clauses of %+v must have an int representation", c)) }
    if c.clause.evalBool(b, e) {
        return c.implication.evalFloat(b, e)
    } else {
        return c.otherwise.evalFloat(b, e)
    }
}

func (c commandIf) evalString(b *Binding, e Extensions) string {
    if !c.hasStringRepresentation() { panic(fmt.Sprintf("both clauses of %+v must have a string representation", c)) }
    if c.clause.evalBool(b, e) {
        return c.implication.evalString(b, e)
    } else {
        return c.otherwise.evalString(b, e)
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

func (c commandIf) hasStringRepresentation() bool {
    return c.implication.hasStringRepresentation() && c.otherwise.hasStringRepresentation()
}
