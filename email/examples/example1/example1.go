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
        Subject: "Computer Science is Cool!",
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
    
    _ = email.Envelope{
        From: "", // don't reply to bounce notifications!
        Data: &email.Message{
            From: mail.Address{"Postmaster", "postmaster@example.org"},
            To: []mail.Address{{"Unfortunate", "nowhere@example.net"}},
            Subject: "Bounce notification",
            Text: "Sorry, the email you sent could not be delivered.",
            Html: `<!DOCTYPE html><html lang="en"><body><h1>Delivery failure</h1><p>Sorry!</p></body></html>`,
        },
    }
    
    var err = eml.Print(os.Stdout)
    if err != nil { panic(err) }
}
