package main

import (
	"github.com/wangxudong123/gorecover/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

type analyzerPlugin struct{}

func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		analyzer.Analyzer,
	}
}

var AnalyzerPlugin *analyzerPlugin

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
