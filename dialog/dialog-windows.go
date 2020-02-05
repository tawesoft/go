// +build windows

package dialog // import "tawesoft.co.uk/go/dialog"

import (
    "fmt"
    "golang.org/x/sys/windows"
)


const cp_utf8 uint32 = 65001;


func toWideChar(input string) []uint16 {

    var buf []uint16
    var required, numchars int32
    var err error

    if len(input) == 0 { goto fail }

    // func MultiByteToWideChar(codePage uint32, dwFlags uint32, str *byte,
    // nstr int32, wchar *uint16, nwchar int32) (nwrite int32, err error)
    required, err = windows.MultiByteToWideChar(
        cp_utf8,
        0,
        &([]byte(input))[0],
        int32(len(input)),
        nil,
        0,
    )

    if (err != nil) { goto fail }

    buf = make([]uint16, required+1)

    numchars, err = windows.MultiByteToWideChar(
        cp_utf8,
        0,
        &([]byte(input))[0],
        int32(len(input)),
        &([]uint16(buf))[0],
        required,
    )

    if (err != nil) { goto fail2 }
    if (len(buf) != int(numchars) + 1) { goto fail2 }

    return buf

    fail2:
        buf = nil
    fail:
        return make([]uint16, 1) // empty null-terminated string
}

func platformAlert(title string, message string, args...interface{}) {
    var msg string

    if len(args) > 0 {
        msg = fmt.Sprintf(message, args...)
    } else {
        msg = message
    }

    var wtitle = toWideChar(title)
    var wmessage = toWideChar(msg)
    var flags uint32 = windows.MB_OK |
        windows.MB_ICONEXCLAMATION |
        windows.MB_SETFOREGROUND |
        windows.MB_TOPMOST

    windows.MessageBox(0, &wmessage[0], &wtitle[0], flags)

    wtitle = nil
    wmessage = nil
}
