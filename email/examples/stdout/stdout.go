// Format an email message and print it to a Writer (here, stdout).
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

