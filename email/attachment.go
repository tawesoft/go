package email

import (
    "bytes"
    "encoding/base64"
    "encoding/json"
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

    // Reader is a lazy reader. e.g. a function that returns the result of os.Open.
    Reader func() (io.ReadCloser, error)
}

type jsonAttachment struct {
    Filename string
    Mimetype string
    Content []byte
}

// Implements the json.Marshal interface. Note that the JSON content is ALWAYS
// Base64 encoded (with whitespace).
func (a *Attachment) MarshalJSON() ([]byte, error) {
    r, err := a.Reader()
    if err != nil {
        return nil, fmt.Errorf("attachment open error: %v", err)
    }
    defer r.Close()

    w := &bytes.Buffer{}
    encoder := base64.NewEncoder(base64.StdEncoding, lineBreaker{writer: w})
    defer encoder.Close()

    _, err = io.Copy(encoder, r)
    if err != nil {
        return nil, fmt.Errorf("attachment read error: %v", err)
    }

    return json.Marshal(jsonAttachment{
        Filename: a.Filename,
        Mimetype: a.Mimetype,
        Content:  w.Bytes(),
    })
}

// Implements the json.Unarshal interface. Note that the JSON content is ALWAYS
// Base64 encoded (with whitespace).
func (a *Attachment) UnmarshalJSON(data []byte) error {
    var j jsonAttachment

    err := json.Unmarshal(data, &j)
    if err != nil { return err }

    a.Filename = j.Filename
    a.Mimetype = j.Mimetype
    a.Reader = func() (io.ReadCloser, error) {
        r := bytes.NewReader(j.Content)
        decoder := base64.NewDecoder(base64.StdEncoding, r)
        return io.NopCloser(decoder), nil
    }

    return nil
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

    _, err = io.Copy(encoder, reader)
    if err != nil { return err }

    fmt.Fprintf(w, "\r\n")
    return nil
}
