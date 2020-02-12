// Package glcaps provides a nice interface to declare OpenGL capabilities you care about, including minimum required
// extensions or capabilities. Glcaps has no dependencies and is agnostic to the exact OpenGL binding used.
package glcaps

import (
    "sort"
)

// Extensions is an ordered list of supported OpenGL extensions.
type Extensions []string

// HaveExtension returns true iff the ordered list of supported OpenGL extensions contains a given extension.
func (extensions Extensions) Contains(key string) bool {
    var index = sort.SearchStrings(extensions, key)
    return (index < len(extensions)) && (extensions[index] == key)
}

// Binding implements a binding between this package and a specific OpenGL implementation (e.g. a specific `go-gl`
// module).
type Binding struct {
    GetIntegerv func(name uint32, data *int32)
    GetFloatv   func(name uint32, data *float32)
    GetStringi  func(name uint32, index uint32) string // required to return a Go string, not a C string!
}

// QueryExtensions returns all extensions supported by the current OpenGL context as a sorted list of strings. It is an
// error to call this method if a current OpenGL context does not exist.
func (b *Binding) QueryExtensions() Extensions {
    var numExtensions int32
    b.GetIntegerv(glconstants["GL_NUM_EXTENSIONS"], &numExtensions)
    if numExtensions <= 0 {
        panic("failed to query OpenGL extensions (is the OpenGL context current?)")
    }

    var xs = make([]string, 0, numExtensions)
    
    for i := uint32(0); i < uint32(numExtensions); i++ {
        x := b.GetStringi(glconstants["GL_EXTENSIONS"], i)
        if len(x) == 0 { continue }
        xs = append(xs, x)
    }

    sort.Strings(xs)
    return xs
}

// Error implements an error result type for reporting a capability that doesn't meet a requirement.
type Error struct {
    Field       string // the name of the field in the struct that failed
    Tag         string // the original tag string
    Requirement requirement // the requirement that failed
    Message     string // a human-readable message
}

type Errors []Error

func (es *Errors) append(e ... Error) {
    if *es == nil && (len(e) > 0) {
        *es = make([]Error, 0)
    }
    
    *es = append(*es, e...)
}




