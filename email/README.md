# email - format multipart RFC 2045 email

```shell script
go get "tawesoft.co.uk/go/"
```

```go
import "tawesoft.co.uk/go/email"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_email] ∙ [docs][docs_email] ∙ [src][src_email] | [MIT][copy_email] | candidate |

[home_email]: https://tawesoft.co.uk/go/email
[src_email]:  https://github.com/tawesoft/go/tree/master/email
[docs_email]: https://godoc.org/tawesoft.co.uk/go/email
[copy_email]: https://github.com/tawesoft/go/tree/master/email/LICENSE.txt

## About

Package email implements the formatting of multipart RFC 2045 e-mail messages,
including headers, attachments, HTML email, and plain text.

File attachments are lazy, and read from disk only at the time the e-mail is
sent.


## Examples


Format an email message and print it to a Writer (here, stdout).
```go
package main

import (
    "net/mail"
    "os"

    "tawesoft.co.uk/go/email"
)

func main() {
    var eml = email.Message{
        ID:    email.NewMessageID("localhost"),
        From:  mail.Address{"Alan Turing", "turing.alan@example.org"},
        To:  []mail.Address{{"Grace Hopper", "amazing.grace@example.net"},},
        Bcc: []mail.Address{{"BCC1", "bcc1@example.net"}, {"BCC2", "bbc2@example.net"}},
        Subject: "Computer Science is Cool! ❤",
        Text: `This is a test email!`,
        Html: `<!DOCTYPE html><html lang="en"><body><p>This is a test email!</p></body></html>`,
        Attachments: []*email.Attachment{
            email.FileAttachment("attachment1.txt"),
        },
        Headers: mail.Header{
            "X-Category": []string{"newsletter", "marketing"},
        },
    }

    var err = eml.Write(os.Stdout)
    if err != nil { panic(err) }
}
```

## Changes

### 2021-03-07

* Breaking changes to this email package, as previously warned, bump the
monorepo tagged version to v0.2.0 and upgrade the email package stability
rating from "unstable" to "candidate". For previous behavior, point your
imports to `tawesoft.co.uk/go/legacy/email`.

* Attachments are now read/written more efficiently.

* Attachments are now closed properly!

* Attachment Reader method is now required to return something satisfying
the io.ReadCloser interface. If no Close is required, wrap the return
value in an `io.NopCloser`.

* The Envelope struct no longer has a message field - instead, use
an (Envelope, Message) 2-tuple where you need both of these items.

* An email's Message-ID header is no longer implicitly generated for an
email. This is left to the mail submission agent.

* If you ARE implementing a mail submission agent, an email's Message-ID
header can be specified by the new ID field on the Message struct type.

* A cryptographically unique Message ID can be generated from the newly
exposed function, NewMessageID.

* The Print method on Message is renamed Write.

* Email message lines longer than 998 characters are now supported in
headers using folding white space. Note that some parsers, such as Go's
`net.mail`, do not understand this syntax (even though it is allowed).

* The new method WriteCompat on Message won't use folding white space to
support long headers and will instead generate an error. Use this method
in preference to Write if you are expecting the consumer of your email
message (e.g. a Go implementation) will be unable to handle folding white
space.



## Getting Help

This package is part of [tawesoft.co.uk/go](https://www.tawesoft.co.uk/go),
a monorepo for small Go modules maintained by Tawesoft®.
Check out that URL for more information about other Go modules from
Tawesoft plus community and commercial support options.