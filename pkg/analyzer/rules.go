package analyzer

import (
	"go/ast"
	"regexp"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var defaultSensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)password`),
	regexp.MustCompile(`(?i)passwd`),
	regexp.MustCompile(`(?i)pwd`),
	regexp.MustCompile(`(?i)token`),
	regexp.MustCompile(`(?i)secret`),
	regexp.MustCompile(`(?i)api[_-]?key`),
	regexp.MustCompile(`(?i)auth`),
	regexp.MustCompile(`(?i)credential`),
	regexp.MustCompile(`(?i)certificate`),
	regexp.MustCompile(`(?i)private.*key`),
	regexp.MustCompile(`(?i)jwt`),
}

type Config struct {
	EnableLowercase     bool     `json:"enable_lowercase"`
	EnableEnglishOnly   bool     `json:"enable_english_only"`
	EnableSpecialChars  bool     `json:"enable_special_chars"`
	EnableSensitiveData bool     `json:"enable_sensitive_data"`
	SensitivePatterns   []string `json:"sensitive_patterns"`
}

func NewAnalyzerWithConfig(config Config) *analysis.Analyzer {
	analyzer := *Analyzer
	analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
		return runWithConfig(pass, config)
	}
	return &analyzer
}

func runWithConfig(pass *analysis.Pass, config Config) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		logCall := extractLogCall(pass, call)

		if logCall != nil {
			checkLogMessageWithConfig(pass, logCall, config)
		}
	})

	return nil, nil
}

func checkLogMessageWithConfig(pass *analysis.Pass, logCall *LogCall, config Config) {
	msg := logCall.Message

	if config.EnableLowercase && len(msg) > 0 {
		firstchar := msg[0]
		if firstchar >= 'A' && firstchar <= 'Z' {
			pass.Reportf(logCall.Pos,
				"лог-сообщение должно начинаться со строчной буквы: %q",
				msg)
		}
	}

	if config.EnableEnglishOnly && !isEnglishOnly(msg) {
		pass.Reportf(logCall.Pos,
			"лог-сообщение должно содержать только английские символы: %q",
			msg)
	}

	if config.EnableSpecialChars && hasSpecialChars(msg) {
		pass.Reportf(logCall.Pos,
			"лог-сообщение не должно содержать спецсимволы или эмодзи: %q",
			msg)
	}

	if config.EnableSensitiveData {
		for _, pattern := range defaultSensitivePatterns {
			if pattern.MatchString(msg) {
				pass.Reportf(logCall.Pos,
					"лог-сообщение не должно содержать потенциально чувствительные данные: %q",
					msg)
				break
			}
		}

		for _, patternStr := range config.SensitivePatterns {
			pattern, err := regexp.Compile("(?i)" + patternStr)
			if err == nil && pattern.MatchString(msg) {
				pass.Reportf(logCall.Pos,
					"лог-сообщение не должно содержать потенциально чувствительные данные: %q",
					msg)
				break
			}
		}
	}
}
