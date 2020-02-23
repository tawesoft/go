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

// A custom DirectX template for storing whatever we want.
//
// Because it's defined at the top of the file, *any* application can parse objects using this type.
template MyCustomType {
    <0123456789ABCDEF0123456789ABCDEF> // example GUID

    STRING fooString;
    DWORD nThings;
    array DWORD intThings[nThings];
    array float floatThings[nThings];
    array Vector vectorThings[nThings];
    Matrix4x4 matrixOfThings;
}

Frame Root { // Frame is a builtin DirectX open (can be extended with anything) template
    MyCustomType {
        "Some Thing"; // fooString
        3; // nThings
        
        // intThings array
        1,2,3;

        // floatThings array
        0.1, 0.2, 0.3;

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
        AnimationKey { // Rotation
            0;
            1;
            0;4;-1.000000, 0.000000, 0.000000, 0.000000;;;
        }
    }
}
`
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
    vectorThings [][3]float32
    matrix       [16]float32
}

// We have to reference our custom type by UUID
var MyCustomTypeUUID = xff.MustHexToUUID("0123456789ABCDEF0123456789ABCDEF")

// MyCustomType2 is a custom DirectX template for storing whatever we want.
//
// This time, we're defining it externally instead of in the file.
//
// This does mean that only applications that can also do this will be able to parse the data though. That's why its
// better to define a template in the DirectX (.x) file itself, as we did earlier.
var MyCustomType2 = xff.Template{
    Name: "MyCustomType2",
    UUID: xff.MustHexToUUID("FEDCBA9876543210FEDCBA9876543210"),
    Mode: 'o', // open, can contain any child objects // TODO restricted template
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
func DecodeMyCustomObject(file *xff.File, data *xff.Data) *MyCustomObject {
    /*
     var obj = &MyCustomObject{
         fooString: data.MustGetNamedSTRING("fooString", file.Templates),
     }
     
    var length = int(data.MustGetNamedDWORD("nThings", file.Templates))
    var index, _, _ = data.MustGetNamedField("intThings", "DWORD", file.Templates)
    arrayIndex, _ := data.MustGetDWORD(index, file.Templates)
    obj.intThings = make([]uint32, length)
    for i := 0; i < int(length); i++ {
        obj.intThings[i] = uint32(data.Arrays[arrayIndex][i*4]) // TODO unpack properly!
    }
    
    return obj
    */
     
     return &MyCustomObject{}
}

// lets recurse over a data object and its children and print stuff out about it
func printData(file *xff.File, data *xff.Data, indent int) {
    var indentStr = strings.Repeat("    ", indent)
    
    if data.Spec == nil {
        fmt.Printf("%sReference: '%s' (to type %s)\n", indentStr, data.Name, data.SpecName())
    } else {
        fmt.Printf("%sChild data: name='%s' of type '%s', %d bytes\n", indentStr, data.Name, data.SpecName(), len(data.Bytes))
        
        if data.Spec.UUID == MyCustomTypeUUID {
            fmt.Printf("%sCool: its on object of our first custom type!\n", indentStr)
            fmt.Printf("%s%+v\n", indentStr, DecodeMyCustomObject(file, data))
        } else if data.Spec.UUID == MyCustomType2.UUID {
            fmt.Printf("%sCool: its an object of our second custom type!\n", indentStr)
        }
    }
    
    for _, child := range(data.Children) {
        printData(file, &child, indent+1)
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
    
    for _, child := range(file.Children) {
        printData(file, &child, 0)
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
