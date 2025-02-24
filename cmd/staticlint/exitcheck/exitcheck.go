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

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		if !strings.HasSuffix(filename, "main.go") {
			continue
		}
		if pass.Pkg.Name() != "main" {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			if funcMain, ok := node.(*ast.FuncDecl); ok && funcMain.Name.Name == "main" {
				for _, body := range funcMain.Body.List {
					if expr, ok := body.(*ast.ExprStmt); ok {
						if call, ok := expr.X.(*ast.CallExpr); ok {
							if selectorExpr, ok := call.Fun.(*ast.SelectorExpr); ok {
								if selectorExpr.Sel.Name == "Exit" {
									if selectorIdent, ok := selectorExpr.X.(*ast.Ident); ok {
										if selectorIdent.Name == "os" {
											pass.Reportf(selectorExpr.Pos(), "os.Exit usage in main function")
										}
									}
								}
							}
						}
					}
				}

			}
			return true
		})
	}
	return nil, nil
}
