package xff

// xff/decode.go decodes the high level structure of a DirectX .x file

import (
    "bufio"
    "encoding/binary"
    "fmt"
    "io"
    "strconv"
    "strings"
)

func Decode(r io.Reader) (err error) {
    var fp = bufio.NewReader(r)
    
    var format byte // 't'/text or 'b'/binary
    var floatSize int // 32 or 64
    
    format, floatSize, err = decodeHeader(fp)
    if err != nil { return fmt.Errorf("error decoding header: %v", err) }
    
    if format != 't' { return fmt.Errorf("format not implemented") }
    if floatSize != 32 { return fmt.Errorf("float size not implemented") }
    
    var templates = make(map[string]*Template)
    for k, v := range Templates {
        templates[k] = v
    }
    
    for {
        var word string
        word, err = readAtom(fp)
        if err != nil { break; }
    
        if word == "template" {
            var t *Template
            t, err = decodeTemplate(fp)
            if err != nil { return err }
            templates[t.Name] = t
        } else {
            var t, ok = templates[word]
            if !ok { return fmt.Errorf("unknown identifier %s", word) }
            
            err = decodeDataBlock(fp, t, templates)
            if err != nil { return err }
        }
    }
    
    if err == io.EOF { err = nil }
    
    return err
}

// DecodeHeader reads the header of a DirectX .x file and returns the format and floatSize on success.
//
// * format will be either 't' (text) or 'b' (binary).
//
// * floatSize will be either 32 or 64.
//
func decodeHeader(r io.Reader) (format byte, floatSize int, err error) {
    /*
    https://docs.microsoft.com/en-us/windows/win32/direct3d9/reserved-words--header--and-comments
    
    The variable-length header is compulsory and must be at the beginning of the data stream. The header contains the
    following data.
    
    Type                Required    Size (in bytes) Value   Description
    Magic Number        x           4               xof
    Version Number      x           2               03      Major version 3
                                                    03      Minor version 3
    Format Type         x           4               txt     Text File
                                                    bin     Binary file
                                                    tzip    MSZip compressed text file
                                                    bzip    MSZip compressed binary file
    Float Size          x       "0064"                      64-bit floats
                        x       "0032"                      32-bit floats
    */
    
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
    defer func() {
        if r := recover(); r != nil {
            err = r.(error)
        }
    }()
    
    t = &Template{}
    t.Name = mustReadAtom(r)
    if mustReadSymbol(r) != '{' { return nil, fmt.Errorf("expected {") }
    if mustReadSymbol(r) != '<' { return nil, fmt.Errorf("expected <") }
    t.UUID = strings.ToUpper(mustReadAtom(r))
    if mustReadSymbol(r) != '>' { return nil, fmt.Errorf("expected >") }
    t.Mode = 'c' // 'c'/closed
    
    t.Members = make([]TemplateMember, 0)
    
    for {
        var fieldName, fieldType string
        var dimensions []string
        
        // Try to end the block
        var symbol = mustReadSymbol(r)
        if symbol == '}' { break; }
        r.UnreadByte()
        
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
                r.UnreadByte()
                if symbol == ';' { break; }
            }
        } else {
            fieldType = word
            fieldName = mustReadAtom(r)
            dimensions = nil
        }
        
        symbol = mustReadSymbol(r)
        if symbol != ';' { return nil, fmt.Errorf("expected ;") }
        
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
func decodeDataBlock(r *bufio.Reader, t *Template, templates map[string]*Template) (err error) {
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
        r.UnreadByte()
        name = mustReadAtom(r)
        if mustReadSymbol(r) != '{' {
            return fmt.Errorf("expected {")
        }
    }
    
    var data = &Data{Name: name, Spec: t}
    
    // Read members first
    for _, member := range t.Members {
        err = decodeValue(r, data, &member, templates)
        if err != nil { return err }
    }
    
    // Read additional data blocks up to closing '}'
    for {
        // More data?
        if mustReadSymbol(r) == '}' { break }
        r.UnreadByte()
        if t.Mode == 'c' { return fmt.Errorf("unexpected extra data in closed data type") }
        
        if mustReadSymbol(r) == '{' {
            // data reference
            var reference = mustReadAtom(r)
            fmt.Printf("got reference %s\n", reference)
            if mustReadSymbol(r) != '}' { return fmt.Errorf("expected } to end reference") }
        } else {
            r.UnreadByte()
        }
        
        var word string
        word = mustReadAtom(r)
        
        var t, ok = templates[word]
        if !ok { return fmt.Errorf("unknown identifier '%s'", word) }
        
        err = decodeDataBlock(r, t, templates)
        if err != nil { return err }
    }
    
    return err
}

func decodeValue(r *bufio.Reader, data *Data, member *TemplateMember, templates map[string]*Template) (err error) {
    
    if member.Dimensions == nil {
        // read a single value
        
        err = decodeSingleValue(r, data, member, templates)
        if err != nil { return err }
        
    } else if len(member.Dimensions) == 1 {
        // read a 1D array
        
        var len, err = strconv.ParseInt(member.Dimensions[0], 10, 32)
        if err != nil {
            // of variable length
            var len32 uint32
            len32, err = data.GetNamedDWORD(member.Dimensions[0])
            if err != nil { return fmt.Errorf("unable to lookup variable dimension length for field %s referencing %s: %v", member.Name, member.Dimensions[0], err) }
            len = int64(len32)
        }
        
        for i := 0; i < int(len); i++ {
            err = decodeSingleValue(r, data, member, templates)
            if err != nil { return err }
            
            if i + 1 < int(len) {
                if mustReadSymbol(r) != ',' { return fmt.Errorf("expected , while reading array") }
            }
        }
        
    } else {
        // Read a multidimensional array
        return fmt.Errorf("multidimensional arrays not yet supported")
    }
    
    if mustReadSymbol(r) != ';' { return fmt.Errorf("expected ;") }
    
    return nil
}

func decodeSingleValue(r *bufio.Reader, data *Data, member *TemplateMember, templates map[string]*Template) (err error) {
    
    if member.PrimitiveType() {
        switch member.Type {
            case "DWORD":
                var dword int64
                dword, err = strconv.ParseInt(mustReadAtom(r), 10, 32)
                if dword < 0 { dword = -dword }
                data.AppendDWORD(uint32(dword))
            case "float":
                var float float64
                float, err = strconv.ParseFloat(mustReadAtom(r), 32)
                data.AppendFloat32(float32(float))
            case "WORD":
                var dword int64
                dword, err = strconv.ParseInt(mustReadAtom(r), 10, 16)
                if dword < 0 { dword = -dword }
                data.AppendWORD(uint16(dword))
            case "FLOAT":  fallthrough
            case "DOUBLE": fallthrough
            case "CHAR":   fallthrough
            case "UCHAR":  fallthrough
            case "BYTE":   fallthrough
            case "STRING":
                // TODO special mustReadString that takes strings with escaped quotes
                if mustReadSymbol(r) != '"' { return fmt.Errorf("expected open quote") }
                var s = mustReadAtom(r)
                fmt.Printf("Read string '%s'\n", s)
                if mustReadSymbol(r) != '"' { return fmt.Errorf("expected close quote") }
            default:
                panic(fmt.Sprintf("primitive type %s not handled (should never happen)", member.Type))
        }
    } else {
        // read a named data type (using a template)
        var subt, ok = templates[member.Type]
        if !ok { return fmt.Errorf("unrecognised named data type %s for %s in %s", member.Type, member.Name, data.Spec.Name) }
        
        var subdata = &Data{Spec: subt}
        
        for _, member := range subt.Members {
            err = decodeValue(r, subdata, &member, templates)
            if err != nil { return err }
        }
    }
    
    return nil
}
