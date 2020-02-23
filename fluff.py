# fluff.py for tawesoft.co.uk/go - generate human-readable fluff
#
# Copyright © 2020 Tawesoft Ltd <open-source@tawesoft.co.uk>
# Copyright © 2020 Ben Golightly <ben@tawesoft.co.uk>
#
# usage: python3 fluff.py
#
# This script uses a central catalog of metadata to generate human-readable
# "fluff" like HTML pages, Markdown documents, package-level docstrings, the
# top-level monorepo go file, license files, etc.
#
# Copying and distribution of this file, with or without modification, are
# permitted in any medium without royalty provided the copyright notice and
# this notice are preserved. This file is offered as-is, without any
# warranty.


import datetime
import textwrap
from dataclasses import dataclass, field
from typing import List
import os.path


@dataclass
class ModuleDesc:
    short: str
    long: str

    def summary(self):
        return self.long.partition("\n\n")[0].strip()


@dataclass
class ModuleLicense:
    id: str
    name: str
    text: str
    commentary: str = ""


@dataclass
class Module:
    id:           str
    desc:         ModuleDesc
    license:      ModuleLicense
    copyright:    str
    example:      str = ""
    exampleFiles: List[str] = field(default_factory=list)
    seeAlso:      List[str] = field(default_factory=list)
    stable:       bool = True

    def name(self):
        return os.path.basename(self.id)

    def slug(self, replacement):
        return self.id.replace('/', replacement)

    def link_table_markdown(self, unique_urls=True):
        return """
|  Links  | License | Stable? | 
|:-------:|:-------:|:-------:| 
| [home][home_{slug}] ∙ [docs][docs_{slug}] ∙ [src][src_{slug}] | [{license}][copy_{slug}] | {stable} |

[home_{slug}]: https://tawesoft.co.uk/go/{id}
[src_{slug}]:  https://github.com/tawesoft/go/tree/master/{id}
[docs_{slug}]: https://godoc.org/tawesoft.co.uk/go/{id}
[copy_{slug}]: https://github.com/tawesoft/go/tree/master/{id}/_COPYING.md
""".format(
            id=self.id,
            slug=(self.slug("_") if unique_urls is True else ""),
            license=self.license.id,
            stable=("✔ yes" if self.stable else "✘ **no**")
        ).strip()


trademarkOpenGL = """OpenGL® and the oval logo are trademarks or registered trademarks of Hewlett Packard Enterprise in
the United States and/or other countries worldwide."""


licenseMIT0 = ModuleLicense("MIT-0", "MIT No Attribution", """
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction,  including without limitation the rights
to use,  copy, modify,  merge,  publish, distribute, sublicense,  and/or sell
copies  of  the  Software,  and  to  permit persons  to whom  the Software is
furnished to do so.

THE SOFTWARE IS PROVIDED  "AS IS",  WITHOUT WARRANTY OF ANY KIND,  EXPRESS OR
IMPLIED,  INCLUDING  BUT  NOT LIMITED TO THE WARRANTIES  OF  MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE  AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
AUTHORS  OR COPYRIGHT HOLDERS  BE LIABLE  FOR ANY  CLAIM,  DAMAGES  OR  OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
""", commentary="""
To be clear: you are free to use this project in binary distributions without
attribution. In the interest of intellectual honesty please attribute source
code extracts or derivatives; it is sufficient to leave a URL to this project
in a source code comment.
""")


licenseMIT = ModuleLicense("MIT", "MIT License", """
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction,  including without limitation the rights
to use,  copy, modify,  merge,  publish, distribute, sublicense,  and/or sell
copies  of  the  Software,  and  to  permit persons  to whom  the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED  "AS IS",  WITHOUT WARRANTY OF ANY KIND,  EXPRESS OR
IMPLIED,  INCLUDING  BUT  NOT LIMITED TO THE WARRANTIES  OF  MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE  AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
AUTHORS  OR COPYRIGHT HOLDERS  BE LIABLE  FOR ANY  CLAIM,  DAMAGES  OR  OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
""",
commentary="""
You may find this FAQ useful: https://www.tawesoft.co.uk/kb/article/mit-license-faq
""")


def licenseBSD3(author: str, commentary: str):
    return ModuleLicense("BSD-3-Clause", """BSD 3-Clause "New" or "Revised" License""", """
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of %s nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
""" % author, commentary=commentary)


licenseGo = licenseBSD3("Google Inc.", commentary="""
This module is based on source code originally by The Go Authors.

The Go Authors and Google Inc. are not affiliated with this project in any way.
""")


def make_base_go():
    """write a dummy tawesoft.go file in the base directory that imports everything as a home for documentation"""
    with open("tawesoft.go", "w") as fp:
        imports = "\n".join(["    _ \"tawesoft.co.uk/go/%s\"" % i.id for i in catalog])
        fp.write("""
/*
A monorepo for small Go modules maintained by Tawesoft®

This is permissively-licensed open source software but exact licenses may vary between modules.

For license information, documentation, source code, support, links, etc. please see
https://tawesoft.co.uk/go
*/
package tawesoft

import (
%s        
)
""".strip() % imports)


def make_module_go():
    """write a module .go file in each module directory for the docstring"""
    for i in catalog:
        docstring = i.desc.long.strip()

        if i.seeAlso:
            docstring += "\n\n"+"\n".join(["See also: %s (https://tawesoft.co.uk/go/%s)" % (x, x) for x in i.seeAlso])

        if i.example:
            docstring += "\n\nExample:\n\n"+textwrap.indent(i.example.strip(), "    ", predicate=lambda _: True)

        with open("%s/%s.go" % (i.id, i.name()), "w") as fp:
            fp.write("""
/*
{docstring}

For license information, documentation, source code, support, links, etc. please see
https://tawesoft.co.uk/go/{id}

This module is part of https://tawesoft.co.uk/go
*/
package {name} // import "tawesoft.co.uk/go/{id}"

// SPDX-License-Identifier: {license.id}

// Code generated by (tawesoft.co.uk/go/) fluff.py: DO NOT EDIT.

""".format(id=i.id, name=i.name(), docstring = docstring, license=i.license).strip())


def make_base_license_txt():
    """write a combined license.txt file in the base directory"""
    with open("LICENSE.txt", "w") as fp:
        divider = "\n%s\n\n" % ("-"*80)
        fp.write("tawesoft.co.uk/go\n")
        for module in catalog:
            fp.write("%stawesoft.co.uk/go/%s\n\n%s\n\n%s\n" % \
                 (divider, module.id, module.copyright.strip(), module.license.text.strip()))


def make_module_license_txt():
    """write a license.txt file in each module directory"""
    for module in catalog:
        with open("%s/LICENSE.txt" % module.id, "w") as fp:
            fp.write("tawesoft.co.uk/go/%s\n\n%s\n\n%s" % \
                 (module.id, module.copyright.strip(), module.license.text.strip()))


def make_module_copying_md():
    """write a copying markdown file with optional commentary in each module directory"""
    for module in catalog:
        with open("%s/_COPYING.md" % module.id, "w") as fp:
            content = "# License\n\n```\ntawesoft.co.uk/go/%s\n\n%s\n\n%s\n```" % \
                 (module.id, module.copyright.strip(), module.license.text.strip())

            if module.license.commentary:
                prefix = "This module is released under the %s (SPDX: %s)." % (module.license.name, module.license.id)
                suffix = "This commentary is not part of the license."
                content += "\n\n## Commentary:\n\n%s\n\n%s\n\n%s" % (prefix, module.license.commentary.strip(), suffix)

            fp.write(content)


def make_base_readme_md():
    """write an index markdown README in the base directory"""

    template = """
[![Tawesoft](https://www.tawesoft.co.uk/media/0/logo-240r.png)](https://tawesoft.co.uk/go)
================================================================================

A monorepo for small Go modules maintained by [Tawesoft®](https://www.tawesoft.co.uk/go)

This is permissively-licensed open source software but exact licenses may vary between modules.

Download
--------

```shell script
go get -u tawesoft.co.uk/go
```

Contents
--------

{contents}

Links
-----

* Home: [tawesoft.co.uk/go](https://tawesoft.co.uk/go)
* Docs hub: [godoc.org/tawesoft.co.uk/go](https://godoc.org/tawesoft.co.uk/go)
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
"""

    fmt = """
### {id}: {desc_short}

`import "tawesoft.co.uk/go/{id}"`

{desc_summary}

{link_table}

---
""".strip()

    contents = "\n\n".join([fmt.format(
        id=i.id,
        slug=i.slug("_"),
        desc_short=i.desc.short,
        desc_summary=i.desc.summary(),
        link_table=i.link_table_markdown().strip(),
    ) for i in catalog])

    with open("README.md", "w") as fp:
        fp.write(template.format(contents=contents).strip())


def make_module_readme_md():
    """write an index markdown README in each module directory"""
    for module in catalog:
        with open("%s/README.md" % module.id, "w") as fp:
            content = """
# {module.id} - {module.desc.short}

## About

{long_desc}

{link_table}

## Download

```shell script
go get -u tawesoft.co.uk/go
```

## Import

```
import tawesoft.co.uk/go/{module.id}
```

""".format(
                module=module,
                long_desc=module.desc.long.strip(),
                link_table=module.link_table_markdown(False),
            ).strip()

            if module.seeAlso:
                content += "\n\n## See Also:\n\n"+"\n".join(["* %s (https://tawesoft.co.uk/go/%s)" % (x, x) for x in module.seeAlso])

            if module.example:
                content += "\n\n## Example:\n\n```go\n%s\n```" % module.example.strip()

            fp.write(content)


def make_base_html():
    template="""
<!doctype html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <title>Tawesoft Go monorepo</title>
    <meta name="description" content="A monorepo for small Go (golang) modules maintained by Tawesoft®" />
    <meta name="go-import" content="tawesoft.co.uk/go git https://github.com/tawesoft/go">
    <style>
        html, body                   {{ color: #323232; background-color: #FFF; font-family: sans-serif; padding: 0; margin: 0; }}
        body                         {{ min-width: 800px; padding: 1em; }}
        h2, h3                       {{ margin-top: 2em; }}
        a                            {{ color: #b63355; text-decoration: none; }}
        a:hover                      {{ text-decoration: underline; }}
        table                        {{ width: 100%; border-spacing: 0 0; }}
        table, thead, tbody, tr      {{ border: 0; border-collapse: separate; text-align: left; }}
        th                           {{ border-bottom: 1px solid #000; }}
        th, td                       {{ padding: 0.5em; }}
        li                           {{ margin: 1em 0; }}
        
        #container                   {{ max-width: 800px; margin: 0 auto; }}
        #footer                      {{ max-width: 1200px; margin: 0 auto; text-align: right; margin-top: 2em; }}
        #footer li                   {{ margin: 0; }}
        
        table.pkg td.stable          {{ text-align: center; }}
        table.pkg td.stable > span   {{ color: #b63355; }}
        table.pkg td.links           {{ padding-left: 0; }}
        table.pkg td.links > a       {{ display: inline-block; background-color: #FEE; padding: 0.5em 1em; margin: 0.25em; text-decoration: none; }}
        table.pkg td.links > a:first-child {{ margin-left: 0; }}
        table.pkg td.links > a:hover {{ background-color: #b63355; color: #FFF; }}
        table.pkg tbody > tr         {{ height: 4em; line-height: 1.5em; }}
        
div#footer {{
  padding-top: 64px;
  color: #686062; }}
  div#footer > div.block {{
    text-align: right; }}
  div#footer a {{
    color: #b73154; }}
  div#footer a:hover {{
    text-decoration: underline; }}
  div#footer ul {{
    list-style-type: none; }}
  div#footer div.smedia {{
    float: left;
    margin-bottom: 1em; }}
    div#footer div.smedia a {{
      display: inline-block;
      width: 64px;
      height: 64px;
      margin-right: 16px;
      background-image: url("smedia.png");
      background-repeat: no-repeat; }}
      div#footer div.smedia a:hover {{
        text-decoration: none !important; }}
      div#footer div.smedia a.twitter {{
        width: 42px;
        background-position: 0px 14px; }}
        div#footer div.smedia a.twitter:hover {{
          background-position: -43px 14px; }}
      div#footer div.smedia a.linkedin {{
        background-position: -86px 0; }}
        div#footer div.smedia a.linkedin:hover {{
          background-position: -150px 0; }}
      div#footer div.smedia a.facebook {{
        background-position: -86px -64px; }}
        div#footer div.smedia a.facebook:hover {{
          background-position: -150px -64px; }}
      div#footer div.smedia a.github {{
        background-position: -86px -128px; }}
        div#footer div.smedia a.github:hover {{
          background-position: -150px -128px; }}
    </style>
</head>
<body>

<div id="container">

<div style="float: right; margin: 1em;">
<a href="https://www.tawesoft.co.uk/products/open-source-software"><img src="https://www.tawesoft.co.uk/media/0/logo-240r.png" alt="Tawesoft Logo" /></a>
</div>

<h1>tawesoft.co.uk/go</h1>

<p>A monorepo for small Go modules maintained by <a href="https://www.tawesoft.co.uk/">Tawesoft<sup>&reg;</sup></a></p>

<p>This is permissively-licensed open source software but exact licenses may vary between modules.</p>

<h2>Download</h2>

<pre>go get -u tawesoft.co.uk/go</pre>

<h2>Packages</h2>

<table class="pkg">
    <thead>
        <tr>
            <th style="width: 50%;">Module</th>
            <th style="width: 20%;">Links</th>
            <th style="width: 10%;">Stable?</th>
            <th style="width: 20%;">License</th>
        </tr>
    </thead>
    <tbody>
{modules}
    </tbody>
</table>

<h2>Links:</h2>

<ul>
    <li>Home: <a href="https://tawesoft.co.uk/go">tawesoft.co.uk/go</a></li>
    <li>Docs hub: <a href="https://godoc.org/tawesoft.co.uk/go">godoc.org/tawesoft.co.uk/go</a></li>
    <li>Repository: <a href="https://github.com/tawesoft/go">github.com/tawesoft/go</a></li>
    <li>Or <a href="https://pkg.go.dev/search?q=tawesoft">search "tawesoft"</a> on <a href="https://go.dev/">go.dev</a></li>
</ul>

<h2>Support:</h2>

<h3>Free and Community Support</h3>

<ul>
    <li><a href="https://github.com/tawesoft/go/issues">GitHub issues</a></li>
    <li>Email <a href="mailto:open-source@tawesoft.co.uk">open-source@tawesoft.co.uk</a> (feedback welcomed, but support is "best effort")</li>
</ul>

<h3>Commercial Support</h3>

<p>Open source software from Tawesoft® is backed by commercial support options.</p>

<p>Email <a href="mailto:open-source@tawesoft.co.uk">open-source@tawesoft.co.uk</a>  or visit
<a href="https://www.tawesoft.co.uk/products/open-source-software">tawesoft.co.uk/products/open-source-software</a> to learn more.</p>

</div>

<div id="footer" class="toplevel block">
    <div class="smedia"
        ><a rel="external nofollow" class="facebook" href="https://www.facebook.com/tawesoft"         title="Find Tawesoft on Facebook">&nbsp;</a
        ><a rel="external nofollow" class="linkedin" href="https://www.linkedin.com/company/10149194" title="Connect with Tawesoft on LinkedIn">&nbsp;</a
        ><a rel="external nofollow" class="github"   href="https://github.com/tawesoft"               title="Code with Tawesoft on GitHub">&nbsp;</a
        ><a rel="external nofollow" class="twitter"  href="https://twitter.com/tawesoft"              title="Follow @tawesoft on Twitter">&nbsp;</a
    ></div>
    <div class="block">
        <p lang="en-gb" class="en">Tawesoft &reg; is a trading style and registered trade mark of <strong>Tawesoft Ltd</strong><br />
        Registered in England and Wales, Company Number 9735741</p>
        <p lang="cy" class="cy">Tawesoft &reg; yn enw masnach a nod masnach cofrestredig o <strong>Tawesoft Cyf</strong><br />
        Cofrestrwyd yng Nghymru a Lloegr, Rhif Cwmni 9735741</p>
        <p>Copyright / <span lang="cy">Hawlfraint</span> &copy; <a href="https://www.tawesoft.co.uk/">Tawesoft Ltd</a> 2019 - 2020</p>
        
        <ul>
            <li><a href="https://www.tawesoft.co.uk/about/tawesoft">About Tawesoft / <span lang="cy">Pwy yw Tawesoft?</span></a></li>
            <li><a href="https://www.tawesoft.co.uk/about/contact">Contact us / <span lang="cy">Cysylltu â ni</span></a></li>
            <li><a href="https://www.tawesoft.co.uk/about/privacy">Privacy / <span lang="cy">Preifatrwydd</span></a></li>
        </ul>
    </div>
</div>

</body>
</html>
""".strip()

    row="""
        <tr>
            <td><b>{module.id}</b><br />{module.desc.short}<br /></td>
            <td class="links">
                <a href="https://godoc.org/tawesoft.co.uk/go/{module.id}">docs</a>
                <a href="https://github.com/tawesoft/go/tree/master/{module.id}">src</a>
            </td>
            <td class="stable">{stable}</td>
            <td>
                <a href="https://github.com/tawesoft/go/tree/master/{module.id}/_COPYING.md">{module.license.id}</a>
            </td>
        </tr>
""".strip()

    modules="\n".join([row.format(
        module=i,
        stable="<span>✔</span> yes" if i.stable else "<span>✘</span> <b>no</b>"
    ) for i in catalog])

    with open("go.html", "w") as fp:
        fp.write(template.format(modules=modules))



def run(fn):
    print("> %s: %s" % (fn.__name__, fn.__doc__))
    fn()


def run_many(*args):
    for i in args:
        run(i)


catalog = [

    Module(
        id="dialog",
        desc=ModuleDesc(
            short="simple cross-platform messagebox",
            long="""
Package dialog implements simple cross platform native MessageBox/Alert dialogs for Go.

Currently only Windows and Linux targets are supported.

On Linux, uses (in order of preference) `zenity`, `xmessage`, or stdio.
"""
        ),
        license=licenseMIT0,
        copyright="""
Copyright © 2019 - 2020 Ben Golightly <ben@tawesoft.co.uk>
Copyright © 2019 - 2020 Tawesoft Ltd <opensource@tawesoft.co.uk>
        """,
        example="""
package main

import "tawesoft.co.uk/go/dialog"

func main() {
    dialog.Alert("Message")
    dialog.Alert("There are %d lights", 4)
}
""",
    ),

    Module(
        id="email",
        desc=ModuleDesc(
            short="format multipart RFC 2045 email",
            long="""
Package email implements the formatting of multipart RFC 2045 e-mail messages,
including headers, attachments, HTML email, and plain text.
""",
        ),
        license=licenseMIT,
        copyright="""
Copyright © 2020 Tawesoft Ltd <open-source@tawesoft.co.uk>
Copyright © 2020 Ben Golightly <ben@tawesoft.co.uk>
""",
        example="""
package main

import (
    "net/mail"
    "os"
    
    "tawesoft.co.uk/go/email"
)

func main() {
    var eml = email.Message{
        From:  mail.Address{"Alan Turing", "turing.alan@example.org"},
        To:  []mail.Address{{"Grace Hopper", "amazing.grace@example.net"}},
        Bcc: []mail.Address{{"BCC1", "bcc1@example.net"}, {"BCC2", "bbc2@example.net"}},
        Subject: "Computer Science is Cool! ❤",
        Text: `This is a test email!`,
        Html: `<!DOCTYPE html><html lang="en"><body><p>This is a test email!</p></body></html>`,
        Attachments: []*email.Attachment{
            //email.FileAttachment("Entscheidungsproblem.pdf"),
            //email.FileAttachment("funny-cat-meme.png"),
        },
        Headers: mail.Header{
            "X-Category": []string{"newsletter", "marketing"},
        },
    }
    
    var err = eml.Print(os.Stdout)
    if err != nil { panic(err) }
}
""",
        exampleFiles=[
            "examples/example1/example1.go",
        ],
    ),

    Module(
        id="glcaps",
        desc=ModuleDesc(
            short="read and check OpenGL capabilities",
            long="""
Package glcaps provides a nice interface to declare OpenGL capabilities you care about, including minimum required
extensions or capabilities. Glcaps has no dependencies and is agnostic to the exact OpenGL binding used.

"""+trademarkOpenGL,
        ),
        license=licenseMIT,
        copyright="""
Copyright © 2020 Tawesoft Ltd <open-source@tawesoft.co.uk>
Copyright © 2020 Ben Golightly <ben@tawesoft.co.uk>
""",
        exampleFiles=[
            "example/example.go",
        ],
    ),

    Module(
        id="humanize",
        desc=ModuleDesc(
            short="lightweight human-readable numbers",
            long="""
Package humanize implements lightweight human-readable numbers for Go.

Why use this package?

Compare tawesoft/go/humanize vs dustin/go-humanize: Tawesoft's parses input about 5 times faster with fewer memory
allocations, and formats output about 25% quicker than dustin's. But dustin's is so far more complete, older, has a
stable API, and has been tested by more people.
""",
        ),
        license=licenseMIT0,
        copyright="""
Copyright © 2020 Tawesoft Ltd <open-source@tawesoft.co.uk>
Copyright © 2020 Ben Golightly <ben@tawesoft.co.uk>
""",
        stable=False, # This API is not yet stable and may be subject to occasional breaking changes.
    ),

    Module(
        id="operator",
        desc=ModuleDesc(
            short="operators as functions",
            long="""
Package operator implements logical, arithmetic, bitwise and comparison
operators as functions (like the Python operator module). Includes unary,
binary, and nary functions with overflow checked variants.
""",
        ),
        license=licenseMIT0,
        copyright="""
Copyright © 2020 Tawesoft Ltd <open-source@tawesoft.co.uk>
Copyright © 2020 Ben Golightly <ben@tawesoft.co.uk>
    """,
        example="""
package main

import (
    "fmt"
    "tawesoft.co.uk/go/operator"
)

func foo(op func(int, int) int, a int, b int) int {
    return op(a, b)
}

func fooChecked(op func(int8, int8) (int8, error), a int8, b int8) (int8, error) {
    return op(a, b)
}

func main() {
    fmt.Println(foo(operator.Int.Binary.Add, 5, 3))
    fmt.Println(foo(operator.Int.Binary.Sub, 5, 3))
    
    var result, err = fooChecked(operator.Int8Checked.Binary.Add, 126, 2) // max int8 is 127!
    if err != nil {
        fmt.Printf("error: %v (expected!)\\n", err)
    } else {
        fmt.Println(result)
    }
}
""",
        exampleFiles=[
            "examples/calculator/calculator.go",
        ],
    ),

    Module(
        id="xff",
        desc=ModuleDesc(
            short="DirectX (.x) file format decoder",
            long="""
Package xff implements a decoder for the DirectX (.x) file format in
Go (Golang).

THIS IS A PREVIEW RELEASE. A few templates are missing and the returned data
is not decoded from bytes to types properly. The testdata is missing. It needs
a refactor. This will be fixed early next week (beginning Feb 24th 2020).

This parser is (or aims to be) a complete implementation. It supports user-defined
templates and should be able to parse *every* well-formed DirectX (.x) file.

There are a few features not yet implemented - not because they are difficult
to implement, but because I can't find any real-world examples using that
feature to test against. Things like object referencing by UIID instead of
name, or the binary encoded DirectX .x file format, arrays of strings, or
multidimensional arrays. Send in a sample if you find something that doesn't
work (and should) and it'll get fixed quickly!

This is free and open source software. It took about five days to write - so
if you find this module useful enough that you use it in a commercial product then
there is a polite expectation that you donate a sum to the author equal to the
retail cost of one copy of your product. To do so, make a payment at
https://www.paypal.me/TawesoftLtd and forward a copy of your receipt to
"xff-thanks@tawesoft.co.uk". Alternatively you can purchase commercial support
through open-source@tawesoft.co.uk. If you're making a computer game I'd
appreciate one free copy, too. This is, of course, entirely optional.

DirectX is a registered trademark of Microsoft Corporation in the United States
and/or other countries.
""",
        ),
        license=licenseMIT,
        copyright="""
Copyright © 2020 Tawesoft Ltd <open-source@tawesoft.co.uk>
Copyright © 2020 Ben Golightly <ben@tawesoft.co.uk>
""",
        stable=False, # This API is not yet stable and may be subject to occasional breaking changes.
    ),

    Module(
        id="ximage",
        desc=ModuleDesc(
            short="extended image types",
            long="""
Package ximage implements Red, RG, and RGB images matching the core
image interface.

Note that there are good reasons these image types aren't in the core image
package. The native image types may have optimized fast-paths for many use
cases.

This package is a tradeoff of these optimizations against lower memory
usage. This package is intended to be used in computer graphics (e.g.
OpenGL) where images are uploaded to the GPU in a specific format (such as
GL_R, GL_RG, or GL_RGB) and we don't care too much about the performance of
native Go image manipulation.

"""+trademarkOpenGL,
        ),
        license=licenseGo,
        copyright="Copyright (c) 2009 The Go Authors. All rights reserved.",
        seeAlso=[
            "ximage/xcolor",
        ],
    ),

    Module(
        id="ximage/xcolor",
        desc=ModuleDesc(
            short="extended color types",
            long="""
Package xcolor implements Red, RedGreen, and RGB color models matching the core
image/color interface.

Note that there are good reasons these color types aren't in the core
image.color package. The native color types may have optimized fast-paths
for many use cases.

This package is a tradeoff of these optimizations against lower memory
usage. This package is intended to be used in computer graphics (e.g.
OpenGL) where images are uploaded to the GPU in a specific format (such as
GL_R, GL_RG, or GL_RGB) and we don't care about the performance of native
Go image manipulation.

"""+trademarkOpenGL,
        ),
        license=licenseGo,
        copyright="Copyright (c) 2009 The Go Authors. All rights reserved.",
        seeAlso=list([
            "ximage",
        ]),
    ),

]

catalog.sort(key=lambda x: x.id)


run_many(
    make_base_html,
    make_base_go,
    make_module_go,
    make_base_license_txt,
    make_module_license_txt,
    make_module_copying_md,
    make_base_readme_md,
    make_module_readme_md,
    make_module_go,
)

