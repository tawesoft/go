package xff

import (
    "encoding/binary"
    "encoding/hex"
    "fmt"
    "math"
    "strings"
)

// UUID is a 128 bit (16 bytes, or 32 hexadecimal digits) ID used in the DirectX (.x) file format to uniquely identify
// templates and optionally objects.
type UUID_t [16]byte

// MustHexToUUID returns a 128 bit UUID from a hexadecimal string, or panics on error (making this version of the
// function better for defining exported constant type values). Hyphens in the string are ignored.
func MustHexToUUID(hexstr string) (uuid UUID_t) {
    uuid, err := HexToUUID(hexstr)
    if err != nil { panic(err) }
    return uuid
}

// MustHexToUUID returns a 128 bit UUID from a hexadecimal string. Hyphens in the string are ignored.
func HexToUUID(hexstr string) (uuid UUID_t, err error) {
    hexstr = strings.Replace(hexstr, "-", "", -1)
    
    decoded, err := hex.Decode(uuid[:], []byte(hexstr))
    
    if err != nil { return uuid, fmt.Errorf("UUID string '%s' is not a valid hexadecimal string: %v", hexstr, err) }
    if (decoded != 16) { return uuid, fmt.Errorf("UUID string '%s' must be exactly 32 hexadecimal digits long", hexstr) }
    
    return uuid, nil
}

// File is represents data in a decoded DirectX (.x) file format
type File struct {
    // Children are zero or more objects
    Children []Data // not nil
    
    // ReferencesByName map object names to objects
    ReferencesByName map[string]*Data
    
    // ReferencesByUUID map object UUIDs to objects
    ReferencesByUUID map[UUID_t]*Data
    
    // floatSize (32 or 64) is used in the binary encoding
    
    templatesByName map[string]*Template
}

func (f *File) appendChild(data *Data) {
    f.Children = append(f.Children, *data)
}

// Data is a decoded object in a DirectX (.x) file format. Each object has a type (DirectX calls this a Template),
// and some values according to that type, and (if the Template is not closed) child objects of any (if the template is
// open) or a restricted set of (if the template is restricted) template type. An object's values can be primitive
// types (like the DWORD, DirectX's version of a uint32), an array of primitive types, a typed object, or an array of
// typed objects, and child objects.
type Data struct {
    
    // Name is the optional name given to the data object. It might just be an empty string.
    Name string
    
    // UUID is the optional UUID given to a data object. It might just be zero bytes.
    // Currently this isn't implemented, just because I haven't got any examples to test with.
    UUID UUID_t
    
    // Soec is the Template type of the object. If nil, the object is just a reference: use the Name field to
    // lookup the referenced object.
    Spec *Template // may be nil
    
    // TODO make these internal
    Bytes []byte
    ArrayData [][]byte // TODO make this flat
    Strings []string
    Children []Data
}

// IsReference returns true if the data object is not a fully instantiated object but instead a reference to another
// object (either by Name or UUID) that may or may not exist and may or may not have been decoded yet. If it is a
// reference, the Spec Template field is a nil pointer, because it doesn't have a Template yet.
func (d *Data) IsReference() bool {
    return d.Spec == nil
}

// SpecName returns a data object's Template's name (useful for debugging) or, if the data object doesn't have a
// Template because instead of being a fully instantiated object it's a reference to another named object (which may
// or may not have been decoded at this point), it returns an empty string. This saves checking for the nil pointer.
func (d *Data) SpecName() string {
    if d.Spec == nil { return "" }
    return d.Spec.Name
}

// TODO get this on a template, not the data!
// getNamedField returns the index (e.g. "the 2nd field"; start counting at zero), offset (bytes) into the packed data,
// and size (bytes) in the packed data of a data object according to a field of a certain name.
func (f *File) getNamedField(data *Data, fieldName string, fieldType string) (index int, offset int, size int, err error) {
    
    for i := 0; i < len(data.Spec.Members); i++ {
        offset += size
        var member = data.Spec.Members[i]
        size = data.Spec.Members[i].size(f.templatesByName)
        
        if member.Name == fieldName {
            if member.Type != fieldType {
                return 0, 0, 0, fmt.Errorf("invalid access to named field %s of object %s and type %s as type %s",
                    fieldName, data.SpecName(), member.Type, fieldType)
            }
            
            return i, offset, size, nil
        }
    }
    
    return 0, 0, 0, fmt.Errorf("named field %s of object %s not found", fieldName, data.SpecName())
}




// ---- OLD VERSIONS BELOW ----

// GetField returns the offset (for GetFloat, GetDWORD, etc) and size (for incrementing offsets in sequential
// access) of a data field in a data block by a known index (e.g. "the second field"; start counting at zero).
//
// Note that GetNamedField (and GetNamedFloat, GetNameDDWORD, etc.) should be preferred where possible because
// these check for type errors.
func (d *Data) GetField(index int, templates map[string]*Template) (offset int, size int, err error) {
    for i := 0; i < len(d.Spec.Members); i++ {
        offset += size
        size = d.Spec.Members[i].size(templates)
        if i == index { return offset, size, nil }
    }
    
    return 0, 0, fmt.Errorf("invalid reference to %s field at index %d", d.SpecName(), index)
}

// GetNamedField returns the index (for GetField), offset (for GetFloat, GetDWORD, etc), size (for incrementing
// offsets in sequential access) of a data field in a data block by name
func (b *Data) GetNamedField(fieldName string, fieldType string, templates map[string]*Template) (index int, offset int, size int, err error) {
    for i := 0; i < len(b.Spec.Members); i++ {
        offset += size
        var member = b.Spec.Members[i]
        size = b.Spec.Members[i].size(templates)
        
        if member.Name == fieldName {
            if member.Type != fieldType {
                return 0, 0, 0, fmt.Errorf("invalid type access %s for named field %s of type %s",
                    fieldType, fieldName, member.Type)
            }
            
            return i, offset, size, nil
        }
    }
    
    return 0, 0, 0, fmt.Errorf("invalid reference to %s named field %s", b.SpecName(), fieldName)
}

// MustGetNamedField is like GetNamedField, but panics on error. This simplifies error handling by enabling the caller
// // to recover over a batch of closely related function calls.
func (b *Data) MustGetNamedField(fieldName string, fieldType string, templates map[string]*Template) (index int, offset int, size int) {
    index, offset, size, err := b.GetNamedField(fieldName, fieldType, templates)
    if err != nil { panic(err) }
    return index, offset, size
}

// GetDWORD unpacks a DWORD field at a given offset. Use the returned size to advance the offset to the start of
// the next field. Note that this is not checked for type errors: GetNamedDWORD is preferred.
func (b *Data) GetDWORD(offset int, templates map[string]*Template) (value uint32, size int, err error) {
    offset, size, err = b.GetField(offset, templates)
    if err != nil { return 0, 0, err }
    return binary.LittleEndian.Uint32(b.Bytes[offset : offset + size]), 4, nil
}

// MustGetDWORD is like GetDWORD, but panics on error. This simplifies error handling by enabling the caller
// to recover over a batch of closely related function calls.
func (b *Data) MustGetDWORD(offset int, templates map[string]*Template) (value uint32, size int) {
    value, size, err := b.GetDWORD(offset, templates)
    if err != nil { panic(err) }
    return value, size
}

// GetNamedDWORD unpacks a DWORD field by a given field name.
func (b *Data) GetNamedDWORD(name string, templates map[string]*Template) (uint32, error) {
    var _, offset, size, err = b.GetNamedField(name, "DWORD", templates)
    if err != nil { return 0, err }
    var value = binary.LittleEndian.Uint32(b.Bytes[offset : offset + size])
    return value, nil
}

// MustGetNamedDWORD is like GetNamedDWORD, but panics on error. This simplifies error handling by enabling the caller
// to recover over a batch of closely related function calls.
func (b *Data) MustGetNamedDWORD(name string, templates map[string]*Template) uint32 {
    var result, err = b.GetNamedDWORD(name, templates)
    if err != nil { panic(err) }
    return result
}

// GetFloat unpacks a float field at a given offset. The size of the float (32 or 64 bit) depends on the
// format specified in the DirectX (.x) file and corresponds to the lowercase "float" datatype. For the explicitly
// sized types, see GetFLOAT and GetDOUBLE. Note that this is not checked for type errors:
// GetNamedFloat is preferred.
func (b *Data) GetFloat(offset int, templates map[string]*Template) (value float64, size int, err error) {
    offset, size, err = b.GetField(offset, templates)
    if err != nil { return 0, 0, err }
    return float64(math.Float32frombits(binary.LittleEndian.Uint32(b.Bytes[offset : offset + size]))), 4, nil
}

// MustGetFloat is like GetFloat, but panics on error. This simplifies error handling by enabling the caller
// to recover over a batch of closely related function calls.
func (b *Data) MustGetFloat(offset int, templates map[string]*Template) (value float64, size int) {
    value, size, err := b.GetFloat(offset, templates)
    if err != nil { panic(err) }
    return value, size
}

// GetNamedFloat unpacks a float field by a given field name. The size of the float (32 or 64 bit) depends on the
// format specified in the DirectX (.x) file and corresponds to the lowercase "float" datatype. For the explicitly
// sized types, see GetNamedFLOAT and GetNamedDOUBLE.
func (b *Data) GetNamedFloat(name string, templates map[string]*Template) (float64, error) {
    var _, offset, size, err = b.GetNamedField(name, "float", templates)
    if err != nil { return 0, err }
    return float64(math.Float32frombits(binary.LittleEndian.Uint32(b.Bytes[offset : offset + size]))), nil
}

// MustGetNamedFloat is like GetNamedFloat, but panics on error. This simplifies error handling by enabling the caller
// to recover over a batch of closely related function calls.
func (b *Data) MustGetNamedFloat(name string, templates map[string]*Template) float64 {
    var result, err = b.GetNamedFloat(name, templates)
    if err != nil { panic(err) }
    return result
}

// GetSTRING unpacks a STRING field at a given offset. Use the returned size to advance the offset to the start of
// the next field. Note that this is not checked for type errors: GetNamedSTRING is preferred.
func (b *Data) GetSTRING(offset int, templates map[string]*Template) (value string, size int, err error) {
    var index uint32
    index, size, err = b.GetDWORD(offset, templates)
    if err != nil { return "", 0, err }
    return b.Strings[index], 4, nil
}

// GetNamedSTRING unpacks a STRING field by a given field name.
func (b *Data) GetNamedSTRING(name string, templates map[string]*Template) (string, error) {
    var index, _, _, err = b.GetNamedField(name, "STRING", templates)
    if err != nil { return "", err }
    
    var value string
    value, _, err = b.GetSTRING(index, templates)
    return value, err
}

// MustGetNamedSTRING is like GetNamedSTRING, but panics on error. This simplifies error handling by enabling the caller
// to recover over a batch of closely related function calls.
func (b *Data) MustGetNamedSTRING(name string, templates map[string]*Template) string {
    var result, err = b.GetNamedSTRING(name, templates)
    if err != nil { panic(err) }
    return result
}

func (b *Data) appendChild(data *Data) {
    b.Children = append(b.Children, *data)
}

func (b *Data) appendWORD(value uint16, arrayIndex int) {
    var buf *[]byte
    var bytes [2]byte
    binary.LittleEndian.PutUint16(bytes[:], value)

    if arrayIndex < 0 {
        buf = &b.Bytes
    } else {
        buf = &b.ArrayData[arrayIndex]
    }
    *buf = append(*buf, bytes[:]...)
}

func (b *Data) appendDWORD(value uint32, arrayIndex int) {
    var buf *[]byte
    var bytes [4]byte
    binary.LittleEndian.PutUint32(bytes[:], value)
    
    if arrayIndex < 0 {
        buf = &b.Bytes
    } else {
        buf = &b.ArrayData[arrayIndex]
    }
    *buf = append(*buf, bytes[:]...)
}

func (b *Data) appendFloat32(value float32, arrayIndex int) {
    var buf *[]byte
    var bytes [4]byte
    binary.LittleEndian.PutUint32(bytes[:], math.Float32bits(value))

    if arrayIndex < 0 {
        buf = &b.Bytes
    } else {
        buf = &b.ArrayData[arrayIndex]
    }
    *buf = append(*buf, bytes[:]...)
}

func (b *Data) appendString(value string, arrayIndex int) {
    b.appendDWORD(uint32(len(b.Strings)), arrayIndex)
    b.Strings = append(b.Strings, value)
}

func (b *Data) appendArray() (index int) {
    b.ArrayData = append(b.ArrayData, nil)
    var length = len(b.ArrayData)
    b.appendDWORD(uint32(length - 1), -1)
    return length - 1
}
