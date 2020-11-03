// Manages markdown/text/golang docstrings across the monorepo
package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "os"
    "sort"
    "strings"
    "unicode"
)

var PackageNames = []string{
    "dialog",
}

func Must(err error) {
    if err != nil { panic(err) }
}

func MustFile(fp *os.File, err error) *os.File {
    if err != nil { panic(err) }
    return fp
}

func MustString(x string, err error) string {
    if err != nil { panic(err) }
    return x
}

func main() {
    sort.Strings(PackageNames)
    packages := make([]Package, 0, len(PackageNames))
    for _, name := range PackageNames {
        packages = append(packages, LoadPackage(name))
    }
    
    for _, pkg := range packages {
        Must(pkg.writeGodoc())
        Must(pkg.writeReadme())
    }
}

type Doc struct {
    SPDXLicenseIdentifier string
    ShortDesc string
    Stable string
    Body string
}

type Package struct {
    Name string
    Doc Doc
    Changes string
    License string
}

func LoadPackage(name string) Package {
    return Package{
        Name:       name,
        Doc:        readDoc(name + "/DESC.txt"),
        Changes:    readOptionalFile(name + "/CHANGES.txt"),
        License:    readFile(name + "/LICENSE.txt"),
    }
}

func readFile(path string) string {
    b, err := ioutil.ReadFile(path)
    if err != nil { panic(err) }
    return string(b)
}

func readOptionalFile(path string) string {
    b, err := ioutil.ReadFile(path)
    if (err != nil) && os.IsNotExist(err) { return "" }
    if err != nil { panic(err) }
    return string(b)
}

func readDoc(path string) Doc {
    kv, body := readMetafile(path)
    stable := kv["stable"]
    if !(
        (stable == "yes") ||
        (stable == "no") ||
        (stable == "candidate")) {
        panic("invalid stable value for " + path)
    }
    
    return Doc{
        SPDXLicenseIdentifier: kv["SPDX-License-Identifier"],
        ShortDesc: kv["short-desc"],
        Stable: stable,
        Body: body,
    }
}

// readMetafile parses a doc.txt file with a k:v store followed by a line of three
// dashes followed by a body.
func readMetafile(path string) (map[string]string, string) {
    fp := MustFile(os.Open(path))
    defer fp.Close()
    
    inBody := false
    kvStore := make(map[string]string)
    bodyParts := make([]string, 0)
    
    scanner := bufio.NewScanner(fp)
    for scanner.Scan() {
        s := scanner.Text()
        
        if !inBody && (s == "---") {
            inBody = true
            continue
        }
        
        if inBody {
            bodyParts = append(bodyParts, s)
            continue
        }
        
        if s == "" {
            continue
        }
        
        parts := strings.SplitN(s, ":", 2)
        if len(parts) != 2 {
            panic(fmt.Sprintf("error parsing k:v store for %s: %+v", path, parts))
        }
        left  := strings.TrimSpace(parts[0])
        right := strings.TrimSpace(parts[1])
        kvStore[left] = right
    }
    Must(scanner.Err())
    
    return kvStore, strings.TrimSpace(strings.Join(bodyParts, "\n"))
}

// writeGodoc generates a doc.go file with a docstring for a given package
// from DESC.txt, CHANGES.txt (optional), and LICENSE.txt
func (p Package) writeGodoc() error {
    data := make([]string, 0)
    
    parts := strings.Split(p.License, "\n")
    for _, part := range parts {
        data = append(data, "// "+part)
    }

    data = append(data, "")
    
    parts = strings.Split(p.Doc.Body, "\n")
    for _, part := range parts {
        data = append(data, "// "+part)
    }
    
    data = append(data,
        "//",
        "// Package Information",
        "//",
        "// License: "+p.Doc.SPDXLicenseIdentifier+" (see LICENSE.txt)",
        "// Stable? "+p.Doc.Stable,
        "//",
        "// For more information, documentation, source code, examples, support, links,",
        "// etc. please see https://www.tawesoft.co.uk/go and ",
        "// https://www.tawesoft.co.uk/go/"+p.Name,
    )
    
    if p.Changes != "" {
        data = append(data,
            "//",
            "// Changes",
            "//",
        )
        
        parts = strings.Split(p.Changes, "\n")
        for _, part := range parts {
            data = append(data, "//     "+part)
        }
    }
    
    data = append(data,
        fmt.Sprintf(`package %[1]s // import "tawesoft.co.uk/go/%[1]s"`, p.Name),
        "",
        "// Code generated by internal. DO NOT EDIT.",
        "// Instead, edit .txt files and `go run internal/doc/doc.go`.",
    )
    
    bdata := []byte(strings.Join(data, "\n"))
    return ioutil.WriteFile(p.Name + "/doc.go", bdata, 0644)
}

// writeGodoc generates a README.md file for a given package from DESC.txt,
// CHANGES.txt (optional), and LICENSE.txt
func (p Package) writeReadme() error {
    data := make([]string, 0)
    
    data = append(data,
        fmt.Sprintf("# %s - %s", p.Name, p.Doc.ShortDesc),
        "",
        "```shell script",
        `go get "tawesoft.co.uk/go/"`,
        "```",
        "",
        "```go",
        `import "tawesoft.co.uk/go/`+p.Name+`"`,
        "```",
        "",
    )
    
    data = append(data, strings.Split(markdownPackageTable(p), "\n")...)
    
    data = append(data,
        "",
        "## About",
        "",
    )
    
    inCodeBlock := false
    codePrefixLen := 0
    var codeData []string
    
    parts := strings.Split(p.Doc.Body, "\n")
    for i, part := range parts {
        if inCodeBlock {
            if strings.TrimSpace(part) == "" {
                codeData = append(codeData, "")
            } else if isGodocCode(part) {
                codeData = append(codeData, part[codePrefixLen:])
            } else {
                data = append(data, trimCodeblock(codeData)...)
                data = append(data, "```", "", part)
                inCodeBlock = false
            }
        } else if isGodocCode(part) {
            inCodeBlock = true
            codePrefixLen = len(part) - len(strings.TrimLeft(part, "\t "))
            codeData = make([]string, 0)
            data = append(data, "```go")
            codeData = append(codeData, part[codePrefixLen:])
        } else if isGodocTitle(i, parts) {
            data = append(data, "## "+part)
        } else {
            data = append(data, part)
        }
    }
    
    if inCodeBlock {
        data = append(data, trimCodeblock(codeData)...)
        data = append(data, "```")
    }
    
    if p.Changes != "" {
        data = append(data, "", "## Changes", "")
        parts = strings.Split(p.Changes, "\n")
        for _, part := range parts {
            if strings.HasPrefix(part, " ") || strings.HasPrefix(part, "\t") {
                data = append(data, strings.TrimSpace(part))
            } else {
                data = append(data, "### "+part)
            }
        }
    }
    
    bdata := []byte(strings.Join(data, "\n"))
    return ioutil.WriteFile(p.Name + "/README.md", bdata, 0644)
}

// isGodocTitle detects titles in a Go docstring. A title is a line that is
// separated from its following line by an empty line, begins with a capital
// letter and doesn't end with punctuation.
func isGodocTitle(i int, lines []string) bool {
    if i + 1 == len(lines) { return false }
    
    line := lines[i]
    nextLine := lines[i+1]
    
    if strings.TrimSpace(nextLine) != "" { return false }
    if len(line) == 0 { return false }
    
    // FIXME: NOT i18n friendly
    firstChar, lastChar := line[0], line[len(line)-1]
    if (firstChar < 'A') || (firstChar > 'Z') { return false }
    if unicode.IsPunct(rune(lastChar)) { return false }
    
    return true
}

// isGodocCode detects a code block in a Go docstring. A codeblock is an
// indented line.
func isGodocCode(line string) bool {
    return strings.HasPrefix(line, "\t") || strings.HasPrefix(line, " ")
}

func trimCodeblock(xs []string) []string {
    x := strings.Join(xs, "\n")
    x = strings.TrimSpace(x)
    return strings.Split(x, "\n")
}

func markdownPackageTable(p Package) string {
    template := strings.TrimSpace(`
|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_%[1]s] ∙ [docs][docs_%[1]s] ∙ [src][src_%[1]s] | [%[2]s][copy_%[1]s] | %[3]s |

[home_%[1]s]: https://tawesoft.co.uk/go/%[1]s
[src_%[1]s]:  https://github.com/tawesoft/go/tree/master/dialog
[docs_%[1]s]: https://godoc.org/tawesoft.co.uk/go/%[1]s
[copy_%[1]s]: https://github.com/tawesoft/go/tree/master/%[1]s/LICENSE.txt
`)

    stable := p.Doc.Stable
    if stable == "yes" {
        stable = "✔ yes"
    } else if stable == "no" {
        stable = "✘ **no**"
    }
    
    return fmt.Sprintf(template, p.Name, p.Doc.SPDXLicenseIdentifier, stable)
}
