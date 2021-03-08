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

// Message defines a multipart e-mail (including headers, HTML body, plain text body, and attachments).
//
// Use the `headers` parameter to specify additional headers. Note that `mail.Header` maps keys to a *list* of
// strings, because some headers may appear multiple times.
type Message struct {
    ID  string // Message-ID header, excluding angle brackets
    From mail.Address
    To   []mail.Address
    Cc   []mail.Address
    Bcc  []mail.Address
    Subject string

    // Headers are additional headers for the message. The combination of a
    // header and a value as strings must not exceed a length of 996
    // characters. Longer values CANNOT be supported with folding white space
    // syntax without advance knowledge (special cases are possible but not
    // currently implemented).
    Headers mail.Header

    // Html is a HTML-encoded version of the message. Lines must not exceed
    // 998 characters.
    Html string

    // Text is a plain-text version of the message. It is your responsibility
    // to ensure word-wrapping. Lines must not exceed 998 characters.
    Text string

    // Attachments is a lazily-loaded sequence of attachments. May be nil.
    Attachments []*Attachment
}

// Envelope wraps an Email with some SMTP protocol information for extra control.
type Envelope struct {
    // From is the sender. Usually this should match the Email From address. In
    // the cause of autoreplies (like "Out of Office" or bounces or delivery
    // status notifications) this should be an empty string to stop an infinite
    // loop of bounces)
    From string

    // ReceiptTo is a list of recipients. This is normally automatically
    // generated from the Email To/CC/BCC addresses.
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

// Write writes a multipart Email to dest.
func (e *Message) Write(dest io.Writer) error {
    return e.write(dest, false)
}

// WriteCompat is like Write, however very long headers (caused, for example,
// by sending a message with many To addresses, or a subject that is too long)
// are not encoded using folding white space and instead cause an error.
func (e *Message) WriteCompat(dest io.Writer) error {
    return e.write(dest, true)
}

func (e *Message) write(dest io.Writer, compat bool) (ret error) {
    defer func() {
        if r := recover(); r != nil {
            if err, ok := r.(error); ok {
                ret = err
            } else {
                ret = fmt.Errorf("panic: %v", r)
            }
        }
    }()

    // folding white space support
    var fws func(value string, keylen int) string
    if compat {
        fws = func(value string, keylen int) string {
            result, err := fwsNone(value, keylen, 998)
            if err != nil { panic(err) }
            return result
        }
    } else {
        fws = func(value string, keylen int) string {
            result, err := fwsWrap(value, keylen, 998)
            if err != nil { panic(err) }
            return result
        }
    }

    // RFC 5332 date format with (comment) for time.Format
    const RFC5332C = "Mon, 02 Jan 2006 15:04:05 -0700 (MST)"

    //var err error
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
        {"From",         fws(e.From.String(), len("From"))},
        {"To",           fws(addresses(e.To), len("To"))},
        {"Cc",           fws(addresses(e.Cc), len("Cc"))},
        {"Bcc",          fws(addresses(e.Bcc), len("Bcc"))},
        {"Date",         time.Now().Format(RFC5332C)},
        {"Subject",      fws(optionalQEncode(e.Subject), len("Subject"))},
        {"MIME-Version", "1.0"},
    }

    if e.ID != "" {
        coreHeaders = append(coreHeaders,
            struct{left string; right string}{"Message-ID", "<"+e.ID+">"})
    }

    for _, v := range coreHeaders {
        fmt.Fprintf(dw, "%s: %s\r\n", textproto.CanonicalMIMEHeaderKey(v.left), v.right)
    }

    for k, vs := range e.Headers {
        for _, v := range vs {
            ck := textproto.CanonicalMIMEHeaderKey(k)
            fwsValue, err := fwsNone(v, len(ck), 998)
            if err != nil { return err }
            fmt.Fprintf(dw, "%s: %s\r\n", ck, fwsValue)
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
