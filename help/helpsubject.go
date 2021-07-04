package help

import (
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"strings"
)

// HelpSubject is a logical grouping of commands and flags.
// the 'global' subject is called 'main', containing the available global commands and flags.
// Each struct mapped to a command is grouped into its own subject, with all its flags (fields) inthat group.
type HelpSubject struct {
	Name     string
	Comment  string
	Commands map[string]string
	Flags    map[string]string
}

func NewHelpSubjects(pkgPath string) ([]*HelpSubject, error) {
	var groups []*HelpSubject

	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, pkgPath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for pn, f := range pkgs {
		// get the docs for the current package
		p := doc.New(f, pkgPath, 2)

		gh := NewHelpSubject(pn, p.Doc)
		groups = append(groups, gh)

		for _, fn := range p.Funcs {
			gh.Commands[strings.ToLower(fn.Name)] = fn.Doc
		}
		for _, v := range p.Vars {
			gh.Flags[strings.ToLower(v.Names[0])] = v.Doc
		}
		for _, t := range p.Types {
			// Each struct Type gets its own help group
			tgh := NewHelpSubject(t.Name, t.Doc)
			groups = append(groups, tgh)
			for _, mt := range t.Methods {
				tgh.Commands[strings.ToLower(mt.Name)] = mt.Doc
			}
			// Vars is empty on the Type, so need to iterate the ast node to get field (flag) comments
			gatherFields(t.Decl, tgh.Flags)
		}
	}
	return groups, nil
}

func NewHelpSubject(name string, comment string) *HelpSubject {
	return &HelpSubject{
		Name:     name,
		Comment:  comment,
		Commands: map[string]string{},
		Flags:    map[string]string{},
	}
}

func gatherFields(f ast.Node, m map[string]string) {
	ast.Inspect(f, func(node ast.Node) bool {
		t, ok := node.(*ast.StructType)
		if ok {
			for _, f := range t.Fields.List {
				m[f.Names[0].Name] = f.Doc.Text()
			}
		}
		return true
	})
}
