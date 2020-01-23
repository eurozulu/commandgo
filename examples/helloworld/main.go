// Greeter is a simple hello world example to show the mainline library in its simpliest form.
// Using a custom struct to define a single flag and using just a single method, this shows how a method func can be used to call a command
package main

import (
	"fmt"
	"github.com/eurozulu/mainline"
)

// Greeter is a simple object which prints greeting from a given name
type Greeter struct {
	Reply bool `flag:"reply,r"`
}

// SayHello will output hello to the given name
func (gr Greeter) SayHello(name string) {
	fmt.Printf("greet %s\n", name)

	// If flag is set, then print addition output
	if gr.Reply {
		fmt.Println("How are you doing?")
	}
}

// main sets up the single command.  To test the example:
// helloworld greet john --reply
func main() {
	mainline.MustAddCommand("hello", Greeter.SayHello)
	mainline.RunCommandLine()
}
