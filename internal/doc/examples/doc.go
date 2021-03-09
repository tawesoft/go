// Generates HTML examples.
//
// NOTE: expects to be run from the same directory as its source file
package main

import (
    "bytes"
    "fmt"
    "html/template"
    "os"
    "path"
    "sort"
    "strings"

    "github.com/alecthomas/chroma"
    "github.com/alecthomas/chroma/formatters/html"
    "github.com/alecthomas/chroma/lexers"
    "github.com/alecthomas/chroma/styles"
)

func main() {
    const theme = "lovelace"

    // monorepo root
    os.Chdir("../../../")

    // err := quick.Highlight(os.Stdout, someSourceCode, "go", "html", theme)

    packageNames := os.Args[1:]
    sort.Strings(packageNames)

    t, err := template.New("page").Parse(
        func() string {
            data, err := os.ReadFile("internal/doc/examples/template.html")
            if err != nil { panic(err) }
            return string(data)
        }(),
    )
    if err != nil { panic(err) }

    // init chroma code highlighter
    lexer := lexers.Get("go")
    lexer = chroma.Coalesce(lexer)
    style := styles.Get(theme)
    if style == nil { panic("missing theme") }
    formatter := html.New(html.WithClasses(true))
    if formatter == nil { panic("missing formatter") }

    type Head struct {
        Title string
        Desc string
        CSS template.CSS
    }

    type Body struct {
        Title template.HTML
        Html  template.HTML
    }

    type Data struct {
        Head Head
        Body Body
    }

    css := func() template.CSS {
        buf := &bytes.Buffer{}
        err := formatter.WriteCSS(buf, style)
        if err != nil { panic(err) }
        return template.CSS(buf.String())
    }()

    for _, name := range packageNames {
        if strings.HasPrefix(name, "legacy") { continue }

        examples, err := os.ReadDir(path.Join(name, "examples"))
        if err != nil {
            if os.IsNotExist(err) { continue }
            panic(err)
        }

        for _, example := range examples {
            if strings.HasPrefix(example.Name(), ".")   { continue }
            if strings.HasPrefix(example.Name(), "_")   { continue }
            if strings.HasPrefix(example.Name(), "dev") { continue }

            src := fmt.Sprintf("%[1]s/examples/%[2]s/%[2]s.go", name, example.Name())
            dest := fmt.Sprintf("doc/%[1]s/examples/%[2]s/index.html", name, example.Name())

            // highlight
            contents, err := os.ReadFile(src)
            if err != nil { panic(err) }
            iterator, err := lexer.Tokenise(nil, string(contents))

            func() {
                w, err := os.Create(dest)
                if err != nil { panic(err) }
                defer w.Close()

                buf := &bytes.Buffer{}

                err = formatter.Format(buf, style, iterator)
                if err != nil { panic(err) }

                data := Data{
                    Head: Head{
                        CSS: css,
                    },
                    Body: Body{
                        Title: template.HTML(fmt.Sprintf(
                            `<a href="https://www.tawesoft.co.uk/go/doc/%[1]s">tawesoft.co.uk/go/%[1]s</a> example: %[2]s.go`,
                        name, example.Name())),
                        Html: template.HTML(buf.String()),
                    },
                }

                err = t.Execute(w, data)
                if err != nil { panic(err) }
            }()
        }
    }

    // contents, err := ioutil.ReadAll(r)
    //

}
