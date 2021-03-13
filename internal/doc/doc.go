// Manages markdown/text/golang docstrings across the monorepo
package main

import (
    "bufio"
    "fmt"
    "html/template"
    "io/ioutil"
    "os"
    "sort"
    "strings"
    textTemplate "text/template"
    "time"
    "unicode"
)

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
    packageNames := os.Args[1:]
    sort.Strings(packageNames)
    packages := make([]Package, 0, len(packageNames))
    packagesExceptLegacy := make([]Package, 0, len(packageNames))
    for _, name := range packageNames {
        p := LoadPackage(name)
        packages = append(packages, p)

        if !strings.HasPrefix(p.Name, "legacy") {
            packagesExceptLegacy = append(packagesExceptLegacy, p)
        }
    }

    for _, pkg := range packages {
        Must(pkg.writeGodoc())
        Must(pkg.writeReadme())
    }

    t, err := template.New("page").Parse(readFile("internal/doc/template.html"))
    if err != nil { panic(err) }

    writeGoIndex(packages)
    writeMarkdownIndex(packagesExceptLegacy)
    writeHtmlIndex(t, packagesExceptLegacy)
    writeLicenseIndex(packages)
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

func (p Package) SPDXLicenseIdentifier() string { return p.Doc.SPDXLicenseIdentifier }
func (p Package) Stable() string { return p.Doc.Stable }
func (p Package) ShortDesc() string { return p.Doc.ShortDesc }
func (p Package) MediumDesc() string {
    return strings.Split(p.Doc.Body, "\n\n")[0]
}

// Returns a name suitable for a Go package identifier
func (p Package) GoName() string {
    // e.g. "ximage/xcolor" => "xcolor"
    parts := strings.Split(p.Name, "/")
    return parts[len(parts)-1]
}

func LoadPackage(name string) Package {
    return Package{
        Name:       name,
        Doc:        readDoc(name + "/DESC.txt"),
        Changes:    readOptionalFile(name + "/CHANGES.txt"),
        License:    strings.TrimSpace(readFile(name + "/LICENSE.txt")),
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

// readGo parses a go source file, returning a 2-tuple of (docstring, body).
func readGo(path string) (string, string) {
    fp := MustFile(os.Open(path))
    defer fp.Close()

    inBody := false
    docParts := make([]string, 0)
    bodyParts := make([]string, 0)

    scanner := bufio.NewScanner(fp)
    for scanner.Scan() {
        s := scanner.Text()

        if !inBody && (strings.HasPrefix(s, "package")) {
            inBody = true
        }

        if inBody {
            bodyParts = append(bodyParts, s)
            continue
        }

        if strings.HasPrefix(s, "//") {
            s = s[2:]
        }
        docParts = append(docParts, strings.TrimSpace(s))
    }
    Must(scanner.Err())

    return strings.TrimSpace(strings.Join(docParts, "\n")), strings.TrimSpace(strings.Join(bodyParts, "\n"))
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
        if strings.HasPrefix(part, "EXAMPLE:") {
            args := strings.SplitN(part, " ", 2)
            if len(args) != 2 { panic("Invalid EXAMPLE: command") }
            example := args[1]
            docstr, _ := readGo(fmt.Sprintf("%[1]s/examples/%[2]s/%[2]s.go",
                p.Name, example))
            for _, p := range strings.Split(docstr, "\n") {
                data = append(data, "// "+p)
            }
            data = append(data, "//")

            // rather than embed the example, provide a link
            // (relative would be preferable but wouldn't automatically get
            // hyperlinked)
            data = append(data, fmt.Sprintf(
                "// https://www.tawesoft.co.uk/go/doc/%s/examples/%s/", p.Name, example))
            /*
            for _, p := range strings.Split(code, "\n") {
                data = append(data, "//     "+p)
            }
             */
            data = append(data, "//")
        } else {
            data = append(data, "// "+part)
        }
    }

    data = append(data,
        "//",
        "// Package Information",
        "//",
        "// License: "+p.Doc.SPDXLicenseIdentifier+" (see LICENSE.txt)",
        "//",
        "// Stable: "+p.Doc.Stable,
        "//",
        "// For more information, documentation, source code, examples, support, links,",
        "// etc. please see https://www.tawesoft.co.uk/go and ",
        "// https://www.tawesoft.co.uk/go/"+p.Name,
    )

    if p.Changes != "" {
        data = append(data,
            "//",
        )

        parts = strings.Split(p.Changes, "\n")
        for _, part := range parts {
            data = append(data, "//     "+part)
        }
    }

    data = append(data,
        fmt.Sprintf(`package %[1]s // import "tawesoft.co.uk/go/%[2]s"`, p.GoName(), p.Name),
        "",
        "// Code generated by internal. DO NOT EDIT.",
        "// Instead, edit DESC.txt and run mkdocs.sh.",
    )

    bdata := []byte(strings.Join(data, "\n"))
    return ioutil.WriteFile(p.Name + "/doc.go", bdata, 0644)
}

// writeReadme generates a README.md file for a given package from DESC.txt,
// CHANGES.txt (optional), and LICENSE.txt
func (p Package) writeReadme() error {
    data := make([]string, 0)

    data = append(data,
        fmt.Sprintf("# %s - %s", p.Name, p.Doc.ShortDesc),
        "",
        "```shell script",
        `go get -u "tawesoft.co.uk/go"`,
        "```",
        "",
        "```go",
        `import "tawesoft.co.uk/go/`+p.Name+`"`,
        "```",
        "",
    )

    data = append(data, strings.Split(fmtMarkdownPackageTable(p), "\n")...)

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
            data = append(data, "", "```go")
            codeData = append(codeData, part[codePrefixLen:])
        } else if strings.HasPrefix(part, "EXAMPLE:") {
            args := strings.SplitN(part, " ", 2)
            if len(args) != 2 { panic("Invalid EXAMPLE: command") }
            example := args[1]
            docstr, code := readGo(fmt.Sprintf("%[1]s/examples/%[2]s/%[2]s.go",
                p.Name, example))
            data = append(data, docstr)
            data = append(data, "```go")
            data = append(data, code)
            data = append(data, "```")
        } else if isGodocTitle(i, parts) {
            data = append(data, "", "## "+part, "")
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
            if strings.HasPrefix(part, "20") { // 20XX-XX-XX
                data = append(data, "### "+part)
            } else {
                data = append(data, strings.TrimSpace(part))
            }
        }
    }

    data = append(data,
        "",
        "## Getting Help",
        "",
        "This package is part of [tawesoft.co.uk/go](https://www.tawesoft.co.uk/go),",
        "a monorepo for small Go modules maintained by Tawesoft®.",
        "Check out that URL for more information about other Go modules from",
        "Tawesoft plus community and commercial support options.",
    )

    bdata := []byte(strings.Join(data, "\n"))
    return ioutil.WriteFile(p.Name + "/README.md", bdata, 0644)
}

func writeGoIndex(packages []Package) {
    f, err := os.Create("tawesoft.go")
    if err != nil { panic(err) }
    defer f.Close()

    t, err := textTemplate.New("go").Parse(
`/*
A monorepo for small Go modules maintained by Tawesoft®

This is permissively-licensed open source software but exact licenses may vary between modules.

For license information, documentation, source code, support, links, etc. please see
https://www.tawesoft.co.uk/go
*/
package tawesoft

import (
{{- range . }}
    _ "tawesoft.co.uk/go/{{ .Name }}"
{{- end }}
)
`)
    if err != nil { panic(err) }

    err = t.Execute(f, packages)
    if err != nil { panic(fmt.Sprintf("template error: %v", err)) }
}

func writeMarkdownIndex(packages []Package) {

    data := make([]string, 0)

    data = append(data,
`[![Tawesoft](https://www.tawesoft.co.uk/media/0/logo-240r.png)](https://tawesoft.co.uk/go)
================================================================================

A monorepo for small Go modules maintained by [Tawesoft®](https://www.tawesoft.co.uk/go)

This is permissively-licensed open source software but exact licenses may vary between modules.

Download
--------

    go get -u tawesoft.co.uk/go

Contents
--------

`)

    for _, p := range packages {
        data = append(data,
            fmt.Sprintf("### %s - %s", p.Name, p.Doc.ShortDesc),
            "",
            p.MediumDesc(),
            "",
            "```go",
            `import "tawesoft.co.uk/go/`+p.Name+`"`,
            "```",
            "",
        )
        data = append(data, strings.Split(fmtMarkdownPackageTable(p), "\n")...)
    }

data = append(data, `
Links
-----

* Home: [tawesoft.co.uk/go](https://tawesoft.co.uk/go)
* Docs hub: [tawesoft.co.uk/go/doc/](https://www.tawesoft.co.uk/go/doc/)
* Repository: [github.com/tawesoft/go](https://github.com/tawesoft/go)
* Or [search "tawesoft"](https://pkg.go.dev/search?q=tawesoft) on [go.dev](https://go.dev/)

Support
-------

### Free and Community Support

* [GitHub issues](https://github.com/tawesoft/go/issues)
* Email open-source@tawesoft.co.uk (feedback welcomed, but support is "best
 effort")

### Commercial Support

Open source software from Tawesoft® backed by commercial support options.

Email open-source@tawesoft.co.uk or visit [tawesoft.co.uk/products/open-source-software](https://www.tawesoft.co.uk/products/open-source-software)
to learn more.
`)

    bdata := []byte(strings.Join(data, "\n"))
    err := ioutil.WriteFile("README.md", bdata, 0644)
    if err != nil { panic(err) }
}

func writeHtmlIndex(t *template.Template, packages []Package) {
    f, err := os.Create("go.html")
    if err != nil { panic(err) }
    defer f.Close()

    type Head struct {
        Title string
        Desc string
    }

    type Body struct {
        Title string
    }

    type Data struct {
        Year int
        Head Head
        Body Body
        Packages []Package
    }

    data := Data{
        Year: time.Now().UTC().Year(),
        Head: Head{
            Title: "Tawesoft® Go monorepo",
            Desc:  "A monorepo for small Go (golang) modules maintained by Tawesoft®",
        },
        Body: Body{
            Title: "tawesoft.co.uk/go",
        },
        Packages: packages,
    }

    err = t.Execute(f, data)
    if err != nil { panic(fmt.Sprintf("template error: %v", err)) }
}

func writeLicenseIndex(packages []Package) {
    f, err := os.Create("LICENSE.txt")
    if err != nil { panic(err) }
    defer f.Close()

    divider := "\n\n--------------------------------------------------------------------------------\n\n"

    for i, p := range packages {
        suffix := "\n"
        if i < len(packages) - 1 {
            suffix = divider
        }

        _, err := f.WriteString(p.License + suffix)
        if err != nil { panic(err) }
    }
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

func fmtMarkdownPackageTable(p Package) string {
    template := strings.TrimSpace(`
|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_%[1]s] ∙ [docs][docs_%[1]s] ∙ [src][src_%[1]s] | [%[2]s][copy_%[1]s] | %[3]s |

[home_%[1]s]: https://tawesoft.co.uk/go/%[1]s
[src_%[1]s]:  https://github.com/tawesoft/go/tree/master/%[1]s
[docs_%[1]s]: https://www.tawesoft.co.uk/go/doc/%[1]s
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
