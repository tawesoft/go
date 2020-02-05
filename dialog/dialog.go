/*
Package dialog implements simple cross platform native MessageBox/Alert
dialogs for Go.

Currently only Windows and Linux targets are supported.

For source code see https://github.com/tawesoft/dialog

For links, news, etc. see https://www.tawesoft.co.uk/go/dialog
*/
package dialog // import "tawesoft.co.uk/go/dialog"

// Alert displays a modal message box with message. The message string can
// be a printf-style format string for an optional sequence of additional
// arguments of any type.
func Alert(message string, args...interface{}) {
    platformAlert("Alert", message, args...)
}
