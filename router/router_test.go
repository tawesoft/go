package router

import (
    "regexp"
    "strconv"
    "testing"
)

// TestParseMethod tests that the first comma separated HTTP verbs parses correctly
func TestParseMethod(t *testing.T) {
    type Row struct {
        s string
        len int // expected
        advance int // expected
    }
    
    rows := []Row{
        {"!!!", 0, 0},
        {"FOO", 3, 3},
        {"FOO, BAR", 3, 5},
    }
    
    for i, row := range rows {
        len, advance := acceptMethod(row.s)
        if (row.len != len) || (row.advance != advance) {
            t.Errorf("%d: unexpected result %d, %d for %s: expected %d, %d",
                i, len, advance, row.s, row.len, row.advance)
        }
    }
}

/*
// TestRouterMethods tests that the comma separated HTTP verbs parse correctly, each route is visited,
// and the sorted list of methods is compiled correclt.y
func TestRouterMethods(t *testing.T) {
    route := &Route{Methods: "GET, POST, DELETE", Children: []Route{
        {Methods: "PUT"},
        {Methods: ""},
        {Methods: "TRACE", Children: []Route{
            {Methods: "VERSION-CONTROL, PATCH"},
        }},
    }}
    
    router, err := New(route)
    if err != nil { panic(err) }
    
    expected := []string{
        "DELETE",
        "GET",
        "PATCH",
        "POST",
        "PUT",
        "TRACE",
        "VERSION-CONTROL",
    }
    
    if diff := deep.Equal(router.methods, expected); diff != nil {
        t.Error(diff)
    }
}
*/

// TestRouteNames<Pass/Fail> tests that route name unique constraints hold
func TestUniqueRouteNamesPass(t *testing.T) {
    route := &Route{Name: "One", Children: []Route{
        {Name: "OneA"},
        {Name: ""},
        {Name: ""},
        {Name: "OneB", Children: []Route{
            {Name: "OneBA"},
            {Name: "OneBB"},
        }},
    }}
    
    _, err := New(route)
    if err != nil { t.Errorf("expected err == nil, got %v", err) }
}

func TestUniqueRouteNamesFail(t *testing.T) {
    route := &Route{Name: "One", Children: []Route{
        {Name: "OneA"}, // not unique
        {Name: ""},
        {Name: ""},
        {Name: "OneB", Children: []Route{
            {Name: "OneBA"},
            {Name: "OneA"}, // not unique
        }},
    }}
    
    _, err := New(route)
    if err == nil { t.Errorf("expected err != nil, but got nil err") }
}

// TestRouteOrdering tests that route children are ordered according to the spec.
func TestRouteOrdering(t *testing.T) {
    
    /*
    // * Empty path (Pattern is nil or empty string; Final has no effect)
    // * Final Exact matches (Pattern is string and Final is true)
    // * Exact matches (Pattern is string and Final is false)
    // * Final Regex matches (Pattern is Regexp and Final is true)
    // * Regex matches (Pattern is Regexp and Final is false)
    // * Wildcard (string "*" and Final is false)
    // * Final Wildcard (string "*" and Final is true)
    */
    
    route := &Route{Children: []Route{
        {Name: "2", Pattern: "foo"},
        {Name: "0", Pattern: ""},
        {Name: "3", Pattern: "bar"},
        {Name: "1", Pattern: "baz", Final: true},
        {Name: "5", Pattern: regexp.MustCompile(`b\d*`)},
        {Name: "4", Pattern: regexp.MustCompile(`b\d*`), Final: true},
        {Name: "7", Pattern: "*", Final: true},
        {Name: "6", Pattern: "*"},
    }}
    
    router, err := New(route)
    if err != nil { panic(err) }
    
    for i, r := range router.root.Children {
        name := strconv.Itoa(i)
        if r.Name != name {
            t.Errorf("%d: got unexpected item ordering with route.Name %s", i, route.Name)
        }
    }
}

// TestRouteMatch tests matching individual path components against a single route.
//
// Note that this test assumes the path component is Final.
func TestRouteMatch(t *testing.T) {
    type Row struct {
        component string
        route Route
        expected bool
    }
    
    table := []Row{
        {"",        Route{Pattern: ""},                             true},
        {"",        Route{Pattern: "foo"},                          false},
        {"foo",     Route{Pattern: "foo"},                          true},
        {"foo",     Route{Pattern: regexp.MustCompile(`^[a-z]+$`)}, true},
        {"bar",     Route{Pattern: "foo"},                          false},
        {"foo/bar", Route{Pattern: "foo"},                          false},
        {"foo",     Route{Pattern: "*"},                            true},
    }
    
    for i, row := range table {
        match := row.route.match(row.component)
        if match != row.expected {
            t.Errorf("%d: expected %t match for path %s and route %+v", i, row.expected, row.component, row.route)
        }
    }
}

// TestRouteParents tests the lookup of a route's parent route
func TestRouteParent(t *testing.T) {
    type Test struct {
        name string
        parentName string
    }
    
    routes := &Route{Name: "Root", Children: []Route{
        {Name: "Home", Pattern: ""},
        {Name: "Foo", Pattern: "foo", Children: []Route{
            {Name: "Bar", Pattern: "bar", Methods: "GET, POST"},
        }},
    }}
    
    tests := []Test{
        {"Root", "NIL"},
        {"Home", "Root"},
        {"Foo", "Root"},
        {"Bar", "Foo"},
    }
    
    router, err := New(routes)
    if err != nil { panic(err) }
    
    for i, test := range(tests) {
        route := router.Named(test.name)
        
        if route == nil {
            t.Errorf("%d: unable to find expected route named %s", i, test.name)
            continue
        }
        
        if route.Name != test.name {
            t.Errorf("%d: unexpected name %s, expected %s", i, route.Name, test.name)
            continue
        }
        
        parent := router.Parent(route)
        if parent == nil {
            if test.parentName != "NIL" {
                t.Errorf("%d: unexpected nil parent of %s", i, route.Name)
            }
        } else {
            if test.parentName != parent.Name {
                t.Errorf("%d: expected parent %s but got %s", i, test.parentName, parent.Name)
            }
        }
    }
}

// TestRouterMatching tests matching paths against a router i.e. against a whole tree of routes
func TestRouterMatching(t *testing.T) {
    
    routes := &Route{Name: "Root", Children: []Route{
        {Name: "Home", Pattern: ""},
        {Name: "Foo", Pattern: "foo", Children: []Route{
            {Name: "Bar", Pattern: "bar", Methods: "GET, POST"},
        }},
        {Name: "Users", Pattern: "users", Children: []Route{
            {Name: "User By ID", Key: "user-id", Pattern: regexp.MustCompile(`^\d+$`), Children: []Route{
                {Name: "User Profile", Pattern: "profile"},
            }},
            {Name: "Me", Pattern: "me", Children: []Route{
                {Name: "My Profile", Pattern: "profile"},
            }},
        }},
    }}
    matches := map[string]string{
        "": "Root",
        "/": "Home",
        
        "/foo": "Foo",
        "/foo/bar": "Bar",
        
        "/users": "Users",
        "/users/123": "User By ID",
        "/users/123/profile": "User Profile",
        "/users/me": "Me",
        "/users/me/profile": "My Profile",
        "/users/invalid/profile": "404",
        
        "other": "404",
        "/other": "404",
    }
    
    router, err := New(routes)
    if err != nil { panic(err) }
    
    // Check POST works
    match := router.Match("POST", "/foo/bar")
    if match == nil || match.Route.Name != "Bar" {
        t.Errorf("incorrect match for POST /foo/bar: %+v", match)
    }
    
    // Check POST to a GET-only route doesn't work
    match = router.Match("GET", "/foo")
    if match == nil || match.Route.Name != "Foo" {
        t.Errorf("incorrect match for GET /foo: %+v", match)
    }
    
    // Check match.Value works
    match = router.Match("GET", "/users/123/profile")
    if match == nil || match.Route.Name != "User Profile" {
        t.Errorf("incorrect match for GET /users/123/profile: %+v", match)
    }
    if match != nil && match.Value("user-id") != "123" && len(match.params) != 1 {
        t.Errorf("incorrect match params for GET /users/123/profile: %v", match.params)
    }
    
    for k, v := range(matches) {
        match := router.Match("GET", k)
        if match == nil && v != "404" {
            t.Errorf("unexpected nil match for expected %v => %v", k, v)
        }
        if match != nil && v != match.Route.Name {
            t.Errorf("unexpected match %v for expected %v => %v", match.Route.Name, k, v)
        }
    }
}

// TestRouteFormatting tests constructing a URL from a route
func TestRouteFormatting(t *testing.T) {
    type Test struct {
        routeName string
        args []string
        expected string
    }
    
    type Row struct {
        route Route
        tests []Test
    }
    
    table := []Row{
        {
            route: Route{Name: "Root", Children: []Route{
                {Name: "Home", Pattern: ""},
                {Name: "Foo", Pattern: "foo", Children: []Route{
                    {Name: "Bar", Pattern: "bar", Methods: "GET, POST"},
                }},
                {Name: "Users", Pattern: "users", Children: []Route{
                    {Name: "User By ID", Pattern: regexp.MustCompile(`^\d+$`), Children: []Route{
                        {Name: "User Profile", Pattern: "profile"},
                    }},
                    {Name: "Me", Pattern: "me", Children: []Route{
                        {Name: "My Profile", Pattern: "profile"},
                    }},
                }},
            }},
            tests: []Test{
                {"Root", []string{}, ""},
                {"Home", []string{}, "/"},
                {"Foo", []string{}, "/foo"},
                {"Bar", []string{}, "/foo/bar"},
                {"User By ID", []string{"123"}, "/users/123"},
                {"User Profile", []string{"123"}, "/users/123/profile"},
            },
        },
    }
    
    for i, row := range table {
        router, err := New(&row.route)
        if err != nil { panic(err) }

        for j, test := range(row.tests) {
            route := router.Named(test.routeName)
            
            path, err := router.Format(route, test.args...)
            if err != nil && test.expected != "ERROR" {
                t.Errorf("%d %d: unexpected error %v", i, j, err)
                continue
            }
            
            if path != test.expected {
                t.Errorf("%d %d: expected route format result %s but got %s", i, j, test.expected, path)
            }
        }
    }
}
