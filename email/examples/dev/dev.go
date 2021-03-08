// Example using the gmime parser to check folding whitespace
package main

import (
    "fmt"
    "strings"

    // sudo apt-get install libgmime-2.6-dev
    // go get github.com/sendgrid/go-gmime
    "github.com/sendgrid/go-gmime/gmime"
)

func must(s string, b bool) string {
    if !b { panic("parse error") }
    return s
}

func main() {
    message := `From: "Alan Turing" <turing.alan@example.org>
To: "Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace
 Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper"
 <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper"
 <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>,"Grace Hopper" <amazing.grace@example.net>
Cc: 
Bcc: "BCC1" <bcc1@example.net>,"BCC2" <bbc2@example.net>
Date: Mon, 08 Mar 2021 04:13:22 +0000 (GMT)
Subject: =?utf-8?q?Computer_Science_is_Cool!_=E2=9D=A4?=
Mime-Version: 1.0
Message-Id: <20210308041322.bf8d41db8300d1cf18c43ac92de6e43c@localhost>
X-Category: newsletter
X-Category: marketing
Content-Type: multipart/mixed; boundary="MXD-'q=Pp+,TTOWm8E,Xd=q-PmSDkoriH4EThcjeci5kiilj'w8E=AGvkIT4uVb9p3wlp6"

--MXD-'q=Pp+,TTOWm8E,Xd=q-PmSDkoriH4EThcjeci5kiilj'w8E=AGvkIT4uVb9p3wlp6
Content-Type: multipart/related; boundary="REL-KaH4a5tvo:i1Dv=?uZxj,6w4dAqHsIahO7MmiNxOgV',QuJr+pp:7533p762PwHB3N"

--REL-KaH4a5tvo:i1Dv=?uZxj,6w4dAqHsIahO7MmiNxOgV',QuJr+pp:7533p762PwHB3N
Content-Type: multipart/alternative; boundary="ALT-iYa/2Sr0iWbedZwuqCTDeVD_?JnepeLdnnvvdn4ZwNzpwEzoHVNeGcyihy6HiebH2j"

--ALT-iYa/2Sr0iWbedZwuqCTDeVD_?JnepeLdnnvvdn4ZwNzpwEzoHVNeGcyihy6HiebH2j
Content-Type: text/plain; charset=utf-8
Content-Transfer-Encoding: quoted-printable

This is a test email!
--ALT-iYa/2Sr0iWbedZwuqCTDeVD_?JnepeLdnnvvdn4ZwNzpwEzoHVNeGcyihy6HiebH2j
Content-Type: text/html; charset=utf-8
Content-Transfer-Encoding: quoted-printable

<!DOCTYPE html><html lang=3D"en"><body><p>This is a test email!</p></body><=
/html>
--ALT-iYa/2Sr0iWbedZwuqCTDeVD_?JnepeLdnnvvdn4ZwNzpwEzoHVNeGcyihy6HiebH2j--
--REL-KaH4a5tvo:i1Dv=?uZxj,6w4dAqHsIahO7MmiNxOgV',QuJr+pp:7533p762PwHB3N--
--MXD-'q=Pp+,TTOWm8E,Xd=q-PmSDkoriH4EThcjeci5kiilj'w8E=AGvkIT4uVb9p3wlp6
Content-Disposition: attachment; filename="attachment1.txt"; filename*="attachment1.txt"
Content-Type: text/plain; charset=utf-8
Content-Transfer-Encoding: base64

SGVsbG8gd29ybGQh
Cg==--MXD-'q=Pp+,TTOWm8E,Xd=q-PmSDkoriH4EThcjeci5kiilj'w8E=AGvkIT4uVb9p3wlp6--
.`

reader := strings.NewReader(message)
parse := gmime.NewParse(reader)

    fmt.Printf(`
    From: %s
    To: %s
    Subject: %s
`, must(parse.From()), parse.To(), must(parse.Subject()))

}


