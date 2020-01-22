package mainline_test

import (
	"fmt"
	"github.com/eurozulu/mainline"
	"net/url"
	"os"
	"testing"
)

type MyCommands struct {
	Verbose bool `flag:"verbose, v"`
	Debug   bool `flag:"d,db"`

	Hostname *url.URL `flag:"host"`
}

func (c MyCommands) DoThis(name string) {
	fmt.Printf("Do This name: %s  isVerbose: %v, isDebug: %v\n", name, c.Verbose, c.Debug)
	if c.Hostname != nil {
		fmt.Printf("%s", c.Hostname)
	}
}
func (c MyCommands) DoThat(name string, lockernumber int) {
	fmt.Printf("Do that name: %s locaker: %d isVerbose: %v, isDebug: %v\n", name, lockernumber, c.Verbose, c.Debug)
}

func TestCommands_AddCommand(t *testing.T) {

	mainline.MustAddCommand("dothis", MyCommands.DoThis)
	mainline.MustAddCommand("dothat", MyCommands.DoThat)

	if err := mainline.RunCommandLine(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
