package xff

// xff/parse.go parses individual tokens of a DirectX .x file

import (
    "bufio"
    "fmt"
    "io"
    "strconv"
)

type lexeme struct {
    token []byte
    index int
}

func (l *lexeme) Append(c byte) {
    l.token = append(l.token, c)
}

func (l *lexeme) String() string {
    return string(l.token)
}

// mustReadExpectedSymbol panics unless the expected symbol is consumed or on I/O error
func mustReadExpectedSymbol(r io.ByteReader, expected byte, desc string) {
    var symbol, err = readSymbol(r)
    if err == io.EOF { panic(fmt.Errorf("unexpected EOF")) }
    if err != nil { panic(err) }
    if symbol != expected {
        panic(fmt.Errorf("expected '%c' for %s but got '%c'", expected, desc, symbol))
    }
}

// mustReadSymbol reads a symbol but panics on I/O error
func mustReadSymbol(r io.ByteReader) byte {
    var symbol, err = readSymbol(r)
    if err == io.EOF { panic(fmt.Errorf("unexpected EOF")) }
    if err != nil { panic(err) }
    return symbol
}

// mustAcceptSymbol returns true if the expected symbol is consumed, or false and the reader is not advanced.
// Note "must" here means "must succeed without I/O error", not "must accept the symbol".
func mustAcceptSymbol(r *bufio.Reader, expected byte) bool {
    var symbol, err = readSymbol(r)
    if err == io.EOF { panic(fmt.Errorf("unexpected EOF")) }
    if err != nil { panic(err) }
    if symbol == expected { return true }
    mustUnreadByte(r)
    return false
}

// mustPeekSymbol looks at the next symbol and unwinds. It panics on I/O error.
func mustPeekSymbol(r *bufio.Reader) byte {
    var symbol, err = readSymbol(r)
    if err == io.EOF { panic(fmt.Errorf("unexpected EOF")) }
    if err != nil { panic(err) }
    mustUnreadByte(r)
    return symbol
}

// mustReadString panics on I/O error
func mustReadString(r *bufio.Reader) string {
    var result, err = readString(r)
    if err == io.EOF { panic(fmt.Errorf("unexpected EOF")) }
    if err != nil { panic(err) }
    return result
}

// mustReadAtom panics on I/O error
func mustReadAtom(r *bufio.Reader) string {
    var word, err = readAtom(r)
    if err == io.EOF { panic(fmt.Errorf("unexpected EOF")) }
    if err != nil { panic(err) }
    return word
}

// mustReadInt panics unless an atom can be read and parsed as an integer
func mustReadInt(r *bufio.Reader, bits int) int64 {
     value, err := strconv.ParseInt(mustReadAtom(r), 0, bits)
     if err != nil { panic(fmt.Errorf("int decode error: %v", err)) }
     return value
}

// mustReadUint panics unless an atom can be read and parsed as an unsigned integer
func mustReadUint(r *bufio.Reader, bits int) uint64 {
     value, err := strconv.ParseUint(mustReadAtom(r), 0, bits)
     if err != nil { panic(fmt.Errorf("int decode error: %v", err)) }
     return value
}

// mustReadFloat panics unless an atom can be read and parsed as an float
func mustReadFloat(r *bufio.Reader, bits int) float64 {
     value, err := strconv.ParseFloat(mustReadAtom(r), bits)
     if err != nil { panic(fmt.Errorf("float decode error: %v", err)) }
     return value
}

// mustUnreadByte should always succeed, because we're never rewinding more than a single byte. If it panics, that's
// a bug with our program not a parse error
func mustUnreadByte(r *bufio.Reader) {
    var err = r.UnreadByte()
    if err != nil { panic(fmt.Errorf("I/O rewind error: %s", err)) }
}

// readSymbol reads a single symbol, ignoring whitespace and comments, and returns it as a byte.
func readSymbol(r io.ByteReader) (symbol byte, err error) {
    var i byte
    var state int
    
    const (
        stateStart   int = 0
        stateComment int = 1
    )
    
    for idx := 0; err == nil; idx++ {
        i, err = r.ReadByte()
        
        switch state {
            case stateStart:
                switch i {
                    case  '/': state = stateComment
                    case  '#': state = stateComment
                    case  ' ': // pass
                    case '\t': // pass
                    case '\r': // pass
                    case '\n': // pass
                    default:
                        return i, nil
                }
            case stateComment:
                switch i {
                    case '\n': state = stateStart
                }
        }
    }
    
    return 0, io.EOF
}

// readString reads a string excluding the enclosing quotes
func readString(r *bufio.Reader) (result string, err error) {
    var i byte
    var l lexeme
    var escape bool = false
    
    for idx := 0; err == nil; idx++ {
        i, err = r.ReadByte()
        
        if (i == '\\') && !escape {
            escape = true
        } else if (i == '"') && !escape {
            r.UnreadByte()
            result = l.String()
            goto done
        } else {
            escape = false
            l.Append(i)
        }
    }
    
    err = io.EOF
    
    done:
        return result, err
}

// readAtom reads a single token delimited by comma, semicolon, whitespace, braces, brackets, quotation marks, and
// comments, and returns it as a string.
func readAtom(r *bufio.Reader) (word string, err error) {
    var i byte
    var state int
    var l lexeme
    
    const (
        stateStart   int = 0
        stateComment int = 1
        stateCapture int = 2
    )
    
    for idx := 0; err == nil; idx++ {
        i, err = r.ReadByte()
        
        switch state {
            case stateStart:
                switch i {
                    case  '/': state = stateComment
                    case  '#': state = stateComment
                    case  ' ': // pass
                    case '\t': // pass
                    case '\r': // pass
                    case '\n': // pass
                    default:
                        state = stateCapture
                        l.Append(i)
                }
            case stateComment:
                switch i {
                    case '\n': state = stateStart
                }
            case stateCapture:
                switch i {
                    case '\n': fallthrough
                    case '\r': fallthrough
                    case '\t': fallthrough
                    case  ' ': fallthrough
                    case  ',': fallthrough
                    case  ';': fallthrough
                    case  '#': fallthrough
                    case  '/': fallthrough
                    case  '{': fallthrough
                    case  '}': fallthrough
                    case  '<': fallthrough
                    case  '>': fallthrough
                    case  '"': fallthrough
                    case  '[': fallthrough
                    case  ']':
                        r.UnreadByte()
                        word = l.String()
                        goto done
                    default:
                        l.Append(i)
                }
        }
    }
    
    if state == stateCapture { err = io.EOF }
    
    done:
        return word, err
}
