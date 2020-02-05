package humanize

import (
    "fmt"
)

type ParseError struct {
    Input string // the input
    Err  error  // the reason
}

func (e *ParseError) Error() string {
    return fmt.Sprintf("error parsing %s: %v", e.Input, e.Err)
}
