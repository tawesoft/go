// tawesoft.co.uk/go/legacy/email
// 
// Copyright © 2020 Tawesoft Ltd <open-source@tawesoft.co.uk>
// Copyright © 2020 Ben Golightly <ben@tawesoft.co.uk>
// 
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction,  including without limitation the rights
// to use,  copy, modify,  merge,  publish, distribute, sublicense,  and/or sell
// copies  of  the  Software,  and  to  permit persons  to whom  the Software is
// furnished to do so, subject to the following conditions:
// 
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
// 
// THE SOFTWARE IS PROVIDED  "AS IS",  WITHOUT WARRANTY OF ANY KIND,  EXPRESS OR
// IMPLIED,  INCLUDING  BUT  NOT LIMITED TO THE WARRANTIES  OF  MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE  AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
// AUTHORS  OR COPYRIGHT HOLDERS  BE LIABLE  FOR ANY  CLAIM,  DAMAGES  OR  OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package email implements the formatting of multipart RFC 2045 e-mail messages,
// including headers, attachments, HTML email, and plain text.
// 
// File attachments are lazy, and read from disk only at the time the e-mail is
// sent.
// 
// Package Stability
// 
// It is likely that this package will change at some point as follows:
// 
// * A Message-ID header will no longer be implicitly generated for a Message.
// 
// * The Envelope struct will no longer have a message field - instead, use
// an (Envelope, Message) 2-tuple where you need both of these items.
// 
// This is a breaking change. As such, when this happens, the old behaviour will
// be made available at tawesoft.co.uk/go/legacy/email.
// 
// Example
// 
// This example demonstrates formatting an email message and printing it to a
// Writer (here, `os.Stdout`).
// 
//     package main
// 
//     import (
//         "net/mail"
//         "os"
// 
//         "tawesoft.co.uk/go/email"
//     )
// 
//     func main() {
//         var eml = email.Message{
//             From:  mail.Address{"Alan Turing", "turing.alan@example.org"},
//             To:  []mail.Address{{"Grace Hopper", "amazing.grace@example.net"}},
//             Bcc: []mail.Address{{"BCC1", "bcc1@example.net"}, {"BCC2", "bbc2@example.net"}},
//             Subject: "Computer Science is Cool! ❤",
//             Text: `This is a test email!`,
//             Html: `<!DOCTYPE html><html lang="en"><body><p>This is a test email!</p></body></html>`,
//             Attachments: []*email.Attachment{
//                 email.FileAttachment("Entscheidungsproblem.pdf"),
//                 email.FileAttachment("funny-cat-meme.png"),
//             },
//             Headers: mail.Header{
//                 "X-Category": []string{"newsletter", "marketing"},
//             },
//         }
// 
//         var err = eml.Print(os.Stdout)
//         if err != nil { panic(err) }
//     }
//
// Package Information
//
// License: MIT (see LICENSE.txt)
//
// Stable: no
//
// For more information, documentation, source code, examples, support, links,
// etc. please see https://www.tawesoft.co.uk/go and 
// https://www.tawesoft.co.uk/go/legacy/email
package email // import "tawesoft.co.uk/go/legacy/email"

// Code generated by internal. DO NOT EDIT.
// Instead, edit DESC.txt and run mkdocs.sh.