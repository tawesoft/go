package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    
    "tawesoft.co.uk/go/xff"
)

// For example's sake, lets invent some DirectX file with a made up custom template
var example = `xof 0303txt 0032

template XYZ {
    <00000000123412341234123412341234> // example GUID
    array DWORD xyz[3];
}

// A custom DirectX template for storing whatever we want.
//
// Because it's defined at the top of the file, *any* application can parse objects using this type.
template MyCustomType {
    <0123456789ABCDEF0123456789ABCDEF> // example GUID

    STRING fooString;
    DWORD nThings;
    array DWORD intThings[nThings];
    array float floatThings[nThings];
    array XYZ xyzThings[nThings];
    array Vector vectorThings[nThings];
    Matrix4x4 matrix;
}

Frame Root { // Frame is a builtin DirectX open (can be extended with anything) template
    MyCustomType {
        "Some Thing"; // fooString
        3; // nThings
        
        // intThings array
        1,2,3;

        // floatThings array
        0.1, 0.2, 0.3;

        // xyzThings array
        11, 12, 13;,
        21, 22, 23;,
        31, 33, 33;;

        // vectorThings array
        1.1; 1.2; 1.3;,
        2.1; 2.2; 2.3;,
        3.1; 3.2; 3.3;;

        // matrixOfThings
        1.000000, 0.000000, 0.000000, 0.000000,
        0.000000, 1.000000, 0.000000, 0.000000,
        0.000000, 0.000000, 1.000000, 0.000000,
        0.000000, 0.000000, 0.000000, 1.000000;;
    }
    MyCustomType ThisOneHasAName {
        "Another \"Thing\""; // fooString
        6; // nThings
        
        // intThings array
        9,8,7,6,5,4;

        // floatThings array
        0.6, 0.5, 0.4, 0.3, 0.2, 0.1;

        // xyzThings array
        11, 12, 13;,
        21, 22, 23;,
        31, 33, 33;,
        41, 42, 43;,
        51, 52, 53;,
        61, 63, 63;;

        // vectorThings array
        1.1; 1.2; 1.3;,
        2.1; 2.2; 2.3;,
        3.1; 3.2; 3.3;,
        4.1; 4.2; 4.3;,
        5.1; 5.2; 5.3;,
        6.1; 6.2; 6.3;;

        // matrixOfThings
        1.1, 1.2, 1.3, 1.4,
        2.1, 2.2, 2.3, 2.4,
        3.1, 3.2, 3.3, 3.4,
        4.1, 4.2, 4.3, 4.4;;
    }

    // We also define our own template type later in Go, called MyCustomType2.
    //
    // Only parsers that know about MyCustomType2 will be able to parse this next bit.
    MyCustomType2 {
        8; // nKeys
        1.1, -2.2, 3.3, -4.4, 5.5, -6.6, 7.7, -8.8; // floatThings array

        MyCustomType2 {
            1; // nKeys
            1.23456789; // floatThings array
        }

        MyCustomType2 {
            2; // nKeys
            1.23456789, -9.87654321; // floatThings array

            MyCustomType2 {
                1; // nKeys
                0.1234; // floatThings array
            }
        }
    }
}

// Here's some objects using DirectX built-in Templates.
//
// The xff decoder knows about these already, so nothing extra is needed to support this.
AnimationSet ArmatureAction {
    Animation {
        {ThisOneHasAName} // Reference to another object by name
        AnimationKey { // Scale
            1;
            4;
            0;3; 2.881100, 2.881100, 2.881100;;,
            1;3; 2.881100, 2.881100, 2.881100;;,
            2;3; 2.881100, 2.881100, 2.881100;;,
            3;3; 2.881100, 2.881100, 2.881100;;;
        }
    }
}
`
/*
AnimationSet ArmatureAction {
    Animation {
        {ThisOneHasAName} // Reference to another object by name
        AnimationKey { // Rotation
            0;
            1;
            0;4;-1.000000, 0.000000, 0.000000, 0.000000;;;
        }
    }
}
*/
// MyCustomObject is a native Go representation of an object of the MyCustomType template defined in the
// DirectX (.x) file.
//
// Note that it doesn't match the template exactly: we don't need nThings because Go has len() and we're using Go
// data types.
//
type MyCustomObject struct {
    fooString    string
    intThings    []uint32
    floatThings  []float32
    xyzThings    [][3]uint32
    vectorThings [][3]float32
    matrix       [16]float32
}

// We have to reference our custom type by UUID
var MyCustomTypeUUID = xff.MustHexToUUID("0123456789ABCDEF0123456789ABCDEF")

// We need to know how to get information from a parsed DirectX (.x) file format Object into our Go struct.
type MyCustomObjectAccessor struct {
    fooString    xff.FieldAccessor
    nThings      xff.FieldAccessor
    intThings    xff.FieldAccessor
    floatThings  xff.FieldAccessor
    xyzThings    xff.FieldAccessor
    vectorThings xff.FieldAccessor
    matrix       xff.FieldAccessor
    
    vectorX      xff.FieldAccessor
    vectorY      xff.FieldAccessor
    vectorZ      xff.FieldAccessor
}

// MyCustomType2 is a custom DirectX template for storing whatever we want.
//
// This time, we're defining it externally instead of in the file.
//
// This does mean that only applications that can also do this will be able to parse the data though. That's why its
// better to define a template in the DirectX (.x) file itself, as we did earlier.
var MyCustomType2 = xff.Template{
    Name: "MyCustomType2",
    UUID: xff.MustHexToUUID("FEDCBA9876543210FEDCBA9876543210"),
    Mode: xff.TemplateOpen, // can contain any type of child objects
    Members: []xff.TemplateMember{
        {
            Name:       "nKeys",
            Type:       "DWORD",
        },
        {
            Name:       "floatThings",
            Type:       "float",
            Dimensions: []string{"nKeys",},
        },
    },
}

// MyCustomObject2 is a native Go representation of an object of the MyCustomType2 template above. Because this one
// is an open template, the object can have child objects.
type MyCustomObject2 struct {
    floatThings  []float32
    
    children []MyCustomObject2
}

// Given a data object of type MyCustomType, let's decode it into a native Go type.
func DecodeMyCustomObject(accessor *MyCustomObjectAccessor, data *xff.Data) *MyCustomObject {

    //fmt.Printf("%+v\n", data)
    
    var nThings = accessor.nThings.MustGetDWORD(data)
    
    // xyzThings is a bit more difficult as an array of arrays...
    var xyzThings = make([][3]uint32, nThings)
    var xyzThingsOuterArray = accessor.xyzThings.MustGetArray(data)
    for i := 0; i < int(nThings); i++ {
        var xyzThingsInnerArray = xyzThingsOuterArray.MustGetArray(data, i)
        for j := 0; j < 3; j++ {
            xyzThings[i][j] = xyzThingsInnerArray.MustGetDWORD(data, j)
        }
    }
    
    // vectorThings is a bit more difficult as an array of templates...
    // FIXME this works, but we'd like to actually get a rich object back e.g. to get x, y, z
    var vectorThings = make([][3]float32, nThings)
    var vectorThingsArray = accessor.vectorThings.MustGetArray(data)
    for i := 0; i < int(nThings); i++ {
        vectorThings[i][0] = vectorThingsArray.MustGetFloat(data, (i*3))
        vectorThings[i][1] = vectorThingsArray.MustGetFloat(data, (i*3)+1)
        vectorThings[i][2] = vectorThingsArray.MustGetFloat(data, (i*3)+2)
    }
    
     // but the other fields are easy!
     var obj = &MyCustomObject{
         fooString: accessor.fooString.MustGetSTRING(data),
         intThings: accessor.intThings.MustGetDWORDArray(data, int(nThings)),
         floatThings: accessor.floatThings.MustGetFloatArray(data, int(nThings)),
         xyzThings: xyzThings,
         vectorThings: vectorThings,
     }
     
     /*
     var matrix = make([]float32, 16)
     var matrixArray = accessor.matrix.MustGetField(data)
     for i := 0; i < 16; i++ {
         matrix[i] = matrixArray.MustGetFloat(data, i)
     }
     copy(obj.matrix[:], matrix)
      */
     
     /*
    var length = int(data.MustGetNamedDWORD("nThings", file.Templates))
    var index, _, _ = data.MustGetNamedField("intThings", "DWORD", file.Templates)
    arrayIndex, _ := data.MustGetDWORD(index, file.Templates)
    obj.intThings = make([]uint32, length)
    for i := 0; i < int(length); i++ {
        obj.intThings[i] = uint32(data.Arrays[arrayIndex][i*4]) // TODO unpack properly!
    }
    
    return obj
    */
     
     return obj
}

// lets recurse over a data object and its children and print stuff out about it
func printData(file *xff.File, accessor *MyCustomObjectAccessor, data *xff.Data, indent int) {
    var indentStr = strings.Repeat("    ", indent)
    
    if data.IsReference() {
        fmt.Printf("%sReference: '%s' (to an object of type %s)\n", indentStr, data.Name, file.ReferencesByName[data.Name].SpecName())
        
    } else {
        fmt.Printf("%sChild data: name='%s' of type '%s', %d bytes\n", indentStr, data.Name, data.SpecName(), len(data.Bytes))
        
        // There are three places a template can be stored, so here's how to select against all of them
        //
        // You might be tempted to compare against the template name or pointer, but please don't: the UUID is
        // used because a template might appear twice - e.g. once in file and repeated as a built-in - so only
        // the UUID can be relied on to match against.
        
        if data.Spec.UUID == MyCustomTypeUUID { // match a constant against the Template UUID in file
            fmt.Printf("%sCool: its on object of our custom type (the type we defined inline in the DirectX file)!\n", indentStr)
            
            if data.Name == "" {
                fmt.Printf("%s%+v\n", indentStr, DecodeMyCustomObject(accessor, data))
            }
            
        } else if data.Spec.UUID == MyCustomType2.UUID { // match Template UUID specified in our custom xff.Template
            fmt.Printf("%sCool: its an object of our custom type (the type we defined in Go)!\n", indentStr)
            
        } else if data.Spec.UUID == xff.TemplateAnimationSet.UUID { // match built-in xff.Template UUID
            fmt.Printf("%sCool: is an animation set type (a type defined as a built-in)!\n", indentStr)
        }
    }
    
    for _, child := range(data.Children) {
        printData(file, accessor, &child, indent+1)
    }
}

// resetwd sets the working directory to the directory of the executable
func resetwd() {
    var exdir, err = os.Executable()
    if err != nil {
        panic(err)
    }
    err = os.Chdir(filepath.Dir(exdir))
    if err != nil { panic(err) }
}

func main() {
    resetwd()
    
    var extraTemplates = []*xff.Template{
        &MyCustomType2,
    }
    
    file, err := xff.Decode(strings.NewReader(example), extraTemplates)
    if err != nil {
        fmt.Printf("Error: %s\n", err.Error())
        return
    }
    
    var accessor = MyCustomObjectAccessor{
        fooString:    file.MustGetFieldAccessor(MyCustomTypeUUID, "fooString"),
        nThings:      file.MustGetFieldAccessor(MyCustomTypeUUID, "nThings"),
        intThings:    file.MustGetFieldAccessor(MyCustomTypeUUID, "intThings"),
        floatThings:  file.MustGetFieldAccessor(MyCustomTypeUUID, "floatThings"),
        xyzThings:    file.MustGetFieldAccessor(MyCustomTypeUUID, "xyzThings"),
        vectorThings: file.MustGetFieldAccessor(MyCustomTypeUUID, "vectorThings"),
        matrix:       file.MustGetFieldAccessor(MyCustomTypeUUID, "matrix"),
        
        vectorX:      file.MustGetFieldAccessor(xff.TemplateVector.UUID, "x"),
        vectorY:      file.MustGetFieldAccessor(xff.TemplateVector.UUID, "y"),
        vectorZ:      file.MustGetFieldAccessor(xff.TemplateVector.UUID, "z"),
    }
    
    for _, child := range(file.Children) {
        printData(file, &accessor, &child, 0)
    }
    
    fmt.Printf("file.ReferencesByName:\n")
    for k, v := range file.ReferencesByName {
        fmt.Printf("  %s (UUID '%x')\n", k, v.UUID)
    }
    
    fmt.Printf("file.ReferencesByUUID:\n")
    for k, v := range file.ReferencesByUUID {
        fmt.Printf("  %x (UUID '%s')\n", k, v.Name)
    }
}
