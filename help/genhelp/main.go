package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/eurozulu/commandgo/help"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

const (
	outputName    = "help.go"
	groupTemplate = `
// Generated code
// Do not edit directly as may be overwritten

package {{.package}}

import "github.com/eurozulu/commandgo/help"

{{range $subject := .subjects}}
var _{{.Name}} = help.HelpSubject{
	Name:     "{{.Name}},",
	Comment:  "{{ clean .Comment -}},",
	Commands: map[string]string{
		{{range $cmd, $cmt := .Commands}} "{{$cmd}}": "{{clean $cmt}}",
		{{end}}
	},
	Flags: map[string]string{
		{{range $flg, $cmt := .Flags}}
			 "{{$flg}}": "{{clean $cmt}}",
		{{end}}
	},
}
{{end}}
`
)

func main() {
	var srcPath = "./"
	if len(os.Args) > 1 {
		srcPath = os.Args[1]
	}

	var outPath = path.Join(srcPath, outputName)
	if len(os.Args) > 1 {
		outPath = os.Args[1]
	}
	err := checkOutputPath(outPath)
	if err != nil {
		log.Fatalln(err)
	}

	hgs, err := help.NewHelpSubjects(srcPath)
	if err != nil {
		log.Fatalln(err)
	}
	by := writeGroups("main", hgs)
	if err := ioutil.WriteFile(outPath, by, 0644); err != nil {
		log.Fatalln(err)
	}
}

func writeGroups(pkgName string, grps []*help.HelpSubject) []byte {
	fm := template.FuncMap{
		"clean": cleanComment,
	}
	t := template.New("help group template").Funcs(fm)
	tm, err := t.Parse(groupTemplate)
	if err != nil {
		log.Fatalln(err)
	}

	out := bytes.NewBuffer(nil)
	m := map[string]interface{}{}
	m["package"] = pkgName
	m["subjects"] = grps
	if err = tm.Execute(out, m); err != nil {
		log.Fatalln(err)
	}
	return out.Bytes()
}

func cleanComment(s string) string {
	return strings.Replace(s, "\n", "\\n", -1)
}

func checkOutputPath(p string) error {
	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	fmt.Printf("%s already exists.  Overwrite?", p)
	in := bufio.NewReader(os.Stdin)
	s, err := in.ReadString('\n')
	if err != nil {
		return err
	}
	s = s[:len(s)-1]
	if !strings.EqualFold(s, "y") && !strings.EqualFold(s, "yes") {
		return os.ErrExist
	}
	return nil
}
