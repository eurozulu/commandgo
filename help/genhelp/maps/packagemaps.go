package maps

import (
	"commandgo/help/genhelp/misc"
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

type PackageMaps interface {
	PackageName() string
	IsEmpty() bool
	Imports() []string
}

// packageMaps scans a set of maps files for instances of a named map type and attempts to build a representation of those maps from maps
type packageMaps struct {
	mapType string
	pkg     *ast.Package
	src     string
	maps    namedMaps
}

type namedMaps map[string]map[string]string

func (pm packageMaps) Imports() []string {
	names := misc.NewStringSet()
	for _, f := range pm.pkg.Files {
		names.Add(importNames(f.Imports)...)
	}
	return names.Get()
}

func (pm packageMaps) PackageName() string {
	return pm.pkg.Name
}

func (pm packageMaps) IsEmpty() bool {
	return len(pm.maps) == 0
}

func (pm *packageMaps) addDeclaration(decl *ast.GenDecl) bool {
	// capture declarations of map type
	// skip past non VAR decls
	if decl.Tok != token.VAR {
		return true
	}

	vs, ok := decl.Specs[0].(*ast.ValueSpec)
	if !ok {
		return true
	}

	var name string
	// is it declaraion of our map type, ensure we know all the names
	if PrintExpr(vs.Type) == pm.mapType {
		for _, n := range vs.Names {
			name = pm.addMap(PrintExpr(n))
		}
	}

	// add any composite type:
	for _, cl := range vs.Values {
		pm.addCompositeLit(name, cl)
	}
	return true
}

// addAssignment checks assignments involving the sought map are captured.
// looks for indexExpr x["a"] = 123, and also for composite literal assignments. x := map[string]int {"one": 1, "two": 2,}
// assignments are only applied to known maps, which have been previously declared
func (pm *packageMaps) addAssignment(stmt *ast.AssignStmt) bool {
	switch as := stmt.Lhs[0].(type) {
	case *ast.IndexExpr:
		// indexed assignment x["y"] = abc
		n := PrintExpr(as.X)
		if pm.isKnownMap(n) {
			n = pm.addToMap(n, as.Index, stmt.Rhs[0])
			return false
		}
		return true

	case *ast.Ident:
		if pm.isKnownMap(as.Name) {
			if n := pm.addCompositeLit(as.Name, stmt.Rhs[0]); n != "" {
				return false
			}
		}
		return true

	default:
		return true
	}
}

func (pm *packageMaps) addCompositeCall(ce *ast.CallExpr) bool {
	// call func with composite  map[string]string{"1": 1, "2", 2}.thisMethod()
	se, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return true
	}
	n := PrintExpr(ce.Fun)
	nn := pm.addCompositeLit(n, se.X)
	return nn == ""
}

func (pm *packageMaps) addCompositeLit(name string, e ast.Expr) string {
	cl, ok := e.(*ast.CompositeLit)
	if !ok {
		return ""
	}
	if PrintExpr(cl.Type) != pm.mapType {
		// these are not the drones we are looking for
		return ""
	}

	for _, e := range cl.Elts {
		kv, ok := e.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		name = pm.addToMap(name, kv.Key, kv.Value)
	}
	return name
}

func (pm *packageMaps) addToMap(name string, key ast.Expr, v ast.Expr) string {
	name = pm.addMap(name)
	m := pm.maps[name]

	k := PrintExpr(key)
	sv := PrintExpr(v)
	// Check if value is a sub map
	if IsCompositeLit(v, pm.mapType) {
		n := pm.addCompositeLit(strings.Join([]string{name, k}, "."), v)
		sv = fmt.Sprintf("{{%s}}", n)
	}
	m[k] = sv
	return name
}

func (pm *packageMaps) addMap(name string) string {
	if pm.isKnownMap(name) {
		return name
	}
	name = pm.fullName(name)
	pm.maps[name] = map[string]string{}
	return name
}

func (pm packageMaps) isKnownMap(name string) bool {
	_, ok := pm.maps[name]
	if ok {
		return true
	}
	_, ok = pm.maps[pm.fullName(name)]
	return ok
}

func (pm packageMaps) fullName(name string) string {
	if strings.Contains(name, ".") {
		return name
	}
	return strings.Join([]string{pm.PackageName(), name}, ".")
}

func (pm *packageMaps) scanPackage() {
	ast.Inspect(pm.pkg, func(node ast.Node) bool {
		switch t := node.(type) {

		case *ast.GenDecl:
			return pm.addDeclaration(t)

		case *ast.AssignStmt:
			// assignment  m["a"] = 123
			return pm.addAssignment(t)

		case *ast.CallExpr:
			return pm.addCompositeCall(t)
		default:
			return true
		}
	})
}

func importNames(ims []*ast.ImportSpec) []string {
	names := make([]string, len(ims))
	for i, f := range ims {
		names[i] = f.Path.Value
	}
	return names
}

func NewPackageMaps(mapType string, srcPath string, pkg *ast.Package) *packageMaps {
	pm := &packageMaps{
		mapType: mapType,
		pkg:     pkg,
		src:     srcPath,
		maps:    namedMaps{},
	}
	pm.scanPackage()
	return pm
}
