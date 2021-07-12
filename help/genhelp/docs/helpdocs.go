package docs

import (
	"commandgo/help"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"strings"
)

func ReadDocs(srcPath string) ([]*help.HelpSubject, error) {
	var subjects []*help.HelpSubject
	// collect package(s) for the given directory
	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, srcPath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	for _, pkg := range pkgs {
		subs := packageHelpSubjects(pkg)
		// get the docs for the current package
		p := doc.New(pkg, srcPath, 2)

		// try to bind each mapped command to its correct docs entry using the "stub" value placed there by packageHelpSubjects
		for _, sub := range subs {
			sub.HelpItem.Comment = findDocs(sub.Name, p)
			for _, it := range sub.HelpItems {
				c := findDocs(it.Comment, p)
				if c == "" {
					c = fmt.Sprintf("--- No help available for %s", it.Comment)
				}
				it.Comment = c
			}
		}
	}
	return subjects, nil
}

func findDocs(name string, docs *doc.Package) string {
	if name == docs.Name {
		return cleanDocs(docs.Doc)
	}

	for _, fn := range docs.Funcs {
		if fn.Name != name {
			continue
		}
		return cleanDocs(fn.Doc)
	}
	for _, v := range docs.Vars {
		if main.inSlice(name, v.Names) {
			return cleanDocs(v.Doc)
		}
	}
	for _, t := range docs.Types {
		if !strings.HasPrefix(name, t.Name) {
			continue
		}
		for _, v := range t.Vars {
			if main.inSlice(name, v.Names) {
				return cleanDocs(v.Doc)
			}
		}
		for _, mt := range t.Methods {
			if strings.HasSuffix(name, mt.Name) {
				return cleanDocs(mt.Doc)
			}
		}
		// Vars is empty on the Type, so need to iterate the ast node to get field (flag) comments
		fc := findField(t.Decl, name)
		if fc != "" {
			return fc
		}
	}
	return ""
}

func findField(f ast.Node, name string) string {
	var found string
	ast.Inspect(f, func(node ast.Node) bool {
		t, ok := node.(*ast.StructType)
		if ok {
			for _, f := range t.Fields.List {
				if f.Names[0].Name == name {
					found = cleanDocs(f.Doc.Text())
					break
				}
			}
		}
		return found != ""
	})
	return found
}

func cleanDocs(d string) string {
	s := strings.ReplaceAll(d, "\"", "\\\"")
	return strings.ReplaceAll(s, "\n", "\\n")
}
