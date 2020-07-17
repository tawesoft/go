// USAGE: run and then get http://localhost:8080/users/1.json
// or http://localhost:8080/users/2.xml

package main

import (
    "encoding/json"
    "encoding/xml"
    "fmt"
    "log"
    "net/http"
    "regexp"
    "strconv"
    
    "tawesoft.co.uk/go/router"
)

type User struct {
    Name string
    DOB string
}

var Users = map[int]User{
    0: {"Alan Turing",  "23 June 1912"},
    1: {"Grace Hopper", "9 December 1906"},
    2: {"Donald Knuth", "10 January 1938"},
}

func HandleUserById(w http.ResponseWriter, r *http.Request, match *router.Match) {
    id, err := strconv.Atoi(match.Submatch("id", 1))
    if err != nil { panic(err) }
    fmt := match.Submatch("id", 2)
    
    user := Users[id] // range handling not done
    var out []byte
    var werr error
    
    switch fmt {
        case "": fallthrough // default to JSON
        case ".json":
            w.Header().Set("Content-Type", "application/json; charset=utf-8")
            out, werr = json.Marshal(user)
        case ".xml":
            w.Header().Set("Content-Type", "application/xml; charset=utf-8")
            out, werr = xml.Marshal(user)
        default:
            panic("format not implemented")
    }
    
    if werr != nil { panic(werr) }
    w.Write(out)
    
}

type MyHandlerType func (http.ResponseWriter, *http.Request, *router.Match)

func main() {
    id := regexp.MustCompile(`^(\d{1,9})(\.\w{1,5})?$`) // e.g. 1234, 1234.json
    
    type MyHandler struct {
        handle MyHandlerType
    }
    
    routes := router.Route{Children: []router.Route{
        {Pattern: "users", Children: []router.Route{
            {Pattern: id, Key: "id", Methods: "GET", Handler: MyHandler{HandleUserById}},
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
