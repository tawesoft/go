godocport=6061

pkgs="\
    dialog \
    drop \
    email \
    glcaps \
    loader \
    lxstrconv \
    operator \
    queue \
    router \
    sqlp \
    sqlp/sqlite3 \
    variadic \
    ximage \
    ximage/xcolor \
    legacy/email \
"

# build go.html, READMEs, etc.
go run internal/doc/doc.go $pkgs

eval `go env`

# add symlinks to default godoc template files (everything except "*.go")
prefix=$GOPATH/src/golang.org/x/tools/godoc/static/
templates=$(find $prefix -type f \( ! -iname '*.go' \))

for i in $templates
do
    # strip prefix
    file=${i/#$prefix}
    dest="internal/godoc/$file"

    # make directories as needed
    dir=$(dirname $file)
    mkdir -p "internal/godoc/$dir"

    # symlink if doesn't already exist
    if [ ! -e $dest ]; then
        echo "$file => $dest"
        ln -s "$prefix$file" "$dest"
    fi
done

# start godoc server with customised templates
godoc -http=:$godocport -templates="internal/godoc" -notes="BUG|TODO|FIXME" &
pid=$!
sleep 1

# extract
# (https://github.com/golang/go/issues/2381#issuecomment-66059484)
mkdir -p doc
rm -r doc
wget -q -nv -e robots=off \
    --cut-dirs 3 \
    -P doc/ \
    -A go,html,css,js,gif,png \
    -r -nH -E -p -k \
    --restrict-file-names=windows \
    -I /src/tawesoft.co.uk/go,/pkg/builtin,/pkg/tawesoft.co.uk/go,/lib \
    http://localhost:$godocport/

# --restrict-file-names=windows - dont use questionarms
# -R reject
# -A accept
# -P dest dir
# -r  : download recursive
# -np : don't ascend to the parent directory
# -nd : no directories
# -nc : no clobber
# -E  : add extension .html to html files (if they don't have)
# -p  : download all necessary files for each page (css, js, images)
# -k  : convert links to relative
# -nv : no verbose
# -nH : Disable generation of host-prefixed directories.

# close server
kill -1 $pid


