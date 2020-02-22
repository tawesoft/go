package xff

// xff/parse.go parses individual tokens of a DirectX .x file

import (
    "bufio"
    "fmt"
    "io"
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

func mustReadExactSymbol(r io.ByteReader, expected byte, desc string) {
    var symbol, err = readSymbol(r)
    if err != nil { panic(err) }
    if symbol != expected {
        panic(fmt.Errorf("expected '%c' for %s but got '%c'", expected, desc, symbol))
    }
}

func mustReadSymbol(r io.ByteReader) byte {
    var symbol, err = readSymbol(r)
    if err != nil { panic(err) }
    return symbol
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

func mustReadString(r *bufio.Reader) string {
    var result, err = readString(r)
    if err != nil { panic(err) }
    return result
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

func mustReadAtom(r *bufio.Reader) string {
    var word, err = readAtom(r)
    if err != nil { panic(err) }
    return word
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
