package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"regexp"
	"strings"
	"unicode"
)

type LogCall struct {
	Call     *ast.CallExpr
	Package  string
	Function string
	Message  string
	Pos      token.Pos
	End      token.Pos
	FullPos  token.Pos
}

var Analyzer = &analysis.Analyzer{
	Name:     "loglinter",
	Doc:      "Анализирует лог-записи в кодe на соответствие правилам: только строчные буквы, английский язык, отсутствие спецсимволов и эмодзи и потенциально чувствительных данных.",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func NewAnalyzer() *analysis.Analyzer {
	config, err := FindConfig()
	if err != nil {
		return Analyzer
	}

	analyzer := *Analyzer
	analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
		return runWithConfig(pass, config)
	}
	return &analyzer
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

	start, end, ok := findStringPos(pass, call)
	if !ok {
		start, end = call.Pos(), call.End()
	}

	return &LogCall{
		Call:     call,
		Package:  pkgName,
		Function: funcName,
		Message:  message,
		Pos:      start,
		End:      end,
		FullPos:  call.Pos(),
	}
}

/// ФУНКЦИИ ПРОВЕРКИ

func checkLogMessage(pass *analysis.Pass, logCall *LogCall) {
	msg := logCall.Message

	// Правило 4: чувствительные данные
	if hasSensitiveData(msg) {
		pass.Reportf(logCall.Pos,
			"лог-сообщение не должно содержать потенциально чувствительные данные: %q",
			msg)
		return
	}

	// Правило 3: спецсимволы и эмодзи
	if !strings.Contains(msg, "%") && hasSpecialChars(msg) {
		pass.Reportf(logCall.Pos,
			"лог-сообщение не должно содержать спецсимволы или эмодзи: %q",
			msg)
		return
	}

	// Правило 2: только английские символы
	if !strings.Contains(msg, "%") && !isEnglishOnly(msg) {
		pass.Reportf(logCall.Pos,
			"лог-сообщение должно содержать только английские символы: %q",
			msg)
		return
	}

	// Правило 1: строчная буква
	if !strings.Contains(msg, "%") && len(msg) > 0 {
		firstChar := msg[0]
		if firstChar >= 'A' && firstChar <= 'Z' {
			pass.Reportf(logCall.Pos,
				"лог-сообщение должно начинаться со строчной буквы: %q",
				msg)
			return
		}
	}
}

func suggestLowercaseFix(pass *analysis.Pass, logCall *LogCall, msg string) {
	fixedMsg := strings.ToLower(string(msg[0])) + msg[1:]

	pass.Report(analysis.Diagnostic{
		Pos:     logCall.FullPos,
		End:     logCall.Call.End(),
		Message: fmt.Sprintf("лог-сообщение должно начинаться со строчной буквы: %q", msg),
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: "сделать первую букву строчной",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     logCall.Pos,
						End:     logCall.End,
						NewText: []byte(fixedMsg),
					},
				},
			},
		},
	})
}

func suggestSpecialCharsFix(pass *analysis.Pass, logCall *LogCall, msg string) {
	var fixedMsg strings.Builder
	for _, r := range msg {
		if strings.ContainsRune(".,:;_ ", r) {
			fixedMsg.WriteRune(r)
			continue
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			fixedMsg.WriteRune(r)
			continue
		}
	}

	if fixedMsg.Len() == 0 {
		fixedMsg.WriteString("message")
	}

	pass.Report(analysis.Diagnostic{
		Pos:     logCall.FullPos,
		End:     logCall.Call.End(),
		Message: fmt.Sprintf("лог-сообщение не должно содержать спецсимволы или эмодзи: %q", msg),
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: "удалить спецсимволы и эмодзи",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     logCall.Pos,
						End:     logCall.End,
						NewText: []byte(fixedMsg.String()),
					},
				},
			},
		},
	})
}

func findStringPos(pass *analysis.Pass, call *ast.CallExpr) (token.Pos, token.Pos, bool) {
	if len(call.Args) == 0 {
		return 0, 0, false
	}

	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return 0, 0, false
	}
	start := lit.Pos() + 1

	end := lit.End() - 1

	return start, end, true
}

func isEnglishOnly(s string) bool {
	for _, r := range s {
		if unicode.IsSpace(r) || unicode.IsDigit(r) {
			continue
		}
		if r > unicode.MaxASCII {
			if r >= 0x0400 && r <= 0x04FF {
				return false
			}
		}
	}
	return true
}

func hasSpecialChars(s string) bool {
	specialChars := "!?@#$%^&*()+={}[]|\\/<>`~"

	allowedChars := ".,:;_"

	for _, r := range s {
		if (r >= 0x1F300 && r <= 0x1F9FF) ||
			(r >= 0x2600 && r <= 0x26FF) ||
			(r >= 0x2700 && r <= 0x27BF) {
			return true
		}

		isAllowed := false
		for _, allowed := range allowedChars {
			if r == allowed {
				isAllowed = true
				break
			}
		}
		if isAllowed {
			continue
		}

		for _, sc := range specialChars {
			if r == sc {
				return true
			}
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
		`\bpassword\b`,
		`\bpasswd\b`,
		`\bpwd\b`,
		`\btoken\b`,
		`\bsecret\b`,
		`\bapi[_-]?key\b`,
		`\bauth\b`,
		`\bcredential\b`,
		`\bcertificate\b`,
		`\bprivate[_\s]?key\b`,
		`\bjwt\b`,
	}

	lowerMsg := strings.ToLower(s)

	for _, pattern := range sensitivePatterns {
		re := regexp.MustCompile(pattern)
		if re.MatchString(lowerMsg) {
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

func init() {
	Analyzer = NewAnalyzer()
}
