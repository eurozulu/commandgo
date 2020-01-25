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

package _help

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const helpData = `{"get":{"packagename":"http.Get","parameters":"","helptext":"func Get(url string) (resp *Response, err error)\n    Get issues a GET to the specified URL. If the response is one of the\n    following redirect codes, Get follows the redirect, up to a maximum of 10\n    redirects:\n\n    301 (Moved Permanently)\n    302 (Found)\n    303 (See Other)\n    307 (Temporary Redirect)\n    308 (Permanent Redirect)\n\n    An error is returned if there were too many redirects or if there was an\n    HTTP protocol error. A non-2xx response doesn't cause an error. Any returned\n    error will be of type *url.Error. The url.Error value's Timeout method will\n    report true if request timed out or was canceled.\n\n    When err is nil, resp always contains a non-nil resp.Body. Caller should\n    close resp.Body when done reading from it.\n\n    Get is a wrapper around DefaultClient.Get.\n\n    To make a request with custom headers, use NewRequest and DefaultClient.Do."}}
`

// CommandDetails represents a single command entry in the data
type commandDetails struct {
	FullName   string `json:"packagename"`
	Parameters string `json:"parameters"`
	HelpText   string `json:"helptext"`
}

// Help displays the available commands and a short description of what they do.
func HelpCommand(arg string) string {
	var m map[string]*commandDetails
	if err := json.Unmarshal([]byte(helpData), &m); err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)
	for k, v := range m {
		if arg == "-" || arg == k {
			if _, err := fmt.Fprintf(buf, "%s\t%s\t\t%s\n\n", k, v.Parameters, v.HelpText); err != nil {
				panic(err)
			}
		}
	}
	return buf.String()
}
