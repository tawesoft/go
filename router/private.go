package router

import (
    "fmt"
    "regexp"
    "sort"
)

// visit applies `f(route)` to a tree of routes (DFS)
func visit(route *Route, f func(route *Route)) {
    f(route)
    for _, child := range route.Children {
        visit(&child, f)
    }
}

// visiterr applies `f(route) => error` to a tree of routes (DFS), aborting on error
func visiterr(route *Route, f func(route *Route) error) error {
    err := f(route)
    if err != nil { return err }
    
    for _, child := range route.Children {
        err = visiterr(&child, f)
        if err != nil { return err }
    }
    
    return nil
}

// visitp applies `f(route, parent)` to a tree of routes (DFS)
func visitp(route *Route, parent *Route, f func(route *Route, parent *Route)) {
    f(route, parent)
    for _, child := range route.Children {
        visitp(&child, parent, f)
    }
}

/*
// sortedStringSetAdd inserts an item into a sorted list of strings if it does not exist already.
func sortedStringSetAdd(list []string, item string) []string {
    var index = sort.SearchStrings(list, item)
    if (len(list) <= index) || (list[index] != item) {
        list = append(list, item)
        sort.Strings(list)
    }
    return list
}
*/

// acceptMethod parses [A-Z-]+ and returns a length plus a length including trailing spaces and commas.
func acceptMethod(s string) (len int, offset int) {
    for _, c := range s {
        if (c >= 'A' && c <= 'Z') || (c == '-') {
            if offset != len { break } // next token appearing after commas/whitespace
            len++
            offset++
        } else if (c == ',') || (c == ' ') || (c == '\t') {
            offset++
        } else {
            break
        }
    }
    
    return len, offset
}

// nextMethod parses a comma-separated list of methods at a given offset and returns the offset of the
// next method.
func nextMethod(methods string, offset int) (string, int) {
    length, advance := acceptMethod(methods[offset:])
    if length > 0 {
        return methods[offset:offset+length], advance
    } else {
        return "", -1
    }
}

// acceptPathComponent parses up to the next '/' and returns a length excluding that '/'
func acceptPathComponent(s string) (offset int) {
    for _, c := range s {
        if (c == '/') { break }
        offset++
    }
    
    return offset
}

/*
// scanMethods visits each route in a tree, building a sorted list of methods.
func scanMethods(root *Route) []string {
    methods := make([]string, 0)
    
    visit(root, func(route *Route) {
        offset := 0
        for {
            method, advance := nextMethod(route.Methods, offset)
            if advance < 0 { break }
            offset += advance
            methods = sortedStringSetAdd(methods, method)
        }
    })
    
    return methods
}
 */

// scnaRouteNames builds a mapping of `name => Route`
func scanRouteNames(root *Route) (map[string]*Route, error) {
    namedRoutes := make(map[string]*Route)
    
    err := visiterr(root, func(route *Route) error {
        name := route.Name
        if len(name) == 0 { return nil }
        
        if namedRoutes[name] == nil {
            namedRoutes[name] = route
        } else {
            return fmt.Errorf("route name %s is not unique", name)
        }
        
        return nil
    })
    
    return namedRoutes, err
}

// builds a mapping of `Route => parent Route`
func scanParents(root *Route) map[*Route]*Route {
    parents := make(map[*Route]*Route)
    
    visitp(root, nil, func(route *Route, parent *Route) {
        parents[route] = parent
    })
    
    return parents
}

// reverseStringList reverses in place.
func reverseStringList(a ... string) {
    // https://github.com/golang/go/wiki/SliceTricks
    for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
        a[left], a[right] = a[right], a[left]
    }
}

type sortableRoutes []Route

func (a sortableRoutes) Len() int             { return len(a) }
func (a sortableRoutes) Swap(i, j int)        { a[i], a[j] = a[j], a[i] }
func (a sortableRoutes) Less(i, j int) bool   { return a[i].order() < a[j].order() }
func (r Route) order() int {
    /*
    // ORDER:
    // 0 Empty path (Pattern is nil or empty string; Final has no effect)
    // 1 Final Exact matches (Pattern is string and Final is true)
    // 2 Exact matches (Pattern is string and Final is false)
    // 3 Final Regex matches (Pattern is Regexp and Final is true)
    // 4 Regex matches (Pattern is Regexp and Final is false)
    // 5 Wildcard (string "*" and Final is false)
    // 6 Final Wildcard (string "*" and Final is true)
     */
    
    switch v := r.Pattern.(type) {
        case nil:
                                      return 0
        case string:
            if v == ""              { return 0 }
            if v == "*" && !r.Final { return 5 }
            if v == "*" && r.Final  { return 6 }
            if r.Final              { return 1 }
            if !r.Final             { return 2 }
        case *regexp.Regexp:
            if r.Final              { return 3 }
            if !r.Final             { return 4 }
        default:
            panic(fmt.Errorf("invalid Route Pattern type %T", v))
    }
    
    return 0
}

// sortRoutes sorts a list of routes in place such that
func sortRoutes(routes []Route) {
    if routes == nil { return }
    if len(routes) == 0 { return }

    sort.Stable(sortableRoutes(routes))
    for _, i := range routes {
        sortRoutes(i.Children)
    }
}

// match returns true if a route matches the given path component
func (r *Route) match(component string) bool {

    switch v := r.Pattern.(type) {
        case nil:
            return component == ""
        case string:
            if v == "" { return component == ""  }
            if v == "*" { return true }
            return v == component
        case *regexp.Regexp:
            return v.MatchString(component)
        default:
            panic(fmt.Errorf("invalid Route Pattern type %T", v))
    }
}


func (router *Router) match(method string, path string, current *Route, params map[string]string) *Route {
    var component string
    var advance int
    
    if current.Final {
        component = path
    } else {
        advance = acceptPathComponent(path)
        component = path[0:advance]
    }
    
    match := current.match(component)
    
    if match && current.Key != "" {
        params[current.Key] = component
    }
    
    if match && current.Final {
        // its the last match because its a Final pattern
        return current
    } else if match && (current.Children == nil || len(current.Children) == 0) {
        // its the last match because its a match with no children AND
        // the full path has been parsed.
        if advance == len(path) { return current }
    } else if match {
        // the full path has been matched, so its this route
        if advance == len(path) { return current }
        
        // otherwise its a partial match, that the child routes have to try and handle
        remainder := path[advance+1:]
        for _, i := range(current.Children) {
            match := router.match(method, remainder, &i, params)
            if match != nil { return match }
        }
    }
    
    return nil
}

func (r *Route) matchMethod(method string, defaultMethods string) bool {
    offset := 0
    
    methods := r.Methods
    if methods == "" { methods = defaultMethods }
    
    for {
        routeMethod, advance := nextMethod(methods, offset)
        if advance < 0 { break }
        offset += advance
        if method == routeMethod { return true }
    }
    
    return false
}
