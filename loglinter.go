package loglinter

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "loglinter",
	Doc:      "Анализирует лог-записи в кодe на соответствие правилам: только строчные буквы, английский язык, отсутствие спецсимволов и эмодзи и потенциально чувствительных данных.",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

var sensitiveWords = []string{
	"password", "passwd", "pwd",
	"token", "api_key", "apikey", "secret",
	"key", "auth", "credential",
}

func run(pass analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)


	})
}

func isLogFunction(pass *analysis.Pass, call *ast.CallExpr) bool{
	if 
}
