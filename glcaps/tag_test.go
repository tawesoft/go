package glcaps

import (
    "testing"
)

func TestParseAtom1(t *testing.T) {
    var atom, index = parseAtom("a b c", 0)
    if atom != "a" { t.Errorf("unexpected result") }
    if index != 1  { t.Errorf("unexpected result") }
}

func TestParseAtom2(t *testing.T) {
    var atom, index = parseAtom("a b c", 1)
    if atom != "b" { t.Errorf("unexpected result") }
    if index != 3  { t.Errorf("unexpected result") }
}

func TestParseAtom3(t *testing.T) {
    var atom, index = parseAtom("a b c", 3)
    if atom != "c" { t.Errorf("unexpected result") }
    if index != 5  { t.Errorf("unexpected result") }
}

func TestParseParts(t *testing.T) {
    var left, right = parseParts("a; b")
    if left  != "a"  { t.Errorf("unexpected result, got '%s'", left) }
    if right != " b" { t.Errorf("unexpected result, got '%s'", right) }
}

func TestCommandExt(t *testing.T) {
    
    var ext = Extensions{
        "FOO",
    }
    
    var command, _, err = parseCommand("ext FOO", 0)
    if err != nil  { t.Errorf("unexpected result") }
    
    var cExt = command.(commandExt)
    if (cExt.name != "FOO") { t.Errorf("unexpected result: expected 'FOO' but got '%s'", cExt.name) }
    
    var result = command.evalBool(nil, ext)
    if result != true { t.Errorf("unexpected result") }
}

func TestParseCommandCompound(t *testing.T) {
    
    var command, _, err = parseCommand("and true or true false", 0)
    if err != nil { t.Failed() }
    
    var cAnd = command.(commandBinaryBoolean)
    var cArg1a = cAnd.a.(commandValue)
    var cArg1b = cAnd.b.(commandBinaryBoolean)
    var cArg2a = cArg1b.a.(commandValue)
    var cArg2b = cArg1b.b.(commandValue)
    if (cArg1a.value != "true") { t.Errorf("unexpected result") }
    if (cArg2a.value != "true") { t.Errorf("unexpected result") }
    if (cArg2b.value != "false") { t.Errorf("unexpected result") }
}

func TestEvalCommandAnd1(t *testing.T) {
    
    var command, _, err = parseCommand("and true true", 0)
    if err != nil { t.Failed() }
    if !command.evalBool(nil, nil) { t.Errorf("unexpected result") }
}

func TestEvalCommandAnd2(t *testing.T) {
    
    var command, _, err = parseCommand("and true false", 0)
    if err != nil { t.Failed() }
    if command.evalBool(nil, nil) { t.Errorf("unexpected result") }
}

func TestEvalCommandOr1(t *testing.T) {
    
    var command, _, err = parseCommand("or true true", 0)
    if err != nil { t.Failed() }
    if !command.evalBool(nil, nil) { t.Errorf("unexpected result") }
}

func TestEvalCommandCompoundBoolean1(t *testing.T) {
    
    var command, _, err = parseCommand("and or true false or false true", 0)
    if err != nil { t.Failed() }
    if !command.evalBool(nil, nil) { t.Errorf("unexpected result") }
}

func TestEvalCommandCompoundBoolean2(t *testing.T) {
    
    var command, _, err = parseCommand("or and true true and false false", 0)
    if err != nil { t.Failed() }
    if !command.evalBool(nil, nil) { t.Errorf("unexpected result") }
}

func TestEvalCommandExtraSpace(t *testing.T) {
    var command, _, err = parseCommand("and  true  true", 0)
    if err != nil { t.Errorf("unexpected result") }
    if !command.evalBool(nil, nil) { t.Errorf("unexpected result") }
}

func TestEvalCommandLeadingSpace(t *testing.T) {
    var command, _, err = parseCommand("  and true true", 0)
    if err != nil { t.Errorf("unexpected result") }
    if !command.evalBool(nil, nil) { t.Errorf("unexpected result") }
}

func TestParseCommandTrailingSpace(t *testing.T) {
    var _, _, err = parseCommand("and true true ", 0)
    if err != nil { t.Errorf("unexpected result") }
}

func TestParseCommandUnfinishedCompound1(t *testing.T) {
    var _, _, err = parseCommand("and true", 0)
    if err == nil { t.Errorf("unexpected result") }
}

func TestParseCommandUnfinishedCompound2(t *testing.T) {
    var _, _, err = parseCommand("and true ", 0)
    if err == nil { t.Errorf("unexpected result") }
}

func TestParseCommandUnfinishedExt(t *testing.T) {
    var _, _, err = parseCommand("ext", 0)
    if err == nil { t.Errorf("unexpected result") }
}

func TestParseTagTrailing(t *testing.T) {
    var _, err = parseTag("and true true trailing")
    if err == nil { t.Errorf("unexpected result - expected an error") }
}

func TestParseTagTrailingSpace(t *testing.T) {
    var _, err = parseTag("and true true ")
    if err != nil { t.Errorf("unexpected result") }
}

func TestParseTagExtraSpace(t *testing.T) {
    var _, err = parseTag("   and  true  true")
    if err != nil { t.Errorf("unexpected result, got %v", err) }
}

func TestEvalCommandLt1(t *testing.T) {
    var command, _, err = parseCommand("lt 1 2", 0)
    if err != nil { t.Failed() }
    if !command.evalBool(nil, nil) { t.Errorf("unexpected result") }
}

func TestEvalCommandLt2(t *testing.T) {
    var command, _, err = parseCommand("lt 2.5 1.0", 0)
    if err != nil { t.Failed() }
    if command.evalBool(nil, nil) { t.Errorf("unexpected result") }
}


/*
func TestParseTagCommand4(t *testing.T) {

//    var command, _, err = parseCommand("gt 123 100")
//    if err != nil { t.Failed() }
    
    //var cGt = command.(commandGT)
    //var cValue1 = command.(commandValue)
    //var cValue2 = command.(commandValue)
}
*/
