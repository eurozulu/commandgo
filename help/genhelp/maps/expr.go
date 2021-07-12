package maps

import (
	"fmt"
	"go/ast"
	"strings"
)

func IsCompositeLit(e ast.Expr, typeName string) bool {
	cl, ok := e.(*ast.CompositeLit)
	return ok && typeName == PrintExpr(cl.Type)
}

func PrintExpr(e ast.Expr) string {
	switch t := e.(type) {
	case nil:
		return ""

	case *ast.Ident:
		return t.Name
	case *ast.BasicLit:
		return t.Value
	case *ast.KeyValueExpr:
		return fmt.Sprintf("%s: %s", PrintExpr(t.Key), PrintExpr(t.Value))
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", PrintExpr(t.X), t.Sel.Name)
	case *ast.IndexExpr:
		return PrintExpr(t.Index)
	case *ast.BinaryExpr:
		return fmt.Sprintf("%s: %s", PrintExpr(t.X), PrintExpr(t.Y))
	case *ast.CallExpr:
		args := make([]string, len(t.Args))
		for i, arg := range t.Args {
			args[i] = PrintExpr(arg)
		}
		return fmt.Sprintf("%s(%s)", PrintExpr(t.Fun), strings.Join(args, ", "))
	case *ast.UnaryExpr:
		return fmt.Sprintf("%s: %s", PrintExpr(t.X), t.Op.String())
	case *ast.TypeAssertExpr:
		var tn string
		if t.Type != nil {
			tn = PrintExpr(t.Type)
		}
		return fmt.Sprintf("%s %s", PrintExpr(t.X), tn)
	case *ast.StarExpr:
		return PrintExpr(t.X)

	default:
		return "????"
	}
}
