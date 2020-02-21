package xff

// xff/decode.go decodes the high level structure of a DirectX .x file

import (
    "bufio"
    "encoding/binary"
    "fmt"
    "io"
    "runtime/debug"
    "strconv"
)

func Decode(r io.Reader) (file *File, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("%v\n%s", r.(error), debug.Stack())
        }
    }()
    
    var fp = bufio.NewReader(r)
    
    var format byte // 't'/text or 'b'/binary
    var floatSize int // 32 or 64
    
    format, floatSize, err = decodeHeader(fp)
    if err != nil { return nil, fmt.Errorf("error decoding DirectX (.x) file header: %v", err) }
    
    if format != 't' { return nil, fmt.Errorf("format not implemented") }
    if floatSize != 32 { return nil, fmt.Errorf("float size not implemented") }
    
    var templates = make(map[string]*Template)
    for k, v := range Templates {
        templates[k] = v
    }
    
    file = &File{
        Children: make([]Data, 0),
        ReferencesByName: make(map[string]*Data),
        ReferencesByUUID: make(map[UUID_t]*Data),
    }
    
    for {
        var word string
        word, err = readAtom(fp)
        if err != nil { break }
    
        if word == "template" {
            var t *Template
            t, err = decodeTemplate(fp)
            if err != nil { break }
            templates[t.Name] = t
        } else {
            var t, ok = templates[word]
            if !ok { return nil, fmt.Errorf("unknown data name %s", word) }
            
            var data *Data
            data, err = decodeDataBlock(fp, file, t, templates)
            if err != nil { return nil, err }
            file.appendChild(data)
        }
    }
    
    if err == io.EOF { err = nil }
    if err != nil { file = nil }
    
    file.Templates = templates
    return file, err
}

// DecodeHeader reads the header of a DirectX .x file and returns the format and floatSize on success.
//
// * format will be either 't' (text) or 'b' (binary).
//
// * floatSize will be either 32 or 64.
//
func decodeHeader(r io.Reader) (format byte, floatSize int, err error) {
    // https://docs.microsoft.com/en-us/windows/win32/direct3d9/reserved-words--header--and-comments
    var record struct {
        Magic     [4]byte
        Version   [4]byte
        Format    [4]byte
        FloatSize [4]byte
    }
    
    err = binary.Read(r, binary.LittleEndian, &record)
    if err != nil { goto fail }
    
    if string(record.Magic[:]) != "xof " { err = fmt.Errorf("invalid magic bytes (not a DirectX .x file)"); goto fail }
    if string(record.Version[:]) != "0303" { err = fmt.Errorf("unsupported version"); goto fail }
    
    switch string(record.Format[:]) {
        case "txt ": format = 't'
        case "bin ": format = 'b'
        default: err = fmt.Errorf("unsupported format"); goto fail
    }

    switch string(record.FloatSize[:]) {
        case "0064": floatSize = 64
        case "0032": floatSize = 32
        default: err = fmt.Errorf("unsupported float size"); goto fail
    }
    
    // fallthrough
    fail:
        return format, floatSize, err
}

// decodeTemplate reads a template section of a DirectX .x file. Note that the reader at this step has already
// consumed the leading "template" word.
func decodeTemplate(r *bufio.Reader) (t *Template, err error) {
    /*
    https://docs.microsoft.com/en-us/windows/win32/direct3d9/dx9-graphics-reference-x-file-textencoding-templates
    
    Template:
    
    template <template-name> {
    <UUID>
        <member 1>;
    ...
        <member n>;
    [restrictions]
    }
    
    Member:
    
    array <data-type> <name>[<dimension-size>];
    
    Restrictions:
    Open: [ ... ]
    Restricted: [ { data-type [ UUID ] , } ... ]
    */
    t = &Template{}
    t.Name = mustReadAtom(r)
    if mustReadSymbol(r) != '{' { return nil, fmt.Errorf("expected {") }
    if mustReadSymbol(r) != '<' { return nil, fmt.Errorf("expected <") }
    t.UUID = MustHexToUUID(mustReadAtom(r)) // TODO handle errors
    if mustReadSymbol(r) != '>' { return nil, fmt.Errorf("expected >") }
    t.Mode = 'c' // 'c'/closed by default
    
    t.Members = make([]TemplateMember, 0)
    
    for {
        var fieldName, fieldType string
        var dimensions []string
        
        // Try to end the block
        var symbol = mustReadSymbol(r)
        if symbol == '}' { break }
        if r.UnreadByte() != nil { return nil, fmt.Errorf("I/O rewind error") }
        
        // otherwise its a data type
        var word = mustReadAtom(r)
        if word == "array" {
            fieldType = mustReadAtom(r)
            fieldName = mustReadAtom(r)
            dimensions = make([]string, 0)
            
            for {
                if mustReadSymbol(r) != '[' { return nil, fmt.Errorf("expected [") }
                dimensions = append(dimensions, mustReadAtom(r))
                if mustReadSymbol(r) != ']' { return nil, fmt.Errorf("expected ]") }
                
                // try to end the line
                symbol = mustReadSymbol(r)
                if r.UnreadByte() != nil { return nil, fmt.Errorf("I/O rewind error") }
                if symbol == ';' { break }
            }
        } else {
            fieldType = word
            fieldName = mustReadAtom(r)
            dimensions = nil
        }
        
        symbol = mustReadSymbol(r)
        if symbol != ';' { return nil, fmt.Errorf("expected ';'") }
        
        t.Members = append(t.Members, TemplateMember{
            Name: fieldName,
            Type: fieldType,
            Dimensions: dimensions,
        })
    }
    
    return t, nil
}

// decodeDataBlock reads a data section of a DirectX .x file. Note that the reader at this step has already
// consumed the leading identifier and successfully matched it to a template
func decodeDataBlock(r *bufio.Reader, f *File, t *Template, templates map[string]*Template) (data *Data, err error) {
/*
        <Identifier> [name] { [<UUID>]
    <member 1>;
...
    <member n>;
}
 */
    
    // Is there a name?
    var name string
    var symbol = mustReadSymbol(r)
    if symbol != '{' {
        if r.UnreadByte() != nil { return nil, fmt.Errorf("I/O rewind error") }
        name = mustReadAtom(r)
        if mustReadSymbol(r) != '{' {
            return nil, fmt.Errorf("expected {")
        }
    }
    
    data = &Data{Name: name, Spec: t}
    
    if len(name) > 0 {
        f.ReferencesByName[name] = data
    }
    
    // Read members first
    err = decodeMembers(r, data, t.Members, templates)
    if err != nil { return nil, err }
    
    // Read additional data blocks up to closing '}'
    for {
        // More data?
        if mustReadSymbol(r) == '}' { break }
        if r.UnreadByte() != nil { return nil, fmt.Errorf("I/O rewind error") } // TODO peek instead
        if t.Mode == 'c' { return nil, fmt.Errorf("unexpected extra data in closed data type") }
        
        if mustReadSymbol(r) == '{' {
            // data reference
            var reference = mustReadAtom(r)

            data.appendChild(&Data{Name: reference, Spec: nil})
            
            if mustReadSymbol(r) != '}' { return nil, fmt.Errorf("expected } to end reference") }
        } else {
            if r.UnreadByte() != nil { return nil, fmt.Errorf("I/O rewind error") } // TODO peek instead
        }
        
        var word string
        word = mustReadAtom(r)
        
        var t, ok = templates[word]
        if !ok { return nil, fmt.Errorf("unknown identifier '%s'", word) }
        
        var subdata *Data
        subdata, err = decodeDataBlock(r, f, t, templates)
        if err != nil { return nil, err }
        data.appendChild(subdata)
    }
    
    return data, err
}

// decodeMembers decodes values according to a template
func decodeMembers(r *bufio.Reader, data *Data, members []TemplateMember, templates map[string]*Template) (err error) {
    for _, member := range members {
        err = decodeMemberValue(r, data, &member, templates)
        if err != nil { return err }
    }
    
    return err
}

// decodeMemberValue decodes a value according to a template member, possibly an array
func decodeMemberValue(r *bufio.Reader, data *Data, member *TemplateMember, templates map[string]*Template) (err error) {
    
    if member.Dimensions == nil {
        // read a single value
        
        err = decodeSingleValue(r, data, member, -1, templates)
        if err != nil { return err }
        
    } else if len(member.Dimensions) == 1 {
        // read a 1D array
        
        var ln, err = strconv.ParseInt(member.Dimensions[0], 10, 32)
        if err != nil {
            // array of variable length
            var len32 uint32
            len32, err = data.GetNamedDWORD(member.Dimensions[0], templates)
            if err != nil { return fmt.Errorf("unable to lookup variable dimension length for field %s referencing %s: %v", member.Name, member.Dimensions[0], err) }
            ln = int64(len32)
        }
        
        var index = data.appendArray()
        
        for i := 0; i < int(ln); i++ {
            err = decodeSingleValue(r, data, member, index, templates)
            if err != nil { return err }
            
            if i + 1 < int(ln) {
                if mustReadSymbol(r) != ',' { return fmt.Errorf("expected , while reading array") }
            }
        }
        
    } else {
        // Read a multidimensional array
        return fmt.Errorf("multidimensional arrays not yet supported")
    }
    
    if mustReadSymbol(r) != ';' { return fmt.Errorf("expected ';'") }
    
    return nil
}

// decodeMemberValue decodes a value according to a template member, but its not an array
func decodeSingleValue(r *bufio.Reader, data *Data, member *TemplateMember, arrayIndex int, templates map[string]*Template) (err error) {
    
    if member.isPrimitiveType() {
        switch member.Type {
        
            case "DWORD":
                var dword int64
                dword, err = strconv.ParseInt(mustReadAtom(r), 10, 32)
                if dword < 0 { dword = -dword }
                data.appendDWORD(uint32(dword), arrayIndex)
                
            case "float":
                var float float64
                float, err = strconv.ParseFloat(mustReadAtom(r), 32)
                data.appendFloat32(float32(float), arrayIndex)
                
            case "WORD":
                var dword int64
                dword, err = strconv.ParseInt(mustReadAtom(r), 10, 16)
                if dword < 0 { dword = -dword }
                data.appendWORD(uint16(dword), arrayIndex)
                
            case "STRING":
                if arrayIndex >= 0 { return fmt.Errorf("string arrays not yet supported") }
                
                // TODO special mustReadString that takes strings with escaped quotes
                if mustReadSymbol(r) != '"' { return fmt.Errorf("expected open quote") }
                data.appendString(mustReadAtom(r), arrayIndex)
                if mustReadSymbol(r) != '"' { return fmt.Errorf("expected close quote") }
                
            case "FLOAT":  fallthrough
            case "DOUBLE": fallthrough
            case "CHAR":   fallthrough
            case "UCHAR":  fallthrough
            case "BYTE":   fallthrough
            default:
                panic(fmt.Sprintf("primitive type %s not handled (should never happen)", member.Type))
        }
    } else {
        // read a named data type (using a template)
        var subt, ok = templates[member.Type]
        if !ok { return fmt.Errorf("unrecognised named data type %s for %s in %s", member.Type, member.Name, data.Spec.Name) }
        
        // TODO should this attach to data instead of subdata?
        var subdata = &Data{Spec: subt}
        err = decodeMembers(r, subdata, subt.Members, templates)
        if err != nil { return err }
        data.appendChild(subdata)
    }
    
    return nil
}
