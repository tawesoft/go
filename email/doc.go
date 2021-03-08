// tawesoft.co.uk/go/email
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
// Examples
// 
// Format an email message and print it to a Writer (here, stdout).
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
//             ID:    email.NewMessageID("localhost"),
//             From:  mail.Address{"Alan Turing", "turing.alan@example.org"},
//             To:  []mail.Address{{"Grace Hopper", "amazing.grace@example.net"},},
//             Bcc: []mail.Address{{"BCC1", "bcc1@example.net"}, {"BCC2", "bbc2@example.net"}},
//             Subject: "Computer Science is Cool! ❤",
//             Text: `This is a test email!`,
//             Html: `<!DOCTYPE html><html lang="en"><body><p>This is a test email!</p></body></html>`,
//             Attachments: []*email.Attachment{
//                 email.FileAttachment("attachment1.txt"),
//             },
//             Headers: mail.Header{
//                 "X-Category": []string{"newsletter", "marketing"},
//             },
//         }
//     
//         var err = eml.Write(os.Stdout)
//         if err != nil { panic(err) }
//     }
//
//
// Package Information
//
// License: MIT (see LICENSE.txt)
//
// Stable: candidate
//
// For more information, documentation, source code, examples, support, links,
// etc. please see https://www.tawesoft.co.uk/go and 
// https://www.tawesoft.co.uk/go/email
//
//     2021-03-07
//     
//         * Breaking changes to this email package, as previously warned, bump the
//           monorepo tagged version to v0.2.0 and upgrade the email package stability
//           rating from "unstable" to "candidate". For previous behavior, point your
//           imports to `tawesoft.co.uk/go/legacy/email`.
//     
//         * Attachments are now read/written more efficiently.
//     
//         * Attachments are now closed properly!
//     
//         * Attachment Reader method is now required to return something satisfying
//           the io.ReadCloser interface. If no Close is required, wrap the return
//           value in an `io.NopCloser`.
//     
//         * The Envelope struct no longer has a message field - instead, use
//           an (Envelope, Message) 2-tuple where you need both of these items.
//     
//         * An email's Message-ID header is no longer implicitly generated for an
//           email. This is left to the mail submission agent.
//     
//         * If you ARE implementing a mail submission agent, an email's Message-ID
//           header can be specified by the new ID field on the Message struct type.
//     
//         * A cryptographically unique Message ID can be generated from the newly
//           exposed function, NewMessageID.
//     
//         * The Print method on Message is renamed Write.
//     
//         * Email message lines longer than 998 characters are now supported in
//           headers using folding white space. Note that some parsers, such as Go's
//           `net.mail`, do not understand this syntax (even though it is allowed).
//     
//         * The new method WriteCompat on Message won't use folding white space to
//           support long headers and will instead generate an error. Use this method
//           in preference to Write if you are expecting the consumer of your email
//           message (e.g. a Go implementation) will be unable to handle folding white
//           space.
//     
//     
package email // import "tawesoft.co.uk/go/email"

// Code generated by internal. DO NOT EDIT.
// Instead, edit DESC.txt and run mkdocs.sh.