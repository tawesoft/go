package parse

const EOF = int(-1)

// Interface ByteStreamL1 defines a byte stream with a single byte of
// lookahead (a L(1) parser).
type ByteStreamL1 interface {
    // Next returns the current and the next uint8 byte in the string, or EOF.
    Next() (int, int)
}

// nullByteStreamL1 implements a ByteStreamL1 where the `Next`
// method always returns `parse.EOF`.
type nullByteStreamL1 struct{}

func (n *nullByteStreamL1) Next() (int, int) {
    return EOF, EOF
}

// stringByteStreamL1 implements a ByteStreamL1 over a string
type stringByteStreamL1 struct {
    str     string
    current int
    next    int
    offset  int
}

// Returns a ByteStreamL1 over a string
func StringByteStreamL1(s string) ByteStreamL1 {
    if len(s) == 0 {
        return &nullByteStreamL1{}
    }
    
    return &stringByteStreamL1{
        str:     s,
        current: EOF,
        next:    int(s[0]),
    }
}

func (s *stringByteStreamL1) Next() (int, int) {
    s.current = s.next
    if len(s.str) < s.offset {
        s.offset++
        s.next = int(s.str[s.offset])
    } else {
        s.next = EOF
    }
    
    return s.current, s.next
}
