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

// Decode parses the DirectX (.x) file format, with an optional list of user-defined templates (may be empty or nil),
// and on success returns a File object containing the decoded data.
func Decode(r io.Reader, templates []*Template) (file *File, err error) {
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
    
    var templatesByName = make(map[string]*Template)
    for _, template := range defaultTemplates {
        templatesByName[template.Name] = template
    }
    if templates != nil {
        for _, template := range templates {
            templatesByName[template.Name] = template
        }
    }
    
    file = &File{
        Children: make([]Data, 0),
        ReferencesByName: make(map[string]*Data),
        ReferencesByUUID: make(map[UUID_t]*Data),
        templatesByName: templatesByName,
        Templates: templatesByName, // TODO remove
    }
    
    for {
        var word string
        word, err = readAtom(fp)
        if err != nil { break }
    
        if word == "template" {
            var t *Template
            t, err = decodeTemplate(fp, templatesByName)
            if err != nil { break }
            templatesByName[t.Name] = t
        } else {
            var t, ok = templatesByName[word]
            if !ok { return nil, fmt.Errorf("unknown data name %s", word) }
            
            var data *Data
            data, err = decodeDataBlock(fp, file, t, templatesByName)
            if err != nil { return nil, err }
            file.appendChild(data)
        }
    }
    
    if err == io.EOF { err = nil }
    if err != nil { file = nil }
    
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
func decodeTemplate(r *bufio.Reader, templates map[string]*Template) (t *Template, err error) {
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
    
    We also check that each template member is defined (no forward references!).
    
    There's no technical reason why forward references cannot be implemented, but it appears that no-one else
    supports them either. All we really want to know is that a template is not recursive but we get that for free by
    checking each template member is defined yet.
    */
    t = &Template{}
    t.Name = mustReadAtom(r)
    if mustReadSymbol(r) != '{' { return nil, fmt.Errorf("expected '{' to begin template") }
    if mustReadSymbol(r) != '<' { return nil, fmt.Errorf("expected '<' to start UUID") }
    t.UUID = MustHexToUUID(mustReadAtom(r)) // TODO handle errors
    if mustReadSymbol(r) != '>' { return nil, fmt.Errorf("expected '>' to close UUID") }
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
        
        //symbol = mustReadSymbol(r)
        //if symbol != ';' { return nil, fmt.Errorf("expected ';'") }
        mustReadExactSymbol(r, ';', "end of template field")
        
        var member = TemplateMember{
            Name: fieldName,
            Type: fieldType,
            Dimensions: dimensions,
        }
        
        // is the type okay?
        if !member.isPrimitiveType() {
            _, exists := templates[fieldType]
            if !exists {
                return nil,
                fmt.Errorf("unknown template field type %s (no forward or recursive references please)", fieldType)
            }
        }
        
        t.Members = append(t.Members, member)
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
    fns, err := decodeMembers(r, data, data, t.Members, templates)
    if err != nil { return nil, err }
    if fns != nil {
        for _, f := range fns {
            f(data)
        }
    }
    
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
func decodeMembers(r *bufio.Reader, data *Data, target *Data, members []TemplateMember, templates map[string]*Template) (fns []func(*Data), err error) {
    for i, member := range members {
        fns, err = decodeMemberValue(r, data, target, i, &member, templates)
        if err != nil { return nil, err }
    }
    
    return fns, err
}

// decodeMemberValue decodes a value according to a template member, possibly an array
func decodeMemberValue(r *bufio.Reader, data *Data, target *Data, index int, member *TemplateMember, templates map[string]*Template) (fns []func(*Data), err error) {
    
    if member.Dimensions == nil {
        // read a single value
        
        f, err := decodeSingleValue(r, data, member, -1, templates)
        if err != nil { return nil, err }
        if f != nil { f(data); }
        
    } else if len(member.Dimensions) == 1 {
        // read a 1D array
        
        var ln, err = strconv.ParseInt(member.Dimensions[0], 10, 32)
        if err != nil {
            // array of variable length
            var len32 uint32
            
            var offset, size int
            _, offset, size, err = target.GetNamedField(member.Dimensions[0], "DWORD", templates)
            
            offset2, _, _ := data.GetField(index, templates)
            offset += offset2
            
            // len32, err = target.GetNamedDWORD(member.Dimensions[0], templates)
            
            if err != nil { return nil, fmt.Errorf("unable to lookup variable dimension length for field %s referencing %s: %v", member.Name, member.Dimensions[0], err) }
            
            fmt.Printf("offset %d, size %d\n", offset, size)
            len32 = binary.LittleEndian.Uint32(data.Bytes[offset : offset + size])
            fmt.Printf("array length %d\n", len32)
            
            ln = int64(len32)
        }
        
        var arrayIndex = data.appendArray()
        
        for i := 0; i < int(ln); i++ {
            f, err := decodeSingleValue(r, data, member, arrayIndex, templates)
            if err != nil { return nil, err }
            if f != nil { f(data); }
            
            if i + 1 < int(ln) {
                mustReadExactSymbol(r, ',', "array item separator")
            }
        }
        
    } else {
        // Read a multidimensional array
        return nil, fmt.Errorf("multidimensional arrays not yet supported")
    }
    
    //if mustReadSymbol(r) != ';' { return nil, fmt.Errorf("expected ';'") }
    mustReadExactSymbol(r, ';', fmt.Sprintf("end of object member value while parsing %s.%s", target.SpecName(), member.Name))
    
    return nil, nil
}

// decodeMemberValue decodes a value according to a template member, but its not an array
func decodeSingleValue(r *bufio.Reader, data *Data, member *TemplateMember, arrayIndex int, templates map[string]*Template) (f func(*Data), err error) {
    
    if member.isPrimitiveType() {
        switch member.Type {
        
            case "DWORD":
                var dword int64
                dword, err = strconv.ParseInt(mustReadAtom(r), 10, 32)
                if dword < 0 { dword = -dword }
                //return func(d *Data) { d.appendDWORD(uint32(dword), arrayIndex) }, nil
                data.appendDWORD(uint32(dword), arrayIndex)
                
            case "float":
                var float float64
                float, err = strconv.ParseFloat(mustReadAtom(r), 32)
                //return func(d *Data) { d.appendFloat32(float32(float), arrayIndex) }, nil
                data.appendFloat32(float32(float), arrayIndex)
                
            case "WORD":
                var word int64
                word, err = strconv.ParseInt(mustReadAtom(r), 10, 16)
                if word < 0 { word = -word }
                // return func(d *Data) { d.appendWORD(uint16(word), arrayIndex) }, nil
                data.appendWORD(uint16(word), arrayIndex)
                
            case "STRING":
                if arrayIndex >= 0 { return nil, fmt.Errorf("string arrays not yet supported") }
                mustReadExactSymbol(r, '"', "open quote")
                var s = mustReadString(r)
                mustReadExactSymbol(r, '"', "close quote")
                //return func(d *Data) { d.appendString(s, arrayIndex) }, nil
                data.appendString(s, arrayIndex)
                
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
        //var subt, ok = templates[member.Type]
        //if !ok { return nil, fmt.Errorf("unrecognised named data type %s for %s in %s", member.Type, member.Name, data.Spec.Name) }
        
        // lookup template by field type is guaranteed to work at this point
        var subt = templates[member.Type]
        
        // TODO this should be refactored so that decodeMembers and decodeData return data without modifying data
        var subdata = &Data{Spec: subt} // just used for spec
        fns, err := decodeMembers(r, data, subdata, subt.Members, templates)
        if err != nil { return nil, err }
        // data.appendChild(subdata)
        return func(d *Data) {
            if fns != nil {
                for _, f := range fns {
                    f(d)
                }
            }
        }, nil
    }
    
    return nil, nil
}
