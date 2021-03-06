// tawesoft.co.uk/go/queue
// 
// Copyright © 2020 Ben Golightly <ben@tawesoft.co.uk>
// Copyright © 2020 Tawesoft Ltd <opensource@tawesoft.co.uk>
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

// Package queue implements simple, durable/ACID, same-process message queues
// with best-effort ordering by priority and/or time.
// 
// Examples
// 
// See examples folder.
//
// Package Information
//
// License: MIT (see LICENSE.txt)
//
// Stable: candidate
//
// For more information, documentation, source code, examples, support, links,
// etc. please see https://www.tawesoft.co.uk/go and 
// https://www.tawesoft.co.uk/go/queue
//
//     2021-07-06
//     
//         * The Queue RetryItem method now takes an Attempt parameter.
//     
//         * Calling the Delete() method on a Queue now attempts to avoid deleting
//           an in-memory database opened as ":memory:".
//     
package queue // import "tawesoft.co.uk/go/queue"

// Code generated by internal. DO NOT EDIT.
// Instead, edit DESC.txt and run mkdocs.sh.