package main

import (
	"github.com/jst-Frenzy/myLinter/internal/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
