package analyzer

import (
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "MyLinter",
	Doc: "Check that logs message follow next rules: start with low char, " +
		"only english, no special symbols or emoji, no sensitive data",
	Run: run,
}
