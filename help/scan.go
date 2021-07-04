package help

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

const packagename = "\"github.com/eurozulu/commandgo"
const packagetag = "commandgo"
const typeName = "Commands"

func ScanSource(srcPath string) error {
	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, srcPath, nil, parser.ImportsOnly)
	if err != nil {
		return err
	}

	var m []*cmdmap
	for _, f := range pkgs {

		m = findMaps(f)
		if m == nil {
			continue
		}

	}
	if m == nil {
		fmt.Println("no package found using %s.%s", packagetag, typeName)
	}
	return nil
}

type cmdmap struct {
	Filename string
	Commands []map[string]ast.Expr
}

func findMaps(f ast.Node) []*cmdmap {
	var maps []*cmdmap
	var cm *cmdmap

	ast.Inspect(f, func(node ast.Node) bool {
		switch t := node.(type) {
		case *ast.File:
			if !hasPackageImport(t) {
				cm = nil
				return false
			}
			cm = &cmdmap{
				Filename: t.Name.String(),
			}
			maps = append(maps, cm)
			return true

		case *ast.CompositeLit:
			if isCommandMap(t) {
				m := buildMap(t)
				if cm == nil {
					cm = &cmdmap{}
				}
				cm.Commands = append(cm.Commands, m)
			}
			return false
		}
		return true
	})
	return maps
}


func buildMap(f *ast.CompositeLit) map[string]ast.Expr {
	m := map[string]ast.Expr{}
	ast.Inspect(f, func(node ast.Node) bool {
		t, ok := node.(*ast.KeyValueExpr)
		if !ok {
			return true
		}
		m[t.Key.(*ast.BasicLit).Value] = t.Value
		return false
	})
	return m
}


func hasPackageImport(f *ast.File) bool {
	for _, im := range f.Imports {
		if strings.HasPrefix(im.Path.Value, packagename) {
			return true
		}
	}
	return false
}

func isCommandMap(f *ast.CompositeLit) bool {
	var pkgFound bool
	var mapFound bool

	ast.Inspect(f, func(node ast.Node) bool {
		switch t := node.(type) {
		case *ast.Ident:
			if t.Name == packagetag {
				pkgFound = true
				break
			}
			if pkgFound && t.Name == typeName {
				mapFound = true
				break
			}
			pkgFound = false
		default:
			return !mapFound
		}
		return false
	})
	return mapFound
}

