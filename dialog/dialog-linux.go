// +build linux

package dialog // import "tawesoft.co.uk/go/dialog"

import (
    "fmt"
    "os"
    "os/exec"
)

type provider struct {
    command string
    alert func(self *provider, title string, message string, args...interface{}) bool
}

var xmessageProvider = provider{
    command: "xmessage",
    alert: func(self *provider, title string, message string, args...interface{}) bool {
        var buf string = wrap(fmt.Sprintf(title+": "+message+"\n", args...), 60)

        var err = exec.Command(self.command, "-center", buf).Run()
        if err != nil { fmt.Printf("system error: %v\n", err) }
        return err == nil
    },
}

var zenityProvider = provider{
    command: "zenity",
    alert: func(self *provider, title string, message string, args...interface{}) bool {
        var buf string = fmt.Sprintf(message, args...)

        var err = exec.Command(self.command,
            "--info", "--no-markup",
            "--title", title,
            "--window-icon", "info",
            "--width=400",
            "--text="+buf,
        ).Run()

        if err != nil {
            var _, ok = err.(*exec.ExitError)
            if ok {
                // fine - we don't care about the return code
                return true
            }
            // other type of error
            fmt.Printf("system error: %v\n", err)
        }
        return err == nil
    },
}

var stdioProvider = provider{
    command: "",
    alert: func(_self *provider, title string, message string, args...interface{}) bool {
        fmt.Fprintf(os.Stderr, "\n===[%s]===\n\n", title)
        fmt.Fprintf(os.Stderr, message, args...)
        fmt.Fprintf(os.Stdout, "\n\n=========\n\n")
        return true
    },
}

var providers []*provider

func (p *provider) register() {
    if len(p.command) == 0 || haveCommand(p.command) {
        providers = append(providers, p)
    }
}

func init() {
    providers = make([]*provider, 0, 4)

    // register in order of preference
    zenityProvider.register()
    xmessageProvider.register()
    stdioProvider.register()
}

func haveCommand(cmd string) bool {
    // TODO use exec.Lookpath instead of which
    return (exec.Command("sh", "-c", "which "+cmd+" > /dev/null 2>&1").Run() == nil)
}

func platformAlert(title string, message string, args...interface{}) {
    for _, provider := range(providers) {
        if provider.alert == nil { continue }
        if provider.alert(provider, title, message, args...) { break }
    }
}
