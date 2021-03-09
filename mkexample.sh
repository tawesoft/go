#!/usr/bin/env bash

godocport=6061

pkgs="\
    dialog \
    drop \
    email \
    glcaps \
    legacy/email \
    loader \
    lxstrconv \
    operator \
    queue \
    router \
    sqlp \
    sqlp/sqlite3 \
    variadic \
    ximage \
    ximage/xcolor
"


# build examples
cd internal/doc/examples
go run doc.go $pkgs
cd ../../../


