//go:build plugin
// +build plugin

package main

import (
	"github.com/MSGR04/loglinter.git/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

type AnalyzerPlugin struct{}

func (AnalyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		analyzer.Analyzer,
	}
}

var Analyzer AnalyzerPlugin
