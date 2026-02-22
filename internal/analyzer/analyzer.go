package analyzer

import "golang.org/x/tools/go/analysis"

var Analyzer = &analysis.Analyzer{
	Name: "MyLinter",
	Doc:  "test version", //TODO: fill this field
	Run:  run,
}
