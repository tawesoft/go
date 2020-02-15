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
[copy_{slug}]: https://github.com/tawesoft/go/tree/master/{id}/COPYING.md
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
        stable=False, # This module API is not yet stable and may be subject to occasional breaking changes.
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
        with open("%s/COPYING.md" % module.id, "w") as fp:
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


def run(fn):
    print("> %s: %s" % (fn.__name__, fn.__doc__))
    fn()


def run_many(*args):
    for i in args:
        run(i)

run_many(
    make_base_go,
    make_module_go,
    make_base_license_txt,
    make_module_license_txt,
    make_module_copying_md,
    make_base_readme_md,
    make_module_readme_md,
    make_module_go,
)

