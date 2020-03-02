package atom

import (
    "testing"
)

func TestAtoms(t *testing.T) {
    var xs = SimpleAtoms()
    
    var a1 = xs.Get("Atom One")
    var b1 = xs.Get("Atom Two")
    
    var a2 = xs.Get("Atom One")
    var b2 = xs.Get("Atom Two")
    
    if a1 == b1 { t.Errorf("fail") }
    if a2 == b2 { t.Errorf("fail") }
    
    if a1 != a2 { t.Errorf("fail") }
    if b1 != b2 { t.Errorf("fail") }
    
    var atom Atom; var exists bool
    atom, exists = xs.Lookup("Atom One"); if atom != a1 { t.Errorf("fail") }; if !exists { t.Errorf("fail") }
    atom, exists = xs.Lookup("Atom Two"); if atom != b1 { t.Errorf("fail") }; if !exists { t.Errorf("fail") }
    _,    exists = xs.Lookup("Atom Three"); if exists { t.Errorf("fail") }
}
