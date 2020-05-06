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

// TestRouteMatch tests matching individual path components against a single route
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
        {"foo/bar", Route{Pattern: "foo", Final: true},             false},
        {"foo/bar", Route{Pattern: "foo/bar", Final: true},         true},
    }
    
    for i, row := range table {
        match := row.route.match(row.component)
        if match != row.expected {
            t.Errorf("%d: expected %t match for path %s and route %+v", i, row.expected, row.component, row.route)
        }
    }
}

// TestRouterMatching tests matching paths against a router i.e. against a whole tree of routes
func TestRouterMatching(t *testing.T) {
    type Row struct {
        route Route
        matches map[string]string
    }
    
    table := []Row{
        {
            route: Route{Name: "Root", Children: []Route{
                {Name: "Home", Pattern: ""},
                {Name: "Foo", Pattern: "foo", Children: []Route{
                    {Name: "Bar", Pattern: "bar", Methods: "GET, POST"},
                }},
            }},
            matches: map[string]string{
                "": "Root",
                "/": "Home",
                "/foo": "Foo",
                "/foo/bar": "Bar",
                "other": "404",
                "/other": "404",
            },
        },
    }
    
    for i, row := range table {
        router, err := New(&row.route)
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
        
        for k, v := range(row.matches) {
            match := router.Match("GET", k)
            if match == nil && v != "404" {
                t.Errorf("%d: unexpected nil match for expected %v => %v", i, k, v)
            }
            if match != nil && v != match.Route.Name {
                t.Errorf("%d: unexpected match %v for expected %v => %v", i, match.Route.Name, k, v)
            }
        }
    }
}
