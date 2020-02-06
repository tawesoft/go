#!/bin/sh

licenses="
dialog/LICENSE.txt
humanize/LICENSE.txt
ximage/LICENSE.txt
ximage/xcolor/LICENSE.txt
"

dest="LICENSE.txt"

printf "tawesoft.co.uk/go\n" > "$dest"
for license in $licenses; do
    printf "\n--------------------------------------------------------------------------------\n\n" >> "$dest"
    cat "$license" >> "$dest"
done
