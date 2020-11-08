// Demonstrates HTTP routing with "submatch" patterns in a path component
// with a server at localhost:8080
package main

import (
    "encoding/json"
    "encoding/xml"
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
    DOB string
}

var Users = map[int]User{
    0: {"Alan Turing",  "23 June 1912"},
    1: {"Grace Hopper", "9 December 1906"},
    2: {"Donald Knuth", "10 January 1938"},
}

func HandleIndex(w http.ResponseWriter, r *http.Request, match *router.Match) {
    const tpl = `
<!DOCTYPE html>
<html>
    <head>
        <meta charset="UTF-8">
        <title>Users</title>
    </head>
    <body>
        <h1>Users</h1>
        {{range $i, $u := . }}
            <div>{{ $u.Name }}:
                <a href="/users/{{ $i }}.json">JSON</a> |
                <a href="/users/{{ $i }}.xml">XML</a>
            </div>
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

func HandleUserById(w http.ResponseWriter, r *http.Request, match *router.Match, id int, format string) {
    user := Users[id] // range handling not done
    var out []byte
    var werr error
    
    switch format {
        case "":
            fallthrough // default to JSON
        case ".json":
            w.Header().Set("Content-Type", "application/json; charset=utf-8")
            out, werr = json.Marshal(user)
        case ".xml":
            w.Header().Set("Content-Type", "application/xml; charset=utf-8")
            out, werr = xml.Marshal(user)
        default:
            panic("invalid format")
    }
    
    if werr != nil { panic(werr) }
    w.Write(out)
    
}

type MyNormalHandlerType func(http.ResponseWriter, *http.Request, *router.Match)
type MyAPIHandlerType func (http.ResponseWriter, *http.Request, *router.Match, int, string)

func main() {
    // Let's make a convention: our API endpoints all end with a numeric ID
    // and an optional format specifier, like .xml or .json e.g.
    // "/users/123.json"
    
    // A regex that captures two submatches
    // e.g. "1234" => "1234", ""
    // "1234.json" => "1234", ".json"
    id := regexp.MustCompile(`^(\d{1,9})(\.\w{1,5})?$`)
    idNumberMatch := 1 // first capture
    idFormatMatch := 2 // second capture
    
    routes := router.Route{Children: []router.Route{
        {Pattern: "users", Children: []router.Route{
            {Pattern: id, Key: "id", Methods: "GET", Handler: MyAPIHandlerType(HandleUserById)},
        }},
        {Handler: MyNormalHandlerType(HandleIndex)},
    }}
    
    router, err := router.New(routes)
    if err != nil { panic(err) }
    
    log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
        match := router.MatchHttpRequest(r)
        
        // how you handle the match is up to you!
        if match != nil {
            
            if apiMatch, ok := match.Route.Handler.(MyAPIHandlerType); ok {

                // Get the integer id as the first submatch of the id component
                id, err := strconv.Atoi(match.Submatch("id", idNumberMatch))
                if err != nil { panic(err) }
                
                // Get the format as the second submatch of the id component
                format := match.Submatch("id", idFormatMatch)
            
                apiMatch(w, r, match, id, format)
            } else {
                match.Route.Handler.(MyNormalHandlerType)(w, r, match)
            }

        } else {
            http.NotFound(w, r)
            fmt.Fprintf(w, "Not Found: Sorry, no matching route!")
        }
    })))
}
