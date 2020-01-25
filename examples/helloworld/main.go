// Copyright 2020 Rob Gilham
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

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
	fmt.Printf("hello %s\n", name)

	// If flag is set, then print addition output
	if gr.Reply {
		fmt.Println("How are you doing?")
	}
}

// main sets up the single command.  To test the example:
// helloworld greet john --reply
func main() {
	mainline.MustAddCommand("greet", Greeter.SayHello)
	mainline.RunCommandLine()
}
