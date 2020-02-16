package glcaps

import (
    "fmt"
    "reflect"
    "strings"
)

// parseAtom parses the next (space-delimited) word in the string starting at character offset, returning a slice to
// the first word and an offset to one-past-the-end of the word, or an error.
func parseAtom(tag string, offset int) (atom string, next int) {
    
    // the first index that isn't a space
    var i = offset
    for (i < len(tag)) && (tag[i] == ' ') { i++ }
    
    // the last index that isn't a space
    var j = i
    for (j < len(tag)) && (tag[j] != ' ') { j++ }
    
    if j - i == 0 {
        return "", -1
    }
    
    return tag[i:j], j
}

// parseCommand2 performs the common task of parsing two commands at once and combining their errors
func parseCommand2(tag string, offset int) (c1 command, c2 command, next int, _err error) {
    var ac, ao, ae = parseCommand(tag, offset)
    if ae != nil { return c1, c2, 0, ae }
    
    var bc, bo, be = parseCommand(tag, ao)
    if be != nil { return c1, c2, 0, be }
    
    return ac, bc, bo, nil
}

// parseBinaryBooleanCopmmand performs the common task of parsing a command that is a function with two boolean
// arguments and returns a boolean (e.g. AND, OR)
func parseBinaryBooleanCommand(tag string, offset int, f func(bool, bool) bool) (c command, next int, _err error) {
    var c1, c2, o, e = parseCommand2(tag, offset)
    if e != nil { return c, 0, e }
    return commandBinaryBoolean{c1, c2, f}, o, nil
}

// parseCCompareCommand performs the common task of parsing a command that is a function with two arguments
// and returns a boolean (e.g. LessThan)
func parseCompareCommand(
    tag string,
    offset int,
    fi func(int, int) bool,
    ff func(float32, float32) bool,
    fs func(string, string) bool,
) (c command, next int, _err error) {
    var c1, c2, o, e = parseCommand2(tag, offset)
    if e != nil { return c, 0, e }
    return commandCompare{c1, c2, fi, ff, fs}, o, nil
}

// parseCommand parses and/or/not/ext/GetIntegerv/GetFloatv/if/eq/neq/lt/lte/gt/gte/value commands and returns an
// offset to the end of the parsed command.
func parseCommand(tag string, _offset int) (c command, next int, err error) {
    var start, offset = parseAtom(tag, _offset)
    if offset < 0 { return c, 0, fmt.Errorf("expected command") }
    
    switch start {
        case "and": return parseBinaryBooleanCommand(tag, offset, operationBooleanAnd)
        case "or":  return parseBinaryBooleanCommand(tag, offset, operationBooleanOr)
        
        case "eq":  return parseCompareCommand(tag, offset, operationIntEq,  operationFloat32Eq,  operationStringEq)
        case "neq": return parseCompareCommand(tag, offset, operationIntNeq, operationFloat32Neq, operationStringNeq)
        case "lt":  return parseCompareCommand(tag, offset, operationIntLt,  operationFloat32Lt,  nil)
        case "lte": return parseCompareCommand(tag, offset, operationIntLte, operationFloat32Lte, nil)
        case "gt":  return parseCompareCommand(tag, offset, operationIntGt,  operationFloat32Gt,  nil)
        case "gte": return parseCompareCommand(tag, offset, operationIntGte, operationFloat32Gte, nil)

        case "if":
            var ac, ao, ae = parseCommand(tag, offset)
            if ae != nil { return c, 0, ae }
            
            var bc, bo, be = parseCommand(tag, ao)
            if be != nil { return c, 0, be }
            
            var cc, co, ce = parseCommand(tag, bo)
            if ce != nil { return c, 0, ce }
            
            return commandIf{ac, bc, cc}, co, nil
        
        case "not":
            var c1, o, e = parseCommand(tag, offset)
            if e != nil { return c, 0, e }
            return commandNot{c1}, o, nil
            
        case "ext":
            var c1, o = parseAtom(tag, offset)
            if o < 0 { return c, 0, fmt.Errorf("expected name after ext") }
            return commandExt{c1}, o, nil
        
        case "GetString":
            var c1, o = parseAtom(tag, offset)
            if o < 0 { return c, 0, fmt.Errorf("expected name after GetString") }
            return commandGetString{c1}, o, nil
            
        case "GetIntegerv":
            var c1, o = parseAtom(tag, offset)
            if o < 0 { return c, 0, fmt.Errorf("expected name after GetIntegerv") }
            return commandGetIntegerv{c1}, o, nil
        
        case "GetFloatv":
            var c1, o = parseAtom(tag, offset)
            if o < 0 { return c, 0, fmt.Errorf("expected name after GetFloatv") }
            return commandGetFloatv{c1}, o, nil
            
        default:
            return commandValue{start}, offset, nil
    }
}

// parseParts splits a string on the first semicolon: for "a; b" it returns "a", " b".
func parseParts(s string) (string, string) {
    var offset = strings.IndexByte(s, ';')
    if offset < 0 {
        return s[0: len(s)], s[0:0]
    } else {
        return s[0: offset], s[offset + 1:]
    }
}

func parseCompareRequirement(
    tag string,
    offset int,
    symbol string,
    fi func(int, int) bool,
    ff func(float32, float32) bool,
    fs func(string, string) bool,
) (r requirement, next int, _err error) {
    var r1, o = parseAtom(tag, offset)
    if o < 0 { return r, 0, fmt.Errorf("expected constant after comparison") }
    return requirementComparison{
        constant:   r1,
        symbol:     symbol,
        operationi: fi,
        operationf: ff,
        operations: fs,
    }, o, nil
}

// parseCommand parses and/or/not/ext/GetIntegerv/GetFloatv/if/eq/neq/lt/lte/gt/gte/value commands and returns an
// offset to the end of the parsed command.
func parseRequirement(tag string, _offset int) (r requirement, next int, err error) {
    var start, offset = parseAtom(tag, _offset)
    if offset < 0 { return r, -1, nil }
    
    switch start {
        case "required":
            return requirementRequired{}, offset, nil
        
        case "eq":  return parseCompareRequirement(tag, offset, "=",  operationIntEq,  operationFloat32Eq,  operationStringEq)
        case "neq": return parseCompareRequirement(tag, offset, "!=", operationIntNeq, operationFloat32Neq, operationStringNeq)
        case "lt":  return parseCompareRequirement(tag, offset, "<",  operationIntLt,  operationFloat32Lt,  nil)
        case "lte": return parseCompareRequirement(tag, offset, "<=", operationIntLte, operationFloat32Lte, nil)
        case "gt":  return parseCompareRequirement(tag, offset, ">",  operationIntGt,  operationFloat32Gt,  nil)
        case "gte": return parseCompareRequirement(tag, offset, ">=", operationIntGte, operationFloat32Gte, nil)
        
        default:
            return r, 0, fmt.Errorf("unknown requirement: '%s'", start)
    }
}

func parseRequirements(t string) (requirements []requirement, err error) {
    var offset int
    var r requirement
    requirements = make([]requirement, 0)
    
    for {
        r, offset, err = parseRequirement(t, offset)
        if err != nil { return requirements, err }
        if offset < 0  { return requirements, nil }
    
        requirements = append(requirements, r)
    }
}

// parseTag parses the command and requirements clauses of a tag
func parseTag(t string) (tag, error) {
    var left, right = parseParts(t)
    
    var command, index, err = parseCommand(left, 0)
    if err != nil { return tag{}, err }
    
    if index < len(left) && strings.TrimSpace(left[index + 1:]) != "" {
        return tag{}, fmt.Errorf("unexpected trailing string after end of command: '%s'", left[index:])
    }
    
    var requirements, rerr = parseRequirements(right)
    if rerr != nil { return tag{}, rerr }

    return tag{
        command: command,
        requirements: requirements,
    }, nil
}

func checkBoolRequirements(field reflect.StructField, result bool, rs []requirement) (errors Errors) {
    for _, r := range rs {
        var err = r.evalBool(field.Name, result)
        if err == nil { continue }
        
        errors.append(Error{
            Field: field.Name,
            Tag:   field.Tag.Get("glcaps"),
            Requirement: r,
            Message: err.Error(),
        })
    }
    
    return errors
}

func checkIntRequirements(field reflect.StructField, result int, rs []requirement)  (errors Errors) {
    for _, r := range rs {
        var err = r.evalInt(field.Name, result)
        if err == nil { continue }
        
        errors.append(Error{
            Field: field.Name,
            Tag:   field.Tag.Get("glcaps"),
            Requirement: r,
            Message: err.Error(),
        })
    }
    
    return errors
}

func checkFloatRequirements(field reflect.StructField, result float32, rs []requirement) (errors Errors) {
    for _, r := range rs {
        var err = r.evalFloat(field.Name, result)
        if err == nil { continue }
        
        errors.append(Error{
            Field: field.Name,
            Tag:   field.Tag.Get("glcaps"),
            Requirement: r,
            Message: err.Error(),
        })
    }
    
    return errors
}

func checkStringRequirements(field reflect.StructField, result string, rs []requirement) (errors Errors) {
    for _, r := range rs {
        var err = r.evalString(field.Name, result)
        if err == nil { continue }
        
        errors.append(Error{
            Field: field.Name,
            Tag:   field.Tag.Get("glcaps"),
            Requirement: r,
            Message: err.Error(),
        })
    }
    
    return errors
}

func parseStructField(binding *Binding, extensions []string, field reflect.StructField, setter reflect.Value, value interface{}) (errors Errors) {

    var parse = func(field reflect.StructField) (_tag tag, ok bool) {
        // fmt.Printf("got float32 %s %s %s\n", field.Name, field.Tag, field.Tag.Get("glcaps"))
        var glcapstag, exists = field.Tag.Lookup("glcaps")
        if !exists { return tag{}, false }
        
        var t, err = parseTag(glcapstag)
        if err != nil {
            errors.append(Error{
                Field: field.Name,
                Tag:   glcapstag,
                Message: fmt.Sprintf("tag parse error: %v", err),
            })
            return t, false
        }
        
        return t, true
    }
    
    var kind = field.Type.Kind()
    
    if kind == reflect.Struct {
        errors.append(parseStruct(binding, extensions, setter)...)
    } else {
        var t, ok = parse(field)
        if ok {
            switch kind {
                case reflect.Bool:
                    var result = t.command.evalBool(binding, extensions)
                    errors.append(checkBoolRequirements(field, result, t.requirements)...)
                    setter.SetBool(result)
                    
                case reflect.Int:
                    var result = t.command.evalInt(binding, extensions)
                    errors.append(checkIntRequirements(field, result, t.requirements)...)
                    setter.SetInt(int64(result))
                    
                case reflect.Float32: fallthrough
                case reflect.Float64:
                    var result = t.command.evalFloat(binding, extensions)
                    errors.append(checkFloatRequirements(field, result, t.requirements)...)
                    setter.SetFloat(float64(result))

                case reflect.String:
                    var result = t.command.evalString(binding, extensions)
                    errors.append(checkStringRequirements(field, result, t.requirements)...)
                    setter.SetString(result)
            }
        }
    }
    
    return errors
}

func parseStruct(binding *Binding, extensions []string, s reflect.Value) (errors Errors) {
    
    if s.Kind() != reflect.Struct {
        panic("target must be a struct or pointer to struct")
    }
    
    for i := 0; i < s.NumField(); i++ {
        errors.append(parseStructField(binding, extensions, s.Type().Field(i), s.Field(i), s.Field(i).Interface())...)
    }
    
    return errors
}

// Parse parses a struct and parses struct tag annotations to identify the required OpenGL information. It fills the
// target struct with the results, and returns zero or more Errors if any defined requirements are not met. It also
// returns a sorted string list of all supported OpenGL extensions.
//
// The struct tag key is `glcaps`. The struct tag syntax is a space-separated list of commands, optionally followed
// by a colon and a space-separated list of requirements.
//
// Commands:
//
//    and command1 command2          - return true if command1 and command2 are true
//    or  command1 command2          - return true if either command1 or command2 are true
//    not command                    - return the boolean opposite of a command
//    ext GL_EXT_name                - return true if the given extension is supported
//    GetIntegerv GL_name            - lookup and return an integer value
//    GetFloatv GL_name              - lookup and return a float value
//    if command1 command2 command3  - if command1 is true, return the result of command2 otherwise return command3
//    eq|neq|lt|lte|gt|gte command1 command2 - return true if command1 ==/!=/</<=/>/>= command2 respectively
//    value                          - a value literal (e.g. true, false, 123, 1.23, 128KiB)
//
// Requirements:
//
//    required                       - generate an error if the result is not true
//    eq|neq|lt|lte|gt|gte value     - generate an error if the command is not ==, !=, <, <=, >, >= value respectively
//
func Parse(binding *Binding, target interface{}) (extensions Extensions, errors Errors) {
    extensions = binding.QueryExtensions()
    return extensions, parseStruct(binding, extensions, reflect.ValueOf(target).Elem())
}

