# email - format multipart RFC 2045 email

## About

Package email implements the formatting of multipart RFC 2045 e-mail messages,
including headers, attachments, HTML email, and plain text.

|  Links  | License | Stable? | 
|:-------:|:-------:|:-------:| 
| [home][home_] ∙ [docs][docs_] ∙ [src][src_] | [MIT][copy_] | ✔ yes |

[home_]: https://tawesoft.co.uk/go/email
[src_]:  https://github.com/tawesoft/go/tree/master/email
[docs_]: https://godoc.org/tawesoft.co.uk/go/email
[copy_]: https://github.com/tawesoft/go/tree/master/email/_COPYING.md

## Download

```shell script
go get -u tawesoft.co.uk/go
```

## Import

```
import tawesoft.co.uk/go/email
```

## Example:

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
            //email.FileAttachment("Entscheidungsproblem.pdf"),
            //email.FileAttachment("funny-cat-meme.png"),
        },
        Headers: mail.Header{
            "X-Category": []string{"newsletter", "marketing"},
        },
    }
    
    var err = eml.Print(os.Stdout)
    if err != nil { panic(err) }
}
```