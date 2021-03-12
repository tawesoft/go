// Format an email message and print it, as well as its JSON serialisation, to
// a Writer (here, stdout).
package main

import (
    "encoding/json"
    "fmt"
    "net/mail"
    "os"
    "strings"

    "tawesoft.co.uk/go/email"
)

func main() {
    eml := email.Message{
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

    fmt.Printf("Formatted email:\n")
    err := eml.Write(os.Stdout)
    if err != nil { panic(err) }

    fmt.Printf("\n\nJSON serialisation:\n")
    encoder := json.NewEncoder(os.Stdout)
    encoder.SetIndent("", "    ")
    encoder.SetEscapeHTML(false)
    err = encoder.Encode(eml)
    if err != nil { panic(err) }

    fmt.Printf("\n\nJSON deserialisation:\n")
    var out email.Message
    decoder := json.NewDecoder(strings.NewReader(`{
    "ID": "20210312143531.183a5bf8f218c9c3e2dc4976a70676d2@localhost",
    "From": {
        "Name": "Alan Turing",
        "Address": "turing.alan@example.org"
    },
    "To": [
        {
            "Name": "Grace Hopper",
            "Address": "amazing.grace@example.net"
        }
    ],
    "Cc": null,
    "Bcc": [
        {
            "Name": "BCC1",
            "Address": "bcc1@example.net"
        },
        {
            "Name": "BCC2",
            "Address": "bbc2@example.net"
        }
    ],
    "Subject": "Computer Science is Cool! ❤",
    "Headers": {
        "X-Category": [
            "newsletter",
            "marketing"
        ]
    },
    "Html": "<!DOCTYPE html><html lang=\"en\"><body><p>This is a test email!</p></body></html>",
    "Text": "This is a test email!",
    "Attachments": [
        {
            "Filename": "attachment1.txt",
            "Mimetype": "text/plain; charset=utf-8",
            "Content": "U0dWc2JHOGdkMjl5YkdRaA=="
        }
    ]
}
`))
    decoder.DisallowUnknownFields()
    err = decoder.Decode(&out)
    if err != nil { panic(err) }
    out.Write(os.Stdout)
}

