package main

import (
	"errors"
	"go/ast"

	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

const (
	mainFuncName      = "main"
	osExitMethod      = "Exit"
	osIdentName       = "os"
	directOsExitError = "direct os.Exit call in main function is not allowed"
)

var errNoAnalysisNeeded = errors.New("no analysis needed")

var NoOsExitAnalyzer = &analysis.Analyzer{
	Name: "noosexit",
	Doc:  "checks for direct os.Exit calls in main function of main package",
	Run:  runNoOsExitAnalyzer,
}

func runNoOsExitAnalyzer(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, errNoAnalysisNeeded
	}
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.Name == mainFuncName {
				checkForOsExitCall(pass, fn.Body)
			}
		}
	}
	return nil, errNoAnalysisNeeded
}

func checkForOsExitCall(pass *analysis.Pass, body *ast.BlockStmt) {
	ast.Inspect(body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			checkOsExitCall(pass, call)
		}
		return true
	})
}

func checkOsExitCall(pass *analysis.Pass, call *ast.CallExpr) {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == osIdentName && sel.Sel.Name == osExitMethod {
			pass.Reportf(call.Pos(), directOsExitError)
		}
	}
}

func main() {
	multichecker.Main(getAnalyzers()...)
}

func getAnalyzers() []*analysis.Analyzer {
	myChecks := []*analysis.Analyzer{
		NoOsExitAnalyzer,
		bodyclose.Analyzer,
		shadow.Analyzer,
		fieldalignment.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		structtag.Analyzer,
		shift.Analyzer,
	}
	return appendStaticcheckAnalyzers(myChecks)
}

func appendStaticcheckAnalyzers(checks []*analysis.Analyzer) []*analysis.Analyzer {
	for _, v := range staticcheck.Analyzers {
		if v.Analyzer.Name[:2] == "SA" || v.Analyzer.Name == "ST" {
			checks = append(checks, v.Analyzer)
		}
	}
	return checks
}
