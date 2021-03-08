package email

import (
    "encoding/base64"
    "fmt"
    "io"
    "mime"
    "os"
    "path"
)

// Attachment defines an e-mail attachment. They are read lazily.
type Attachment struct {
    // Filename is the name to give the attachment in the email
    Filename string

    // Mimetype e.g. "application/pdf".
    // If an empty string, then attempts to automatically detect based on filename extension.
    Mimetype string

    // Reader is a lazy reader. e.g. return the result of os.Open.
    Reader func() (io.ReadCloser, error)
}

// FileAttachment returns an Attachment from a file path. The file at that path is lazily opened at the time the
// attachment is sent.
func FileAttachment(src string) *Attachment {

    var mimetype string
    var base = path.Base(src)

    var reader = func() (io.ReadCloser, error) {
        return os.Open(src)
    }

    mimetype =  mime.TypeByExtension(path.Ext(src))
    if mimetype == "" {
        mimetype = "application/octet-stream"
    }

    return &Attachment{
        Filename: base,
        Mimetype: mimetype,
        Reader: reader,
    }
}

// write encodes an attachment as part of a RFC 2045 MIME Email
func (a *Attachment) write(w io.Writer) error {
    fmt.Fprintf(w, "Content-Disposition: attachment; filename=\"%[1]s\"; filename*=\"%[1]s\"\r\n",
        optionalQEncode(a.Filename))
    fmt.Fprintf(w, "Content-Type: %s\r\n", a.Mimetype)
    fmt.Fprintf(w, "Content-Transfer-Encoding: base64\r\n")
    fmt.Fprintf(w, "\r\n")

    var reader, err = a.Reader()
    if err != nil {
        return fmt.Errorf("attachment open error: %v", err)
    }
    defer reader.Close()

    var encoder = base64.NewEncoder(base64.StdEncoding, lineBreaker{writer: w})
    defer encoder.Close()

    io.Copy(encoder, reader)

    fmt.Fprintf(w, "\r\n")
    return nil
}
