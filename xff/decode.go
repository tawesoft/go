package xff

// xff/decode.go decodes the high level structure of a DirectX .x file

import (
    "bufio"
    "encoding/binary"
    "fmt"
    "io"
    "strconv"
)

// Decode parses the DirectX (.x) file format, with an optional list of user-defined templates (may be empty or nil),
// and on success returns a File object containing the decoded data. A DirectX (.x) file may define its own templates.
func Decode(r io.Reader, templates []*Template) (file *File, err error) {
    // Parse errors just panic and are recovered here; finer grained control is not part of the public interface
    // so just catching the panic is nice and simple.
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("DirectX (.x) file format decode error: %v\n", r.(error)) //debug.Stack())
        }
    }()
    
    var fp = bufio.NewReader(r)
    
    var format, _ = decodeHeader(fp)
    if format != 't' { panic(fmt.Errorf("non-text formats not implemented (yet)")) }
    
    // templates are replaced such that in-file templates are prioritised first, caller-supplied templates second,
    // and default templates last.
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
        Children: make([]Data, 0), // never nil
        ReferencesByName: make(map[string]*Data),
        ReferencesByUUID: make(map[UUID_t]*Data),
        templatesByName: templatesByName,
    }
    
    for {
        var word string
        word, err = readAtom(fp) // EOF is allowed here
        if err != nil { break }
    
        if word == "template" {
            var t = decodeTemplate(fp, templatesByName)
            templatesByName[t.Name] = t
        } else {
            var t, ok = templatesByName[word]
            if !ok { panic(fmt.Errorf("unknown object type '%s'", word)) }
            file.appendChild(decodeObject(fp, file, t))
        }
    }
    
    return file, nil
}

// DecodeHeader reads the header of a DirectX .x file and returns the format and floatSize on success:
// `format` will be either 't' (text) or 'b' (binary) and `floatSize` will be either 32 or 64.
func decodeHeader(r io.Reader) (format byte, floatSize int) {
    var record struct {
        Magic     [4]byte
        Version   [4]byte
        Format    [4]byte
        FloatSize [4]byte
    }
    
    var err = binary.Read(r, binary.LittleEndian, &record)
    if (err != nil) || string(record.Magic[:]) != "xof " {
        panic(fmt.Errorf("invalid magic bytes (not a DirectX .x file)"))
    }
    
    switch string(record.Version[:]) {
        case "0303":
            // fine
        case "0302":
            // probably fine!
            // what's the difference?!
            // who knows! its not an open spec!
            // at a guess, v3.2 probably has fewer default templates defined by default but that doesn't matter
        default:
            panic(fmt.Errorf("unsupported file version %v", record.Version))
    }
    
    switch string(record.Format[:]) {
        case "txt ": format = 't'
        case "bin ": format = 'b'
        default:
            panic(fmt.Errorf("unsupported file type %v", record.Format))
    }

    switch string(record.FloatSize[:]) {
        case "0064": floatSize = 64
        case "0032": floatSize = 32
        default:
            // don't actually care what this value is in text mode but lets be strict and validate it anyway
            panic(fmt.Errorf("unsupported file float size %v", record.FloatSize))
    }
    
    return format, floatSize
}

// decodeTemplate reads a template section of a DirectX .x file. Note that the reader at this step has already
// consumed the leading "template" word.
func decodeTemplate(r *bufio.Reader, templates map[string]*Template) *Template {
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
    checking if each template member is defined yet.
    */
    var t = &Template{}
    t.Name = mustReadAtom(r)
    mustReadExpectedSymbol(r, '{', "start of template")
    mustReadExpectedSymbol(r, '<', "start of UUID")
    t.UUID = MustHexToUUID(mustReadAtom(r))
    mustReadExpectedSymbol(r, '>', "close of UUID")
    t.Mode = 'c' // 'c'/closed by default
    
    t.Members = make([]TemplateMember, 0)
    
    for {
        var fieldName, fieldType string
        var dimensions []string
        
        // Try to end the block
        var symbol = mustReadSymbol(r)
        if symbol == '}' { break }
        mustUnreadByte(r)
        
        // otherwise its a data type
        var word = mustReadAtom(r)
        if word == "array" {
            fieldType = mustReadAtom(r)
            fieldName = mustReadAtom(r)
            dimensions = make([]string, 0)
            
            for {
                mustReadExpectedSymbol(r, '[', "array dimension start")
                dimensions = append(dimensions, mustReadAtom(r))
                mustReadExpectedSymbol(r, ']', "array dimension end")
                
                // try to end the line
                symbol = mustPeekSymbol(r)
                if symbol == ';' { break }
            }
        } else {
            fieldType = word
            fieldName = mustReadAtom(r)
            dimensions = nil
        }
        
        //symbol = mustReadSymbol(r)
        //if symbol != ';' { return nil, fmt.Errorf("expected ';'") }
        mustReadExpectedSymbol(r, ';', "end of template field")
        
        var member = TemplateMember{
            Name: fieldName,
            Type: fieldType,
            Dimensions: dimensions,
        }
        
        // is the type okay?
        if !member.isPrimitiveType() {
            _, exists := templates[fieldType]
            if !exists {
                panic(fmt.Errorf("unknown template field type %s (no forward or recursive references please)", fieldType))
            }
        }
        
        t.Members = append(t.Members, member)
    }
    
    return t
}

// decodeObject reads an object in a DirectX .x file. Note that the reader at this step has already
// consumed the leading identifier and successfully matched it to a template
func decodeObject(r *bufio.Reader, f *File, t *Template) *Data {
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
        mustUnreadByte(r)
        name = mustReadAtom(r)
        mustReadExpectedSymbol(r, '{', "start of object")
    }
    
    var data = &Data{Name: name, Spec: t}
    
    if len(name) > 0 {
        f.ReferencesByName[name] = data
    }
    
    // Read members first
    decodeMembers(r, data, data.Spec, 0, t.Members, f.templatesByName)
    
    // Read additional data blocks up to closing '}'
    for {
        // More data?
        if mustAcceptSymbol(r, '}') { break }
        if t.Mode == TemplateClosed { panic(fmt.Errorf("unexpected extra data in closed object type")) }
        
        if mustAcceptSymbol(r, '{') {
            data.appendReference(mustReadAtom(r))
            mustReadExpectedSymbol(r, '}', "end of reference")
        }
        
        var word = mustReadAtom(r)
        
        var t, ok = f.templatesByName[word]
        if !ok { panic(fmt.Errorf("unknown object type '%s'", word)) }
        
        data.appendChild(decodeObject(r, f, t))
    }
    
    return data
}

// decodeMembers decodes values according to a template. `target` is where the data is written to; spec specifies how
// the data is decoded. The spec need not be the target's spec e.g. when a field's type is itself a Template.
func decodeMembers(r *bufio.Reader, target *Data, spec *Template, suboffset int, members []TemplateMember, templates map[string]*Template) {
    for _, member := range members {
        decodeMemberValue(r, target, spec, suboffset, &member, templates)
    }
}

// decodeMemberValue decodes a value according to a template member, possibly an array of such
func decodeMemberValue(r *bufio.Reader, target *Data, spec *Template, suboffset int, member *TemplateMember, templates map[string]*Template) {
    
    if member.Dimensions == nil {
        decodeSingleValue(r, target, suboffset, member, -1, templates)
        
    } else if len(member.Dimensions) == 1 {
        // read a 1D array
        
        var ln, err = strconv.ParseInt(member.Dimensions[0], 10, 32)
        if err != nil {
            // array of variable length
            var len32 uint32
            
            var offset, size int
            _, offset, size, err = (&Data{Spec: spec}).GetNamedField(member.Dimensions[0], "DWORD", templates)
            // len32, err = target.GetNamedDWORD(member.Dimensions[0], templates)
            offset += suboffset
            
            if err != nil {
                panic(fmt.Errorf("unable to lookup variable dimension length for field %s referencing %s: %v", member.Name, member.Dimensions[0], err))
            }
            
            len32 = binary.LittleEndian.Uint32(target.Bytes[offset : offset + size])
            
            ln = int64(len32)
        }
        
        var arrayIndex = target.appendArray()
        
        for i := 0; i < int(ln); i++ {
            decodeSingleValue(r, target, suboffset, member, arrayIndex, templates)
            
            if i + 1 < int(ln) {
                mustReadExpectedSymbol(r, ',', "array item separator")
            }
        }
        
    } else {
        // Read a multidimensional array
        panic(fmt.Errorf("multidimensional arrays not yet supported"))
    }
    
    mustReadExpectedSymbol(r, ';', fmt.Sprintf("end of object member value while parsing %s.%s", target.SpecName(), member.Name))
}

// decodeMemberValue decodes a value according to a template member, but its not an array
func decodeSingleValue(r *bufio.Reader, data *Data, suboffset int, member *TemplateMember, arrayIndex int, templates map[string]*Template) {
    
    if member.isPrimitiveType() {
        switch member.Type {
        
            case "DWORD":   data.appendDWORD(   uint32( mustReadUint (r, 32)), arrayIndex)
            case "float":   data.appendFloat32( float32(mustReadFloat(r, 32)), arrayIndex)
            case "WORD":    data.appendWORD(    uint16( mustReadUint (r, 16)), arrayIndex)
            case "FLOAT":   data.appendFloat32( float32(mustReadFloat(r, 32)), arrayIndex)
            case "DOUBLE":  fallthrough
            case "CHAR":    fallthrough  // data.appendWORD(int8(mustReadInt(r, 8)), arrayIndex)
            case "UCHAR":   fallthrough // data.appendWORD(int8(mustReadInt(r, 8)), arrayIndex)
            case "BYTE":    fallthrough
            
            case "STRING":
                if arrayIndex >= 0 { panic(fmt.Errorf("string arrays not yet supported")) }
                mustReadExpectedSymbol(r, '"', "open quote")
                var s = mustReadString(r)
                mustReadExpectedSymbol(r, '"', "close quote")
                data.appendString(s, arrayIndex)
            default:
                panic(fmt.Sprintf("primitive type %s not handled (should never happen)", member.Type))
        }
    } else {
        // read a named data type (using a template)

        var spec = templates[member.Type] // guaranteed to work at this point
        
        decodeMembers(r, data, spec, len(data.Bytes), spec.Members, templates)
    }
}
