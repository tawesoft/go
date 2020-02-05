package dialog // import "tawesoft.co.uk/go/dialog"

import (
    "strings"
)

// quickly wraps a message for error output
func wrap(message string, length int) string {
    var atoms = strings.Fields(strings.TrimSpace(message))
    var results = make([]string, 0, 16)
    var currentLength int

    for _, atom := range atoms {
        // special case for an atom longer than a whole line
        if (currentLength == 0) && (len(atom) >= length) {
            results = append(results, atom)
            results = append(results, "\n")
            currentLength = 0
            continue
        }

        // will overflow?
        if currentLength + len(atom) + 1 > length {
            results = append(results, "\n")
            currentLength = 0
        }

        // mid-line?
        if currentLength > 0 {
            results = append(results, " ")
            currentLength += 1
        }

        results = append(results, atom)
        currentLength += len(atom)
    }

    return strings.TrimSpace(strings.Join(results, ""))
}
