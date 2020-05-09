package router

import (
    "fmt"
    "net/http"
    "regexp"
    "strings"
)

// Router is a general purpose router of a method or verb (like HTTP "GET") and
// path (like "/foo/bar") to a tree of Route objects.
type Router struct {
    // DefaultMethods is the comma-separated list of methods to use if a route
    // has none specified e.g. HTTP "GET, POST".
    //
    // router.New() sets this to only "GET" by default.
    DefaultMethods string
    
    root *Route
    routesByName map[string]*Route
    parents map[*Route]*Route
}

// Route describes a mapping between methods (like HTTP "GET") and path
// component (like "foo" in "/foo/bar") to a user-supplied handler value.
//
// Each route may contain child routes which map to subsequent path components.
// For example a route may match "foo" in "/foo/bar", and a child route could
// match "bar".
//
// A path such as "/user/123/profile" resolves into the components ["", "user",
// "123", "profile"]. A path such as "/user/123/profile/" resolves into the
// components ["", "user", "123", "profile", ""]. As such, the root Route
// and folder-index Routes will normally be configured to match the empty
// string "".
//
// A route can be Final, in which case it matches all remaining path components
// at once.
type Route struct {
    // Name (optional) is a way of uniquely identifying a specific route. The
    // name must be unique to the whole router tree.
    Name string
    
    // Key (optional) is a way of uniquely identifying a path parameter (e.g.
    // the path "/user/123/profile" might have a parameter for the user ID).
    //
    // This key is used to query the value of the path parameter by name in a
    // Match result.
    //
    // All keys must be unique to a route and its parents, but does not have to
    // be unique across the entire router tree.
    Key string
    
    // Pattern (optional) is a way of matching against a path component. This
    // may be a literal string e.g. "user", or or a compiled regular expression
    // (https://golang.org/pkg/regexp/) such as regexp.MustCompile(`$\d{1,5}^`)
    // (this matches one to five ASCII digits).
    //
    // Note that the regexp implementation provided by regexp is guaranteed to
    // run in time linear in the size of the input so it is generally not
    // essential to limit the length of the path component. Path components are
    // limited to 255 characters regardless.
    //
    // If the pattern is nil or an empty string, the route matches an empty
    // path component. This is useful for the root node, or index pages for
    // folders.
    //
    // As a special case, the string "*" matches everything. This is useful for
    // a per-directory catch-all.
    //
    // Patterns are matched in the following order first, then in order of
    // sequence defined:
    //
    // 1 Empty path (Pattern is nil or empty string; Final has no effect)
    // 2 Final Exact matches (Pattern is string and Final is true)
    // 3 Exact matches (Pattern is string and Final is false)
    // 4 Final Regex matches (Pattern is Regexp and Final is true)
    // 5 Regex matches (Pattern is Regexp and Final is false)
    // 6 Wildcard (string "*" and Final is false)
    // 7 Final Wildcard (string "*" and Final is true)
    Pattern interface{}
    
    // Methods (optional) is a comma-separated string of methods or verbs
    // accepted by the route e.g. HTTP "GET, POST". If left empty, implies
    // only the Router's DefaultMethods.
    //
    // Note that "OPTIONS" is not handled automatically - this is up to you.
    Methods string
    
    // Handler (optional) is the caller-supplied information attached to a
    // route e.g. a HTTP Handler or any custom value.
    //
    // It is up to the caller to do something with the resulting Handler -
    // nothing is called automatically for you.
    Handler interface{}
    
    // Children (optional) are any child Routes for subsequent path components.
    // For example, the path "/foo/bar" decomposes into the path components
    // ["", "foo", "bar"], and can be matched by a route with pattern "",
    // child route with pattern "foo", and grandchild route "bar".
    Children []Route
    
    // Final (optional; default false), if true, indicates that the pattern
    // should be matched against the entire remaining path, not just the
    // current path component.
    //
    // For example, a Route might have a pattern to accept arbitrary
    // files in subdirectories, like "static/assets/foo.png". It would be
    // Final with a Pattern value of something like
    // regexp.MustCompile(`(\w+/)*\w+\.\w`)
    Final bool
}

// Match is the result of a routing query - may be nil if no match found
type Match struct {
    // Route is the route picked - never nil
    Route *Route
    
    // Router is the router that was matched on
    Router *Router
    
    // mapping of route keys => path parameters; may be nil: use Params()
    params map[string]string
}

// Params returns a non-nil mapping of the route keys => path parameters
func (match Match) Values() map[string]string {
    if match.params == nil { match.params = make(map[string]string) }
    return match.params
}

// Value returns a a parsed path parameter identified by a Route Key.
//
// For example, you might create a tree of routes such that the path
// "/user/123/profile" captures the path component "123" and associates it
// with a "user-id" route.Key.
func (match Match) Value(key string) string {
    if match.params == nil { return "" }
    return match.params[key]
}

// Creates a new router with DefaultMethods set to GET
func New(root Route) (*Router, error) {
    router := &Router{}
    
    router.DefaultMethods = "GET"
    
    // we take a copy of the routes rather than just hold a pointer because
    // we modify the route Children by sorting them, and we return
    
    router.root = &root
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
    
    return &Match{route, router, params}
}

// Parent returns the parent Route of a Route in the router's tree of Routes. May be nil (i.e. at the root).
func (router *Router) Parent(route *Route) *Route {
    return router.parents[route]
}

// Format creates a URL for a route - the opposite of routing. Any
// regexp.Regexp patterns are replaced using each arg in sequence.
//
// WARNING: The return value from this function may be controlled by the
// User Agent. Escape it as necessary.
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

// Format creates a URL for a route - the opposite of routing. Any
// regexp.Regexp patterns are replaced using each the route Key to lookup
// a value in the the mapping argument. Pass in Match Values() to have
// this automatically filled based on the parsed request.
//
// WARNING: The return value from this function may be controlled by the
// User Agent. Escape it as necessary.
func (router *Router) FormatMap(route *Route, mapping map[string]string) (string, error) {
    components := make([]string, 0)
    current := route
    
    for current != nil {
        
        switch v := current.Pattern.(type) {
            case nil:
                components = append(components, "")
            case string:
                components = append(components, v)
            case *regexp.Regexp:
                components = append(components, mapping[current.Key])
            default:
                panic(fmt.Errorf("invalid Route Pattern type %T", v))
        }
        
        current = router.Parent(current)
    }
    
    reverseStringList(components...)
    
    return strings.Join(components, "/"), nil
}

// MustFormatMap is like FormatMap, but panics on error.
func (router *Router) MustFormatMap(route *Route, mapping map[string]string) string {
    result, err := router.FormatMap(route, mapping)
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
