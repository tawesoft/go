/*
Package router is a general purpose router of methods (e.g. HTTP "GET") and paths (e.g. "/user/123/profile") to
some value e.g. a controller.

Supports named routes, route parameters, constructing a path from a route, etc.

Although built with HTTP routing in mind, this is a general purpose implementation that can route to any type
of value - it is not limited to HTTP handlers.

Example (router/examples/example1/example1.go):

    package main
    
    import (
        "fmt"
        "html/template"
        "log"
        "net/http"
        "regexp"
        "strconv"
        
        "tawesoft.co.uk/go/router"
    )
    
    type User struct {
        Name string
    }
    
    const Me = 1
    
    var Users = map[int]User{
        0: {"Alan Turing"},
        1: {"Grace Hopper"},
        2: {"Donald Knuth"},
    }
    
    func HandleIndex(w http.ResponseWriter, r *http.Request, match *router.Match) {
        const tpl = `
    <!DOCTYPE html>
    <html>
        <head>
            <meta charset="UTF-8">
            <title>Home</title>
        </head>
        <body>
            <h1>Home</h1>
            <a href="/users">Users</a>
        </body>
    </html>`
        fmt.Fprintf(w, tpl)
    }
    
    func HandleUsersIndex(w http.ResponseWriter, r *http.Request, match *router.Match) {
        const tpl = `
    <!DOCTYPE html>
    <html>
        <head>
            <meta charset="UTF-8">
            <title>Users</title>
        </head>
        <body>
            <h1>Users</h1>
            <div><a href="/users/me"><b>My Profile</b></a></div>
            {{range $i, $u := .}}
                <div><a href="/users/{{ $i }}">{{ $u.Name }}</a></div>
            {{else}}
                <div><strong>no rows</strong></div>
            {{end}}
        </body>
    </html>`
        t, err := template.New("webpage").Parse(tpl)
        if err != nil { panic(err) }
        err = t.Execute(w, Users)
        if err != nil { panic(err) }
    }
    
    func HandleUserMe(w http.ResponseWriter, r *http.Request, match *router.Match) {
        fmt.Fprintf(w, "I am %s", Users[Me].Name)
    }
    
    func HandleUserById(w http.ResponseWriter, r *http.Request, match *router.Match) {
        id, err := strconv.Atoi(match.Value("id"))
        if err != nil { panic(err) }
        fmt.Fprintf(w, "I am %s", Users[id].Name)
    }
    
    type MyHandlerType func (http.ResponseWriter, *http.Request, *router.Match)
    
    func main() {
        id := regexp.MustCompile(`^\d{1,9}$`)
        
        type MyHandler struct {
            handle MyHandlerType
        }
        
        routes := &router.Route{Name: "Root", Children: []router.Route{
            {Name: "Home", Methods: "GET", Handler: MyHandler{HandleIndex}},
            {Name: "Users", Pattern: "users", Methods: "GET, POST", Handler: MyHandler{HandleUsersIndex}, Children: []router.Route{
                {Pattern: "me", Methods: "GET", Name: "My Profile", Handler: MyHandler{HandleUserMe}},
                {Pattern: id, Key: "id", Methods: "GET", Name: "User Profile", Handler: MyHandler{HandleUserById}},
            }},
        }}
        
        router, err := router.New(routes)
        if err != nil { panic(err) }
        
        log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
            match := router.MatchHttpRequest(r)
            
            // how you handle the match is up to you!
            if match != nil {
                match.Route.Handler.(MyHandler).handle(w, r, match)
            } else {
                http.NotFound(w, r)
                fmt.Fprintf(w, "Not Found: Sorry, no matching route!")
            }
        })))
    }



For license information, documentation, source code, support, links, etc. please see
https://tawesoft.co.uk/go/router

This module is part of https://tawesoft.co.uk/go
*/
package router // import "tawesoft.co.uk/go/router"

// SPDX-License-Identifier: MIT

// Code generated by (tawesoft.co.uk/go/) fluff.py: DO NOT EDIT.