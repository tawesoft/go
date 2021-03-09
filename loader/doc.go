// tawesoft.co.uk/go/loader
// 
// Copyright © 2021 Tawesoft Ltd <open-source@tawesoft.co.uk>
// Copyright © 2021 Ben Golightly <ben@tawesoft.co.uk>
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

// Package loader implements the ability to define a graph of tasks and
// dependencies, classes of synchronous and concurrent workers, and limiting
// strategies, and solve the graph incrementally or totally.
// 
// For example, this could be used to implement a loading screen for a computer
// game with a progress bar that updates in real time, with images being decoded
// concurrently with files being loaded from disk, and synchronised with the main
// thread for safe OpenGL operations such as creating texture objects on the GPU.
// 
// While this package is generally suitable for use in real world applications,
// we are waiting to get some experience with how it works for us in an internal
// application before polishing or committing to a stable API.
// 
// TODO: doesn't yet free temporary results
// 
// TODO: refactor the load loop to always send/receive at the same time
// 
// TODO: clean up generally
// 
// TODO: not decided about the API for Loader.Result (but loader.MustResult is ok)
// 
// TODO: a step to simplify the DAG to remove passthrough loader.NamedTask steps
// 
// Examples
// 
// Configure the Loader with a Strategy to limit concurrent connections per host
//
// https://www.tawesoft.co.uk/go/doc/loader/examples/limit-connections-per-host/
//
//
// Package Information
//
// License: MIT (see LICENSE.txt)
//
// Stable: no
//
// For more information, documentation, source code, examples, support, links,
// etc. please see https://www.tawesoft.co.uk/go and 
// https://www.tawesoft.co.uk/go/loader
package loader // import "tawesoft.co.uk/go/loader"

// Code generated by internal. DO NOT EDIT.
// Instead, edit DESC.txt and run mkdocs.sh.