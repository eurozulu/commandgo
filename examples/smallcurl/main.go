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

// smallcurl is another triviol example which builds on the mini curl.
// Still keeping it simple, it adds a struct to handle the response.
// With the struct it can use flags, so adds a flag to get the headers as well as the body
package main

import (
	"bytes"
	"fmt"
	"github.com/eurozulu/mainline"
	"github.com/eurozulu/mainline/examples/smallcurl/_help"
	"io"
	"net/http"
	"net/url"
)

type SmallCurl struct {
	// Header, when true, displays the response headers at the start of the body
	Header bool `flag:"header,i"`
	// Nobody, when true, will NOT return the body stream.  Used with Header to get just headers, or on its own gets response code.
	Nobody bool
}

// Get performs a http GET on the given url, returning the BODY
func (sc SmallCurl) Get(u *url.URL) (string, error) {
	r, err := http.Get(u.String())
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(nil)
	if sc.Header {
		for k, v := range r.Header {
			buf.WriteString(fmt.Sprintf("%s = %v\n", k, v))
		}
		fmt.Println()
	}

	if !sc.Nobody {
		io.Copy(buf, r.Body)
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
	}()

	if buf.Len() == 0 {
		buf.WriteString(r.Status)
	}
	return buf.String(), nil
}

func main() {
	mainline.AddCommand("get", SmallCurl.Get)
	mainline.MustAddCommand("help", _help.HelpCommand)
	mainline.RunCommandLine()
}
