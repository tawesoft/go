package email

import (
    "bufio"
    "crypto/rand"
    "fmt"
    "io"
    "mime/quotedprintable"
    "net/mail"
    "net/textproto"
    "strings"
    "time"
)

// Messages defines a multipart e-mail (including headers, HTML body, plain text body, and attachments).
//
// Use the `headers` parameter to specify additional headers. Note that `mail.Header` maps keys to a *list* of
// strings, because some headers may appear multiple times.
type Message struct {
    From mail.Address
    To   []mail.Address
    Cc   []mail.Address
    Bcc  []mail.Address
    Subject string
    Headers mail.Header
    Html string
    Text string
    Attachments []*Attachment
}

// Envelope wraps an Email with some SMTP protocol information for extra control.
type Envelope struct {
    // From is the sender. Usually this should match the Email From address. In the cause of autoreplies (like "Out of
    // Office" or bounces or delivery status notifications) this should be an empty string to stop an infinite loop
    // of bounces)
    From string
    
    // Data is just a pointer to an Email struct
    Data *Message
    
    // ReceiptTo is normally automatically generated from the Email To/CC/BCC addresses
    ReceiptTo []string
}

// bound defines an email message boundary
type bound struct {
    label string
}

func (b bound) start(w io.Writer, contentType string) {
    fmt.Fprintf(w, "Content-Type: %s; boundary=\"%s\"\r\n", contentType, b.label)
    fmt.Fprintf(w, "\r\n")
    fmt.Fprintf(w, "--%s\r\n", b.label)
}

func (b bound) next(w io.Writer) {
    fmt.Fprintf(w, "--%s\r\n", b.label)
}

func (b bound) end(w io.Writer) {
    fmt.Fprintf(w, "--%s--\r\n", b.label)
}

// boundary returns a randomly-generated RFC 1341 boundary with a label of exactly 70 characters with a given prefix.
// The randomly generated prefix is cryptographically secure iff `rand` is `crypto/rand`.
func boundary(prefix string) (bound, error) {
    
    // RFC 1341 sets this maximum length for a boundary
    const maxlen int = 70
    
    // NOTE: RFC 1341 says we can use space in a boundary as long as it isn't a trailing space, but for simplicity
    // of implementation we avoid this.
    const bchars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ'()+_,-./:=?"
    
    // This is a purely arbitrary limit but if the prefix is too long then we run the risk of collisions between the
    // boundary and a message.
    if (len(prefix) > 16) { return bound{}, fmt.Errorf("boundary prefix too long") }
    
    var randlen = maxlen - len(prefix)
    xs := make([]byte, randlen)
    
    _, err := rand.Read(xs)
    if err != nil { return bound{}, fmt.Errorf("boundary random source read error: %v", err) }
    
    // for bytes in range [0-255] cast down to bcharsnospace
    for i, x := range xs {
        xs[i] = bchars[int(x) % len(bchars)]
    }
    
    return bound{prefix + string(xs)}, nil
}

// Print writes a multipart Email to dest.
//
// NOTE: the maximum length of a email message line is 998 characters. If sending emails to multiple addresses
// the caller should keep this limit in mind and divide the addresses over multiple calls to this function.
func (e *Message) Print(dest io.Writer) error {

    // RFC 5332 date format with (comment) for time.Format
    const RFC5332C = "Mon, 02 Jan 2006 15:04:05 -0700 (MST)"
    
    var err error
    var qp *quotedprintable.Writer
    var bufferedDest = bufio.NewWriter(dest)
    var dw = textproto.NewWriter(bufferedDest).DotWriter()
    defer dw.Close()
    
    // format a list of mail.Address objects as a comma-separated string
    var addresses = func (xs []mail.Address) string {
        var s = make([]string, 0, len(xs))
        
        for _, x := range xs {
            s = append(s, x.String())
        }
        
        return strings.Join(s, ",")
    }
    
    var coreHeaders = []struct{left string; right string} {
        {"From",         e.From.String()},
        {"To",           addresses(e.To)},
        {"Cc",           addresses(e.Cc)},
        {"Bcc",          addresses(e.Bcc)},
        {"Date",         time.Now().Format(RFC5332C)},
        {"Subject",      optionalQEncode(e.Subject)},
        {"MIME-Version", "1.0"},
        {"Message-ID",   msgid(e.From)},
    }
    
    for _, v := range coreHeaders {
        fmt.Fprintf(dw, "%s: %s\r\n", textproto.CanonicalMIMEHeaderKey(v.left), v.right)
    }
    
    for k, vs := range e.Headers {
        for _, v := range vs {
            fmt.Fprintf(dw, "%s: %s\r\n", textproto.CanonicalMIMEHeaderKey(k), v)
        }
    }
    
    bndMxd, err := boundary("MXD-")
    if err != nil { return err }
    bndMxd.start(dw, "multipart/mixed")
    
    bndRel, err := boundary("REL-")
    if err != nil { return err }
    bndRel.start(dw, "multipart/related")
    
    bndAlt, err := boundary("ALT-")
    if err != nil { return err }
    bndAlt.start(dw, "multipart/alternative")
    
    fmt.Fprintf(dw, "Content-Type: text/plain; charset=utf-8\r\n")
    fmt.Fprintf(dw, "Content-Transfer-Encoding: quoted-printable\r\n")
    fmt.Fprintf(dw, "\r\n")
    qp = quotedprintable.NewWriter(dw)
    io.WriteString(qp, strings.TrimSpace(e.Text))
    qp.Close()
    fmt.Fprintf(dw, "\r\n")
    
    bndAlt.next(dw)
    
    fmt.Fprintf(dw, "Content-Type: text/html; charset=utf-8\r\n")
    fmt.Fprintf(dw, "Content-Transfer-Encoding: quoted-printable\r\n")
    fmt.Fprintf(dw, "\r\n")
    qp = quotedprintable.NewWriter(dw)
    io.WriteString(qp, e.Html)
    qp.Close()
    fmt.Fprintf(dw, "\r\n")
    
    bndAlt.end(dw)
    bndRel.end(dw)
    
    for _, attachment := range e.Attachments {
        bndMxd.next(dw)
        err = attachment.write(dw)
        if err != nil {
            return fmt.Errorf("error writing attachment %s: %v", attachment.Filename, err)
        }
    }
    
    bndMxd.end(dw)
    
    return nil
}
