package mainline_test

import (
	"fmt"
	"github.com/eurozulu/mainline"
	"net/url"
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
func (c MyCommands) DoTheNumbers(i int, f float32, b bool) {
	fmt.Printf("the numbers are: int: %v float: %v bool: %v\n", i, f, b)
}

func TestCommands_AddCommand(t *testing.T) {

	mainline.MustAddCommand("dothis", MyCommands.DoThis)
	mainline.MustAddCommand("dothat", MyCommands.DoThat)

	if err := mainline.RunCommand("dothis", "world"); err != nil {
		t.Fatal(err)
	}

	if err := mainline.RunCommand("dothat", "world", "22"); err != nil {
		t.Fatal(err)
	}
	if err := mainline.RunCommand("dothis"); err == nil {
		t.Fatal(fmt.Errorf("expected error of not enough params"))
	}
	if err := mainline.RunCommand("dothis", "hello", "world"); err == nil {
		t.Fatal(fmt.Errorf("expected error of too many params"))
	}

	if err := mainline.RunCommand("dothat", "hello", "world"); err == nil {
		t.Fatal(fmt.Errorf("expected unparsable int error"))
	}
}

func TestCommands_NumberParams(t *testing.T) {
	mainline.MustAddCommand("numbers", MyCommands.DoTheNumbers)

	if err := mainline.RunCommand("numbers", "1", "1.2", "true"); err != nil {
		t.Fatal(err)
	}
}
