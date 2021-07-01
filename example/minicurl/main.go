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
package main

import (
	"bytes"
	"fmt"
	"github.com/eurozulu/commandgo"
	"github.com/eurozulu/commandgo/flags"
	"github.com/eurozulu/commandgo/help"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Verbose output.  Global flag, available to all commands
var Verbose bool

// minicurl is a simple example to perform http GET and POSTs from the command line.
func main() {

	// Add global variables, available  , regardless of the command
	flags.GlobalFlags["verbose"] = &Verbose
	flags.GlobalFlags["v"] = &Verbose

	// Flags mapped to struct fields are automatically detected and parsed.
	// e.g. the PostCommand ContentType flag will become active only when the PostURL method is being called.
	// Otherwise, (such as Get being called, it will throw unknown flag error
	cmds := commandgo.Commands{
		"get":  GetURL,
		"post": PostCommand.PostURL,

		// empty key will be called if no arguments ar given
		"": help.ListCommands,
	}

	if err := cmds.Run(os.Args...); err != nil {
		fmt.Println(err)
	}
}

// GetURL requires a http url, e.g. http://localhost:8080/
// Returns the body response.
// Use -v to get headers
func GetURL(u *url.URL) error {
	r, err := http.Get(u.String())
	if err != nil {
		return err
	}
	if Verbose {
		for k, v := range r.Header {
			fmt.Printf("%s:\t%s", k, strings.Join(v, ", "))
		}
	}
	return writePresponse(r, os.Stdout)
}

// PostCommand is a struct used for the Post command.
// Using a method, rather than a global function, allows flags specific to the methods on the function.
// e.g. Post uses a 'ContentType' flag, which is only specific to the post command.
type PostCommand struct {
	// specify a content type for posting using MIME type. Defaults to text/plain
	ContentType string `flag:"contenttype,content-type,ct"`
}

// PostURL to the given url, the given data
// post requires a http url and a data string, e.g. http://localhost:8080/ 'mydata to post'
//Returns the body response. Use -v to get headers
func (pc PostCommand) PostURL(u *url.URL, data string) error {
	if pc.ContentType == "" {
		pc.ContentType = "text/plain"
	}
	r, err := http.Post(u.String(), pc.ContentType, bytes.NewBufferString(data))
	if err != nil {
		return err
	}
	return writePresponse(r, os.Stdout)
}

// writePresponse reads the Body of the given response and pushes it into the given out.
// If Versobe flag is true, also writes the response headers, prior to writing the body
func writePresponse(r *http.Response, out io.Writer) error {
	if Verbose {
		if err := r.Header.WriteSubset(out, map[string]bool{}); err != nil {
			return err
		}
	}
	by, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err := r.Body.Close(); err != nil {
		log.Println(err)
	}
	by = append(by, '\n')
	if _, err = os.Stdout.Write(by); err != nil {
		return err
	}
	return nil
}
