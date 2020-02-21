package xff

type TemplateMember struct {
    Name string
    Type string
    Dimensions []string
}

func (m *TemplateMember) PrimitiveType() bool {
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
        default:        return false
    }
}

func (m *TemplateMember) Size() int {
    switch m.Type {
        case "WORD":   return 2
        case "DWORD":  return 4
        case "FLOAT":  return 0
        case "DOUBLE": return 0
        case "CHAR":   return 0
        case "UCHAR":  return 0
        case "BYTE":   return 0
        case "STRING": return 0 // index to strings
        case "float":  return 4
        default:       return 0 // index to children
    }
}

type Template struct {
    Name string
    UUID string
    Mode byte // 'o'pen, 'c'losed, 'r'estricted
    Members []TemplateMember
}

var TemplateAnimation = Template{
    Name: "Animation",
    UUID: "3D82AB4F-62DA-11cf-AB39-0020AF71E433",
    Mode: 'o', // open
}

var TemplateAnimationKey = Template{
    Name: "AnimationKey",
    UUID: "10DD46A8-775B-11CF-8F52-0040333594A3",
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

// TemplateAnimationSet Contains one or more Animation objects
var TemplateAnimationSet = Template{
    Name: "AnimationSet",
    UUID: "3D82AB50-62DA-11cf-AB39-0020AF71E433",
    Mode: 'r', // restricted to Animation objects (TODO)
}

var TemplateCoords2d = Template{
    Name: "Coords2D",
    UUID: "F6F23F44-7686-11cf-8F52-0040333594A3",
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
    UUID: "10DD46A9-775B-11cf-8F52-0040333594A3",
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
    UUID: "3D82AB46-62DA-11CF-AB39-0020AF71E433",
    Mode: 'o', // open
}

var TemplateFrameTransformMatrix = Template{
    Name: "FrameTransformMatrix",
    UUID: "F6F23F41-7686-11cf-8F52-0040333594A3",
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
    UUID: "F6F23F45-7686-11cf-8F52-0040333594A3",
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
    UUID: "3D82AB44-62DA-11CF-AB39-0020AF71E433",
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
    UUID: "3D82AB5F-62DA-11cf-AB39-0020AF71E433",
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
    UUID: "F6F23F43-7686-11cf-8F52-0040333594A3",
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
    UUID: "F6F23F40-7686-11cf-8F52-0040333594A3",
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
    UUID: "F406B180-7B3B-11cf-8F52-0040333594A3",
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
    UUID: "3D82AB5E-62DA-11cf-AB39-0020AF71E433",
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

var Templates = map[string]*Template{
    TemplateAnimation.Name:             &TemplateAnimation,
    TemplateAnimationKey.Name:          &TemplateAnimationKey,
    TemplateAnimationSet.Name:          &TemplateAnimationSet,
    TemplateCoords2d.Name:              &TemplateCoords2d,
    TemplateFloatKeys.Name:             &TemplateFloatKeys,
    TemplateFrame.Name:                 &TemplateFrame,
    TemplateFrameTransformMatrix.Name:  &TemplateFrameTransformMatrix,
    TemplateMatrix4x4.Name:             &TemplateMatrix4x4,
    TemplateMesh.Name:                  &TemplateMesh,
    TemplateMeshFace.Name:              &TemplateMeshFace,
    TemplateMeshNormals.Name:           &TemplateMeshNormals,
    TemplateMeshTextureCoords.Name:     &TemplateMeshTextureCoords,
    TemplateTimedFloatKeys.Name:        &TemplateTimedFloatKeys,
    TemplateVector.Name:                &TemplateVector,
}
