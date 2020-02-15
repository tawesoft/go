package dialog

// Alert displays a modal message box with message. The message string can
// be a printf-style format string for an optional sequence of additional
// arguments of any type.
func Alert(message string, args...interface{}) {
    platformAlert("Alert", message, args...)
}
