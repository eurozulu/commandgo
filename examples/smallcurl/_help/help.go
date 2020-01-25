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

const helpData = `{"get":{"packagename":"SmallCurl.Get","parameters":"string, error","helptext":"Get performs a http GET on the given url, returning the BODY"}}
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
