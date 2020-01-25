package helper

const HelpImport = `import (
	"encoding/json"
	"github.com/eurozulu/mainline"
)
`

const HelpFunction = `

// Help displays the available commands and a short description of what they do.
func HelpCommand(arg string) string {
	by, err := Asset("_help/help.json")
	if err != nil {
		return "Help information is not available"
	}

	var m map[string]commandDetails
	if err := json.Unmarshal(by, &m); err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)
	for k, v := range m {
		if arg == "" || arg == k {
			fmt.Fprintf(buf, "%s\t%s\t\t%s\n\n", k, v.Parameters, v.HelpText)
		}
	}
	return buf.String()
}

// CommandDetails represents a single command mapping to its function
type commandDetails struct {
	FullName   string ` + "`json:\"packagename\"`" + `
	Parameters string ` + "`json:\"parameters\"`" + `
	HelpText   string ` + "`json:\"helptext\"`" + `
}
`
