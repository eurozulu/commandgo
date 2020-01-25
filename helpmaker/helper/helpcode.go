package helper

const HelpImport = `import (
	"bytes"
	"encoding/json"
	"fmt"
)
`

const HelpFunction = `

// CommandDetails represents a single command entry in the data
type commandDetails struct {
	FullName   string ` + "`json:\"packagename\"`" + `
	Parameters string ` + "`json:\"parameters\"`" + `
	HelpText   string ` + "`json:\"helptext\"`" + `
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
			fmt.Fprintf(buf, "%s\t%s\t\t%s\n\n", k, v.Parameters, v.HelpText)
		}
	}
	return buf.String()
}

`
