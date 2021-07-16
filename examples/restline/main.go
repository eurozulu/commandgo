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
	"fmt"
	"github.com/eurozulu/commandgo"
	"github.com/eurozulu/commandgo/examples/restline/restutils"
	"log"
)

// Sample data for additional info using ShowAbout. (To demo the Verbose flag usage)
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
	// These are our application model objects we will be mapping into
	var g = &restutils.URLGet{LocalFileRoot: "./content"}
	var p = &restutils.URLPost{LocalFileRoot: "./content", LocalFilePermissions: 0640}

	// top level flags and commands, available on all commands, usually map to global variables and functions
	var cmds = commandgo.Commands{
		"--verbose": &restutils.Verbose,
		"-v":        &restutils.Verbose,
		"version":   ShowAbout,

		// Default mapping to show about.  Invoked when no arguments are given
		"": ShowAbout,

		// map the 'get' command in subcommand so it has its own flags, seperate from the global flags.
		// maps to the URLGet instance (g), using default "" mapping to the Get method.
		// Additional command "get local ..." maps to a second method on the same instance.
		// Has a single assignment flag to show response headers, mapped twice for a short and long name.
		"get": commandgo.Commands{
			"":          g.Get,
			"local":     g.GetLocal,
			"--headers": &g.ShowHeaders,
			"-I":        &g.ShowHeaders,
		},

		// map the post command to the URLPost instance (p), using default "" on Post method.
		// Has two assignment flags mapped to multiple names, ContentType and LocalFilePermissions
		"post": commandgo.Commands{
			"":               p.Post,
			"local":          p.PostLocal,
			"--content-type": &p.ContentType,
			"--contenttype":  &p.ContentType,
			"-ct":            &p.ContentType,
			"--permissions":  &p.LocalFilePermissions,
			"--perm":         &p.LocalFilePermissions,
			"-p":             &p.LocalFilePermissions,
		},
	}

	// Call using the os.CommandLine argument
	r, err := cmds.RunArgs()
	if err != nil {
		log.Fatalln(err)
	}

	// output any results from the call
	for _, l := range r {
		fmt.Println(l)
	}
}

// ShowAbout gives version and copyright information about the application
// A simple local function invoked on the root command map. (Also default, no arguments mapping)
// Uses the Verbose flag to show full copyright data when true.
func ShowAbout() string {
	var fullText string
	if restutils.Verbose {
		fullText = fullVersion
	}
	return fmt.Sprintf("restline.  version 0.0\tcopyright 2021 eurozulu@github.com%s", fullText)
}
