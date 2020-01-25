// minicurl is a silly example of how small you can make a functional tool.
// Although crude, using a few lines of code, we can get a http response header with one command.
// get <url>
//
// This example shows how base functions can be used, on packages other than your own application.
package main

import (
	"github.com/eurozulu/mainline"
	"github.com/eurozulu/mainline/examples/minicurl/_help"
	"net/http"
)

func main() {
	mainline.MustAddCommand("get", http.Get)
	mainline.MustAddCommand("help", _help.HelpCommand)
	mainline.RunCommandLine()
}
