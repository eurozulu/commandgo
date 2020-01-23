package main

import (
	"github.com/eurozulu/mainline"
	"net/http"
)

func main() {
	mainline.MustAddCommand("get", http.Get)
	mainline.RunCommandLine()
}
