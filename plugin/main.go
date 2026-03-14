//go:build plugin
// +build plugin

package main

import (
	"github.com/yourusername/loglinter/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

// AnalyzerPlugin экспортируется для golangci-lint
type AnalyzerPlugin struct{}

// GetAnalyzers возвращает список анализаторов
func (AnalyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		analyzer.Analyzer,
	}
}

// Экспортируем переменную для golangci-lint
var Analyzer AnalyzerPlugin
