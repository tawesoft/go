// tawesoft.co.uk/go/dialog
// 
// Copyright © 2019 - 2020 Ben Golightly <ben@tawesoft.co.uk>
// Copyright © 2019 - 2020 Tawesoft Ltd <opensource@tawesoft.co.uk>
// 
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction,  including without limitation the rights
// to use,  copy, modify,  merge,  publish, distribute, sublicense,  and/or sell
// copies  of  the  Software,  and  to  permit persons  to whom  the Software is
// furnished to do so.
// 
// THE SOFTWARE IS PROVIDED  "AS IS",  WITHOUT WARRANTY OF ANY KIND,  EXPRESS OR
// IMPLIED,  INCLUDING  BUT  NOT LIMITED TO THE WARRANTIES  OF  MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE  AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
// AUTHORS  OR COPYRIGHT HOLDERS  BE LIABLE  FOR ANY  CLAIM,  DAMAGES  OR  OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package dialog implements simple cross platform native MessageBox/Alert
// dialogs for Go.
// 
// Currently, only supports Windows and Linux targets.
// 
// On Linux, uses (in order of preference) `zenity`, `xmessage`, or stdio.
// 
// Example
// 
// Usage is quite simple:
// 
//     package main
// 
//     import "tawesoft.co.uk/go/dialog"
// 
//     func main() {
//         dialog.Alert("Hello world!")
//         dialog.Alert("There are %d lights", 4)
//     }
//
// Package Information
//
// License: MIT-0 (see LICENSE.txt)
//
// Stable: yes
//
// For more information, documentation, source code, examples, support, links,
// etc. please see https://www.tawesoft.co.uk/go and 
// https://www.tawesoft.co.uk/go/dialog
//
//     2019-11-16
//     
//         * Fix incorrect formatting of multiple arguments in Linux stdio fallback
//     
//     2019-10-16
//     
//         * Remove title argument from Alert function
//     
//     2019-10-01
//     
//         * Fix string formatting bug in Windows build
//     
//     2019-10-01
//     
//         * Support Unicode in UTF16 Windows dialogs
//         * Use "golang.org/x/sys/windows" to provide WinAPI
//         * Removes CGO and windows.h implementation
//         * Linux stdio fallback alert no longer blocks waiting for input
//     
//     2019-09-30
//     
//         * First release
//     
package dialog // import "tawesoft.co.uk/go/dialog"

// Code generated by internal. DO NOT EDIT.
// Instead, edit DESC.txt and run mkdocs.sh.