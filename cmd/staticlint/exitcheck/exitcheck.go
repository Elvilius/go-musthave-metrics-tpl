package exitcheck

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "Check for os.Exit usage in main function",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		if !strings.HasSuffix(filename, "main.go") {
			continue
		}
		if pass.Pkg.Name() != "main" {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			funcMain, ok := node.(*ast.FuncDecl)
			if !ok || funcMain.Name.Name != "main" {
				return true
			}

			for _, body := range funcMain.Body.List {
				expr, ok := body.(*ast.ExprStmt)

				if !ok {
					continue
				}
				call, ok := expr.X.(*ast.CallExpr)
				if !ok {
					continue
				}

				selectorExpr, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					continue
				}

				if selectorExpr.Sel.Name != "Exit" {
					continue
				}

				if selectorIdent, ok := selectorExpr.X.(*ast.Ident); ok && selectorIdent.Name == "os" {
					pass.Reportf(selectorExpr.Pos(), "os.Exit usage in main function")

				}

			}
			return true
		})
	}
	return nil, nil
}
