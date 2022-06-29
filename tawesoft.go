/*
A monorepo for small Go modules maintained by Tawesoft®

This is permissively-licensed open source software but exact licenses may vary between modules.

For license information, documentation, source code, support, links, etc. please see
https://www.tawesoft.co.uk/go

FROZEN - PLEASE MIGRATE

These packages are moving to https://github.com/tawesoft/golib.

This is to increase security against possible supply chain attacks such as our domain name expiring in the future and being registered by someone else.

Please migrate to https://github.com/tawesoft/golib (when available) instead.

Most programs relying on a package in this monorepo, such as the dialog or lxstrconv packages, will continue to work for the foreseeable future.

Rarely used packages have been hidden for now - they are in the git commit history at https://github.com/tawesoft/go if you need to resurrect one.

*/
package tawesoft

import (
    _ "tawesoft.co.uk/go/dialog"
    _ "tawesoft.co.uk/go/glcaps"
    _ "tawesoft.co.uk/go/humanizex"
    _ "tawesoft.co.uk/go/lxstrconv"
    _ "tawesoft.co.uk/go/operator"
)
