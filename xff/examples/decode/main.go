package main

import (
    "fmt"
    "os"
    "strings"
    
    "tawesoft.co.uk/go/xff"
)

// TODO TODO TODO
// Refactor xff.File to intern templates and parse-time float precision
// TODO TODO TODO

type SkinWeight struct {
    transformNodeName string
    vertexIndices []uint32
    weights []float32
    matrix [16]float32
}

func (s *SkinWeight) Decode(data *xff.Data, templates map[string]*xff.Template) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = r.(error)
        }
    }()
    
    s.transformNodeName = data.MustGetNamedSTRING("transformNodeName", templates)
    
    var length = int(data.MustGetNamedDWORD("nWeights", templates))
    var index, _, _ = data.MustGetNamedField("vertexIndices", "DWORD", templates)
    arrayIndex, _ := data.MustGetDWORD(index, templates)
    s.vertexIndices = make([]uint32, length)
    for i := 0; i < int(length); i++ {
        s.vertexIndices[i] = uint32(data.Arrays[arrayIndex][i*4]) // TODO unpack properly!
    }
    
    index, _, _ = data.MustGetNamedField("weights", "float", templates)
    arrayIndex, _ = data.MustGetDWORD(index, templates)
    s.weights = make([]float32, length)
    for i := 0; i < int(length); i++ {
        s.weights[i] = 1.0 * float32(data.Arrays[arrayIndex][i*4]) // TODO unpack properly!
    }
    
    var matrix = data.Children[0]
    
    index, _, _ = matrix.MustGetNamedField("matrix", "float", templates)
    arrayIndex, _ = matrix.MustGetDWORD(index, templates)
    
    for i := 0; i < 16; i++ {
         //var f, _ = data.MustGetFloat(index + i, templates)
         s.matrix[i] = 1.0 * float32(matrix.Arrays[arrayIndex][i*4]) // TODO unpack properly!
    }
    
    // etc.
    return nil
}


func printData(data *xff.Data, indent int, templates map[string]*xff.Template) {
    var indentStr = strings.Repeat("  ", indent)
    
    if data.Spec == nil {
        fmt.Printf("%sReference: '%s'\n", indentStr, data.Name)
    } else {
        fmt.Printf("%sData: name='%s' of type '%s', %d bytes\n", indentStr, data.Name, data.Spec.Name, len(data.Bytes))
        
        if data.Spec.UUID == xff.MustHexToUUID("6f0d123bbad24167a0d080224f25fabb") {
            fmt.Printf("Strings: %+v\n", data.Strings)
            
            // SkinWeights extension
            var weight SkinWeight
            var err = weight.Decode(data, templates)
            if err != nil { panic(err) }
            fmt.Printf("%s> SkinWeight.transformNodeName: '%s'\n", indentStr, weight.transformNodeName)
            fmt.Printf("%s> SkinWeight.vertexIndices (%d): %v\n", indentStr, len(weight.vertexIndices), weight.vertexIndices)
            fmt.Printf("%s> SkinWeight.weights (%d): %v\n", indentStr, len(weight.weights), weight.weights)
            fmt.Printf("%s> SkinWeight.matrix (%d): %v\n", indentStr, len(weight.matrix), weight.matrix)
        }
    }
    
    for i, child := range(data.Children) {
        printData(&child, indent+1, templates)
        if (i > 5) && (i + 1 < len(data.Children)) {
            // fmt.Printf("%s(and %d more children)...\n", indentStr, len(data.Children) - 5)
            // break
        }
    }
}

func main() {
    in, err := os.Open("assets/test/person.x")
    if err != nil { panic(err) }
    
    xfile, err := xff.Decode(in)
    if err != nil { panic(err) }
    
    for _, child := range(xfile.Children) {
        printData(&child, 0, xfile.Templates)
    }
    
    fmt.Printf("ReferencesByName:\n")
    for k, v := range xfile.ReferencesByName {
        fmt.Printf("  %s (UUID '%x')\n", k, v.UUID)
    }
    
    fmt.Printf("ReferencesByUUID:\n")
    for k, v := range xfile.ReferencesByUUID {
        fmt.Printf("  %x (UUID '%s')\n", k, v.Name)
    }
}
