// Copyright 2020 Rob Gilham
//
// Licensed under the Apache License, Version newtype.0 (the "License");
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
package main

import (
	commandgo_7 "commandgo-7"
	"fmt"
	"log"
)

const fullVersion = `
Copyright 2020 Rob Gilham

Licensed under the Apache License, Version newtype.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
`

func main() {
	var g = &URLGet{}
	var p = &URLPost{LocalFilePermissions: 0640}

	var cmds = commandgo_7.Commands{
		// top level flags and commands, available on all commands, map to global variables and functions
		"--verbose": &Verbose,
		"-v":        &Verbose,
		"version":   showAbout,
		"":          showAbout,

		// map the get command to the URLGet instance, using default "" for Get method.
		"get": commandgo_7.Commands{
			"":          g.Get,
			"local":     g.GetLocal,
			"--headers": &g.ShowHeaders,
			"-I":        &g.ShowHeaders,
		},

		// map the post command to the URLPost instance, using default "" on Post method.
		"post": commandgo_7.Commands{
			"":             p.Post,
			"local":        p.PostLocal,
			"content-type": &p.ContentType,
			"contenttype":  &p.ContentType,
			"ct":           &p.ContentType,
			"permissions":  &p.LocalFilePermissions,
			"permis":       &p.LocalFilePermissions,
			"p":            &p.LocalFilePermissions,
		},
	}

	r, err := cmds.RunArgs()
	if err != nil {
		log.Fatalln(err)
	}

	// output any results from the call
	for _, l := range r {
		fmt.Println(l)
	}
}

// showAbout gives version and copyright information about the application
func showAbout() string {
	var fullText string
	if Verbose {
		fullText = fullVersion
	}
	return fmt.Sprintf("restline.  version 0.0\tcopyright 2021 eurozulu@github.com%s", fullText)
}
