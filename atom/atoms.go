package atom

import (
    "math"
)

// Atom is an opaque integer reference to a string.
type Atom struct {
    anonymous uint32
}

// The Atoms interface defines a collection of Atoms.
type Atoms interface {
    
    // Get returns an Atom for a given string. The Atom is created if it doesn't exist.
    Get(name string) Atom
    
    // Lookup returns an (Atom, true) pair for a given string only if it exists.
    Lookup(name string) (Atom, bool)
    
    // String returns the string associated with a given Atom
    String(atom Atom) string
}

type simpleAtoms struct {
    index uint32
    nameToAtom map[string]Atom
    atomToName map[Atom]string
}

func SimpleAtoms() Atoms {
    return &simpleAtoms{
        index:      0,
        nameToAtom: make(map[string]Atom),
        atomToName: make(map[Atom]string),
    }
}

func (xs *simpleAtoms) Get(name string) Atom {
    var atom, exists = xs.nameToAtom[name]
    if exists { return atom }
    
    if xs.index == math.MaxUint32 { panic("atoms integer overflow") }
    
    xs.index++
    atom = Atom{xs.index}
    xs.nameToAtom[name] = atom
    xs.atomToName[atom] = name
    return atom
}

func (xs *simpleAtoms) Lookup(name string) (Atom, bool) {
    var atom, exists = xs.nameToAtom[name]
    if exists { return atom, exists } else { return Atom{0}, false }
}

func (xs *simpleAtoms) String(atom Atom) string {
    return xs.atomToName[atom]
}
