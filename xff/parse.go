package xff

// xff/parse.go parses individual tokens of a DirectX .x file

import (
    "bufio"
    "fmt"
    "io"
)

type lexeme struct {
    token [256]byte
    index int
}

func (l *lexeme) Append(c byte) bool {
    if l.index == 255 { return false }
    l.token[l.index] = c
    l.index++
    return true
}

func (l *lexeme) String() string {
    return string(l.token[:l.index])
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
                        if !l.Append(i) { err = fmt.Errorf("token too long"); goto done }
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
                        if !l.Append(i) { err = fmt.Errorf("token too long"); goto done }
                }
        }
    }
    
    if state == stateCapture { err = io.EOF }
    
    done:
        return word, err
}
