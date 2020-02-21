package xff

import (
    "encoding/binary"
    "fmt"
    "math"
)

type Data struct {
    Name string
    Spec *Template
    Bytes []byte
}

func (b Data) GetField(index int) (offset int, size int, err error) {
    for i := 0; i < len(b.Spec.Members); i++ {
        offset += size
        size = b.Spec.Members[i].Size()
        if i == index { return offset, size, nil }
    }
    
    return 0, 0, fmt.Errorf("invalid reference to field by index %d", index)
}

func (b Data) GetNamedField(name string) (buf []byte, offset int, size int, err error) {
    for i := 0; i < len(b.Spec.Members); i++ {
        offset += size
        size = b.Spec.Members[i].Size()
        if b.Spec.Members[i].Name == name { return b.Bytes, offset, size, nil }
    }
    
    return nil, 0, 0, fmt.Errorf("invalid reference to named field %s in %s", name, b.Spec.Name)
}

func (b Data) GetDWORD(offset int) (uint32, int, error) {
    var o, size, err = b.GetField(offset)
    if err != nil { return 0, 0, err }
    return binary.LittleEndian.Uint32(b.Bytes[o : o + size]), 4, nil
}

func (b Data) GetNamedDWORD(name string) (uint32, error) {
    var buf, o, size, err = b.GetNamedField(name)
    if err != nil { return 0, err }
    return binary.LittleEndian.Uint32(buf[o : o + size]), nil
}

func (b *Data) AppendDWORD(value uint32) {
    var bytes [4]byte
    binary.LittleEndian.PutUint32(bytes[:], value)
    b.Bytes = append(b.Bytes, bytes[:]...)
}

func (b *Data) AppendFloat32(value float32) {
    var bytes [4]byte
    binary.LittleEndian.PutUint32(bytes[:], math.Float32bits(value))
    b.Bytes = append(b.Bytes, bytes[:]...)
}

func (b *Data) AppendWORD(value uint16) {
    var bytes [2]byte
    binary.LittleEndian.PutUint16(bytes[:], value)
    b.Bytes = append(b.Bytes, bytes[:]...)
}

