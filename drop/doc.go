// tawesoft.co.uk/go/drop
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
// 

// Package drop implements the ability to start a process as root, open
// privileged resources as files, drop privileges to become a given user account,
// and inherit file handles across the dropping of privileges.
// 
// Examples
// 
// Opens privileged files and ports as root, then drops privileges
//
// https://www.tawesoft.co.uk/go/doc/drop/drop/
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
// https://www.tawesoft.co.uk/go/drop
//
//     2020-11-27
//     
//         * Drop() functionality has been moved to tawesoft.co.uk/go/drop with
//           changes to Inheritables from a struct to an interface.
//     
package drop // import "tawesoft.co.uk/go/drop"

// Code generated by internal. DO NOT EDIT.
// Instead, edit DESC.txt and run mkdocs.sh.