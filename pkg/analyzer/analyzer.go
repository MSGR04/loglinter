package analyzer

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"strings"
	"unicode"
)

type LogCall struct {
	Call     *ast.CallExpr
	Package  string
	Function string
	Message  string
	Pos      token.Pos
}

var Analyzer = &analysis.Analyzer{
	Name:     "loglinter",
	Doc:      "Анализирует лог-записи в кодe на соответствие правилам: только строчные буквы, английский язык, отсутствие спецсимволов и эмодзи и потенциально чувствительных данных.",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		logCall := extractLogCall(pass, call)
		if logCall != nil {
			checkLogMessage(pass, logCall)
		}
	})
	return nil, nil
}

func extractLogCall(pass *analysis.Pass, call *ast.CallExpr) *LogCall {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok {
		return nil
	}

	pkgName := pkgIdent.Name
	funcName := sel.Sel.Name

	if !isKnownLogger(pkgName, funcName) {
		return nil
	}
	if len(call.Args) == 0 {
		return nil
	}

	message := extractStringArg(pass, call.Args[0])
	if message == "" {
		return nil
	}

	return &LogCall{
		Call:     call,
		Package:  pkgName,
		Function: funcName,
		Message:  message,
		Pos:      call.Pos(),
	}
}

/// ФУНКЦИИ ПРОВЕРКИ

func checkLogMessage(pass *analysis.Pass, logCall *LogCall) {
	msg := logCall.Message
	if len(msg) > 0 {
		firstchar := msg[0]
		if firstchar >= 'A' && firstchar <= 'Z' {
			pass.Reportf(logCall.Pos,
				"лог-сообщение должно начинаться со строчной буквы: %q",
				msg)
		}
	}

	if !isEnglishOnly(msg) {
		pass.Reportf(logCall.Pos,
			"лог-сообщение должно содержать только английские символы: %q",
			msg)
	}

	if hasSpecialChars(msg) {
		pass.Reportf(logCall.Pos,
			"лог-сообщение не должно содержать спецсимволы или эмодзи: %q",
			msg)
	}

	if hasSensitiveData(msg) {
		pass.Reportf(logCall.Pos,
			"лог-сообщение не должно содержать потенциально чувствительные данные: %q",
			msg)
	}
}

func isEnglishOnly(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			if r >= 0x0400 && r <= 0x04FF {
				return false
			}
		}
	}
	return true
}

func hasSpecialChars(s string) bool {
	specialChars := "!?.,:;@#$%^&*()-+={}[]|\\/<>`~"

	for _, r := range s {
		for _, sc := range specialChars {
			if r == sc {
				return true
			}
		}

		if (r >= 0x1F300 && r <= 0x1F9FF) ||
			(r >= 0x2600 && r <= 0x26FF) ||
			(r >= 0x2700 && r <= 0x27BF) ||
			(r >= 0x1FA70 && r <= 0x1FAFF) {
			return true
		}
	}

	return false
}

func isKnownLogger(pkg, fn string) bool {
	knownPackages := map[string]bool{
		"log":  true,
		"slog": true,
		"zap":  true,
	}

	logFunctions := map[string]bool{
		"Print":   true,
		"Printf":  true,
		"Println": true,
		"Fatal":   true,
		"Fatalf":  true,
		"Panic":   true,
		"Panicf":  true,
		"Debug":   true,
		"Info":    true,
		"Warn":    true,
		"Error":   true,
		"DPanic":  true,
	}

	return knownPackages[pkg] && logFunctions[fn]
}

func hasSensitiveData(s string) bool {
	sensitivePatterns := []string{
		"password",
		"passwd",
		"pwd",
		"token",
		"secret",
		"api_key",
		"apikey",
		"auth",
		"credential",
		"certificate",
		"private key",
		"jwt",
	}
	lowerMsg := strings.ToLower(s)
	for _, pattern := range sensitivePatterns {
		if strings.Contains(lowerMsg, pattern) {
			return true
		}
	}

	return false
}

func extractStringArg(pass *analysis.Pass, arg ast.Expr) string {
	switch v := arg.(type) {
	case *ast.BasicLit:
		if v.Kind == token.STRING {
			return strings.Trim(v.Value, "\"")
		}
	case *ast.BinaryExpr:
		if left, ok := v.X.(*ast.BasicLit); ok && left.Kind == token.STRING {
			return strings.Trim(left.Value, "\"")
		}
	}
	return ""
}
