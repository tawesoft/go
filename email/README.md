# email - format multipart RFC 2045 email

```shell script
go get "tawesoft.co.uk/go/"
```

```go
import "tawesoft.co.uk/go/email"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_email] ∙ [docs][docs_email] ∙ [src][src_email] | [MIT][copy_email] | ✘ **no** |

[home_email]: https://tawesoft.co.uk/go/email
[src_email]:  https://github.com/tawesoft/go/tree/master/email
[docs_email]: https://godoc.org/tawesoft.co.uk/go/email
[copy_email]: https://github.com/tawesoft/go/tree/master/email/LICENSE.txt

## About

Package email implements the formatting of multipart RFC 2045 e-mail messages,
including headers, attachments, HTML email, and plain text.

File attachments are lazy, and read from disk only at the time the e-mail is
sent.


## Package Stability


It is likely that this package will change at some point as follows:

* A Message-ID header will no longer be implicitly generated for a Message.

* The Envelope struct will no longer have a message field - instead, use
an (Envelope, Message) 2-tuple where you need both of these items.

This is a breaking change. As such, when this happens, the old behaviour will
be made available at tawesoft.co.uk/go/legacy/email.


## Example


This example demonstrates formatting an email message and printing it to a
Writer (here, `os.Stdout`).


```go
package main

import (
    "net/mail"
    "os"

    "tawesoft.co.uk/go/email"
)

func main() {
    var eml = email.Message{
        From:  mail.Address{"Alan Turing", "turing.alan@example.org"},
        To:  []mail.Address{{"Grace Hopper", "amazing.grace@example.net"}},
        Bcc: []mail.Address{{"BCC1", "bcc1@example.net"}, {"BCC2", "bbc2@example.net"}},
        Subject: "Computer Science is Cool! ❤",
        Text: `This is a test email!`,
        Html: `<!DOCTYPE html><html lang="en"><body><p>This is a test email!</p></body></html>`,
        Attachments: []*email.Attachment{
            email.FileAttachment("Entscheidungsproblem.pdf"),
            email.FileAttachment("funny-cat-meme.png"),
        },
        Headers: mail.Header{
            "X-Category": []string{"newsletter", "marketing"},
        },
    }

    var err = eml.Print(os.Stdout)
    if err != nil { panic(err) }
}
```

## Getting Help

This package is part of [tawesoft.co.uk/go](https://www.tawesoft.co.uk/go),
a monorepo for small Go modules maintained by Tawesoft®.
Check out that URL for more information about other Go modules from
Tawesoft plus community and commercial support options.