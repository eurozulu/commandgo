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

