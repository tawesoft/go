package xff

type TemplateMember struct {
    Name string
    Type string
    Dimensions []string
}

func (m *TemplateMember) isPrimitiveType() bool {
    switch m.Type {
        case "WORD":   return true
        case "DWORD":  return true
        case "FLOAT":  return true
        case "DOUBLE": return true
        case "CHAR":   return true
        case "UCHAR":  return true
        case "BYTE":   return true
        case "STRING": return true
        case "float":  return true
        default:       return false
    }
}

func (m *TemplateMember) size(templates map[string]*Template) int {
    if m.Dimensions != nil { return 4 } // array
    switch m.Type {
        case "WORD":   return 2
        case "DWORD":  return 4
        case "FLOAT":  return 4
        case "DOUBLE": return 8
        case "CHAR":   return 1
        case "UCHAR":  return 1
        case "BYTE":   return 1
        case "STRING": return 4
        case "float":  return 4
        default:
            // templates[m.Type] is guaranteed to succeed at this point
            return templates[m.Type].size(templates)
    }
}

// A DirectX (.x) file format Template defines the structure of an object. Some Templates are built-in, some Templates
// are defined ad-hoc at the start of some files, and they can even be passed in to the Decode function at runtime.
// A Template should be considered constant and read-only once instantiated, because a decoded object will need to hold
// a reference to the Template used. The preferred way to identify which Template is associated with a decoded object
// is to use the UUID not pointer equality because a template may be defined by either of the three methods listed
// above. You can get a UUID of a built-in template easily e.g. TemplateAnimation.UUID.
type Template struct {
    Name string
    UUID UUID_t
    Mode byte // 'o'pen, 'c'losed, 'r'estricted
    Members []TemplateMember
}

func (t *Template) size(templates map[string]*Template) (acc int) {
    for _, member := range t.Members {
        acc += member.size(templates)
    }
    return acc
}

var TemplateAnimation = Template{
    Name: "Animation",
    UUID: MustHexToUUID("3D82AB4F-62DA-11cf-AB39-0020AF71E433"),
    Mode: 'o', // open
}

var TemplateAnimationKey = Template{
    Name: "AnimationKey",
    UUID: MustHexToUUID("10DD46A8-775B-11CF-8F52-0040333594A3"),
    Mode: 'c', // closed
    Members: []TemplateMember{
        {
            Name:       "keyType",
            Type:       "DWORD",
        },
        {
            Name:       "nKeys",
            Type:       "DWORD",
        },
        {
            Name:       "keys",
            Type:       "TimedFloatKeys",
            Dimensions: []string{"nKeys",},
        },
    },
}

// TemplateAnimationSet is a DirectX (.x) file Template for an AnimationSet object that contains one or more
// Animation objects
var TemplateAnimationSet = Template{
    Name: "AnimationSet",
    UUID: MustHexToUUID("3D82AB50-62DA-11cf-AB39-0020AF71E433"),
    Mode: 'r', // restricted to Animation objects (TODO)
}

// TemplateCoords2D is a DirectX (.x) file Template for a Coords2D object that contains a (u, v) component
var TemplateCoords2D = Template{
    Name: "Coords2D",
    UUID: MustHexToUUID("F6F23F44-7686-11cf-8F52-0040333594A3"),
    Mode: 'c', // closed
    Members: []TemplateMember{
        {
            Name:       "u",
            Type:       "float",
        },
        {
            Name:       "v",
            Type:       "float",
        },
    },
}

var TemplateFloatKeys = Template{
    Name: "FloatKeys",
    UUID: MustHexToUUID("10DD46A9-775B-11cf-8F52-0040333594A3"),
    Mode: 'c', // closed
    Members: []TemplateMember{
        {
            Name:       "nValues",
            Type:       "DWORD",
        },
        {
            Name:       "values",
            Type:       "float",
            Dimensions: []string{"nValues",},
        },
    },
}

var TemplateFrame = Template{
    Name: "Frame",
    UUID: MustHexToUUID("3D82AB46-62DA-11CF-AB39-0020AF71E433"),
    Mode: 'o', // open
}

var TemplateFrameTransformMatrix = Template{
    Name: "FrameTransformMatrix",
    UUID: MustHexToUUID("F6F23F41-7686-11cf-8F52-0040333594A3"),
    Mode: 'c', // closed
    Members: []TemplateMember{
        {
            Name:       "frameMatrix",
            Type:       "Matrix4x4",
        },
    },
}

var TemplateMatrix4x4 = Template{
    Name: "Matrix4x4",
    UUID: MustHexToUUID("F6F23F45-7686-11cf-8F52-0040333594A3"),
    Mode: 'c', // closed
    Members: []TemplateMember{
        {
            Name:       "matrix",
            Type:       "float",
            Dimensions: []string{"16",},
        },
    },
}

var TemplateMesh = Template{
    Name: "Mesh",
    UUID: MustHexToUUID("3D82AB44-62DA-11CF-AB39-0020AF71E433"),
    Mode: 'o', // open
    Members: []TemplateMember{
        {
            Name:       "nVertices",
            Type:       "DWORD",
        },
        {
            Name:       "vertices",
            Type:       "Vector",
            Dimensions: []string{"nVertices",},
        },
        {
            Name:       "nFaces",
            Type:       "DWORD",
        },
        {
            Name:       "faces",
            Type:       "MeshFace",
            Dimensions: []string{"nFaces",},
        },
    },
}

var TemplateMeshFace = Template{
    Name: "MeshFace",
    UUID: MustHexToUUID("3D82AB5F-62DA-11cf-AB39-0020AF71E433"),
    Mode: 'c', // closed
    Members: []TemplateMember{
        {
            Name:       "nFaceVertexIndices",
            Type:       "DWORD",
        },
        {
            Name:       "faceVertexIndices",
            Type:       "DWORD",
            Dimensions: []string{"nFaceVertexIndices",},
        },
    },
}

var TemplateMeshNormals = Template{
    Name: "MeshNormals",
    UUID: MustHexToUUID("F6F23F43-7686-11cf-8F52-0040333594A3"),
    Mode: 'c', // closed
    Members: []TemplateMember{
        {
            Name:       "nNormals",
            Type:       "DWORD",
        },
        {
            Name:       "normals",
            Type:       "Vector",
            Dimensions: []string{"nNormals",},
        },
        {
            Name:       "nFaceNormals",
            Type:       "DWORD",
        },
        {
            Name:       "faceNormals",
            Type:       "MeshFace",
            Dimensions: []string{"nFaceNormals",},
        },
    },
}

var TemplateMeshTextureCoords = Template{
    Name: "MeshTextureCoords",
    UUID: MustHexToUUID("F6F23F40-7686-11cf-8F52-0040333594A3"),
    Mode: 'c', // closed
    Members: []TemplateMember{
        {
            Name:       "nTextureCoords",
            Type:       "DWORD",
        },
        {
            Name:       "textureCoords",
            Type:       "Coords2D",
            Dimensions: []string{"nTextureCoords"},
        },
    },
}

var TemplateTimedFloatKeys = Template{
    Name: "TimedFloatKeys",
    UUID: MustHexToUUID("F406B180-7B3B-11cf-8F52-0040333594A3"),
    Mode: 'c', // closed
    Members: []TemplateMember{
        {
            Name:       "time",
            Type:       "DWORD",
        },
        {
            Name:       "tfkeys",
            Type:       "FloatKeys",
        },
    },
}

var TemplateVector = Template{
    Name: "Vector",
    UUID: MustHexToUUID("3D82AB5E-62DA-11cf-AB39-0020AF71E433"),
    Mode: 'c', // closed
    Members: []TemplateMember{
        {
            Name:       "x",
            Type:       "float",
        },
        {
            Name:       "y",
            Type:       "float",
        },
        {
            Name:       "z",
            Type:       "float",
        },
    },
}

var defaultTemplates = []*Template{
    &TemplateAnimation,
    &TemplateAnimationKey,
    &TemplateAnimationSet,
    &TemplateCoords2D,
    &TemplateFloatKeys,
    &TemplateFrame,
    &TemplateFrameTransformMatrix,
    &TemplateMatrix4x4,
    &TemplateMesh,
    &TemplateMeshFace,
    &TemplateMeshNormals,
    &TemplateMeshTextureCoords,
    &TemplateTimedFloatKeys,
    &TemplateVector,
}
