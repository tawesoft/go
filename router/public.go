package router

import (
    "fmt"
    "net/http"
    "regexp"
    "strings"
)

// Router is a general purpose router of a HTTP method (like "GET") and path (like "/foo/bar") to a tree of
// Route objects.
type Router struct {
    // DefaultMethods is the comma-separated list of HTTP methods to use if a route has none specified.
    // `router.New()` sets this to only "GET" by default.
    DefaultMethods string
    
    root *Route
    // methods []string // sorted
    routesByName map[string]*Route
    parents map[*Route]*Route
}

// Route describes a mapping between a component of a request path and a way to handle the request. Each route
// may contain child routes which map to child paths: each route is a tree. Each route works on a single
// component of the path recursively. A path such as "/user/123/profile" resolves into the components
// ["", "user", "123", "profile"]. A path such as "/user/123/profile/" resolves into the components
// ["", "user", "123", "profile", ""].
type Route struct {
    // Name (optional) is a way of uniquely identifying a specific route. It is used to generate a path from
    // a route (i.e. the reverse of routing). The name must be unique to the whole router.
    Name string
    
    // Key (optional) is a way of uniquely identifying a path parameter (e.g. the path "/user/123/profile" might have
    // a parameter for the user ID). It is used to query the value of the path parameter by name in a Match result.
    // The name must be unique to a route and its parents, but may be shared across different subtrees of routes.
    Key string
    
    // Pattern (optional) is a way of matching against a path component. This may be a literal string e.g. "user", or
    // or a compiled regular expression (https://golang.org/pkg/regexp/) such as regexp.MustCompile(`\d{1,5}`) (one to
    // five ASCII digits). Note that the regexp implementation provided by regexp is guaranteed to run in time linear
    // in the size of the input so it is generally not essential to limit the length of the path component. Path
    // components are limited to 255 characters regardless.
    //
    // If the pattern is nil or an empty string, the route matches an empty path component. This is useful for the
    // root node, or index pages for folders.
    //
    // As a special case, the string "*" matches everything. This is useful for a per-directory catch-all.
    //
    // Patterns are matched in the following order first, then in order of sequence defined:
    //
    // * Empty path (Pattern is nil or empty string; Final has no effect)
    // * Final Exact matches (Pattern is string and Final is true)
    // * Exact matches (Pattern is string and Final is false)
    // * Final Regex matches (Pattern is Regexp and Final is true)
    // * Regex matches (Pattern is Regexp and Final is false)
    // * Wildcard (string "*" and Final is false)
    // * Final Wildcard (string "*" and Final is true)
    Pattern interface{}
    
    // Methods (optional) is a comma-separated string of HTTP methods accepted by the route e.g. "GET, POST".
    // If left empty, implies only "GET". Note that "OPTIONS" is not handled automatically - this is up to you.
    Methods string
    
    // Handler (optional) is the information attached to a route e.g. a HTTP Handler. If nil, the route is
    // never a final match, but its child routes may match as normal. It is up to the caller to do what they will
    // with the resulting handler.
    Handler interface{}
    
    // Children (optional) are any child nodes that route child path components. For example, the path "/foo/bar"
    // decomposes into the path components ["", "foo", "bar"], and can be matched by a route with pattern "" and
    // child route with pattern "foo" which itself has a child route "bar".
    Children []Route
    
    // Final (optional; default false), if true, indicates that the pattern should be matched against the entire
    // remaining path, not just the current path component.
    //
    // For example, a Final route might have a pattern to accept arbitrary files in subdirectories, like
    // "static/assets/foo.png": regexp.MustCompile(`(\w+/)*\w+\.\w`)
    Final bool
}

// Match is the result of a routing query.
type Match struct {
    // Route is the route picked, or nil if no matching route was found.
    Route *Route
    
    // mapping of route keys => path parameters; may be nil
    params map[string]string
}

// Value returns a a parsed path parameter (e.g. "/user/123/profile") identified by a Route Key.
func (match Match) Value(key string) string {
    if match.params == nil { return "" }
    return match.params[key]
}

// Creates a new router with DefaultMethods set to GET
func New(root *Route) (*Router, error) {
    router := &Router{}
    
    router.DefaultMethods = "GET"
    
    router.root = root
    sortRoutes(router.root.Children) // note, must be done first!
    
    routesByName, err := scanRouteNames(router.root)
    if err != nil { return nil, err }
    router.routesByName = routesByName
    
    // router.methods = scanMethods(router.root)
    router.parents = scanParents(router.root)
    
    return router, nil
}

// Matches a HTTP request to a route. See Match
func (router *Router) MatchHttpRequest(r *http.Request) *Match {
    return router.Match(r.Method, r.URL.Path)
}

// Match attempts to match a method (e.g. a HTTP method like "GET") and path (e.g. "/foo/bar") to a route in
// a router's tree of routes. In the event that there is no match, returns nil.
func (router *Router) Match(method string, path string) *Match {
    var params map[string]string
    
    route := router.match(method, path, router.root, &params)
    if route == nil { return nil }
    
    if !route.matchMethod(method, router.DefaultMethods) { return nil }
    
    return &Match{route, params}
}

// Parent returns the parent Route of a Route in the router's tree of Routes. May be nil (i.e. at the root).
func (router *Router) Parent(route *Route) *Route {
    return router.parents[route]
}

// Format creates a URL for a route - the opposite of routing. Any regexp.Regexp patterns are replaced using
// each arg in sequence.
func (router *Router) Format(route *Route, args ... string) (string, error) {
    components := make([]string, 0)
    current := route
    index := len(args) - 1
    
    for current != nil {
        
        switch v := current.Pattern.(type) {
            case nil:
                components = append(components, "")
            case string:
                components = append(components, v)
            case *regexp.Regexp:
                if index < 0 {
                    return "", fmt.Errorf("not enough arguments to format path")
                }
                components = append(components, args[index])
                index--
            default:
                panic(fmt.Errorf("invalid Route Pattern type %T", v))
        }
        
        current = router.Parent(current)
    }
    
    //fmt.Println(components)
    reverseStringList(components...)
    
    return strings.Join(components, "/"), nil
}

// MustFormat is like Format, but panics on error.
func (router *Router) MustFormat(route *Route, args ... string) string {
    result, err := router.Format(route, args...)
    if err != nil {
        panic(fmt.Errorf("unable to format a route: %v", err))
    }
    return result
}
// Named returns a named Route where `name == route.Name`. Nil if not found.
func (router *Router) Named(name string) *Route {
    return router.routesByName[name]
}

// MustNamed is like Named, but panics if not found.
func (router *Router) MustNamed(name string) *Route {
    route := router.routesByName[name]
    if route == nil {
        panic(fmt.Errorf("named route %q not found", name))
    }
    return route
}
