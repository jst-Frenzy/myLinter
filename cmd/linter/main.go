package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"
	"myLinter/internal/analyzer"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
