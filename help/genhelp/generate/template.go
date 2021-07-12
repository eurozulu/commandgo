package generate

import (
	"bytes"
	"commandgo/help"
	"text/template"
)

const (
	//commandPackage = "github.com/eurozulu/commandgo/help"
	commandPackage = "commandgo/help"

	groupTemplate = `// Generated code
// Do not edit directly as may be overwritten

package {{.packagename}}
import help "{{.commandhelppackage}}"
{{range $grp := .subjects}}
var {{.Name}}Help = &help.HelpSubject{
	Name:     "{{.Name}}",
	Comment:  "{{.Comment}}",
	HelpItems: []*HelpItem{
		{{range $hi := .HelpItems -}} &help.HelpItem{ Name: "{{.Name}}", Comment: "{{.Comment}}", Aliases: {{.Aliases}}} {{- end}}
	},
}
{{end}}

func init() {
	help.HelpLibrary = append(help.HelpLibrary, []*help.HelpSubject {
		{{range $grp := .subjects}}{{.Name}}Help,{{end}}
	}...)
}
`
)

func WriteTemplate(packageName string, subjects []*help.HelpSubject) ([]byte, error) {
	t := template.New("commandgo help template")
	tm, err := t.Parse(groupTemplate)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	m["packagename"] = packageName
	m["commandhelppackage"] = commandPackage
	m["subjects"] = subjects

	out := bytes.NewBuffer(nil)
	if err = tm.Execute(out, m); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
