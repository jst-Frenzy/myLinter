package main

import (
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"golang.org/x/tools/go/analysis"
	"myLinter/internal/analyzer"
)

func New(settings map[string]interface{}) *linter.Config {
	if settings != nil {
		if configPath, ok := settings["config"].(string); ok && configPath != "" {
			analyzer.SetConfigPath(configPath)
		}
	}

	return &linter.Config{
		Linter: goanalysis.NewLinter(
			analyzer.Analyzer.Name,
			analyzer.Analyzer.Doc,
			[]*analysis.Analyzer{analyzer.Analyzer},
			nil,
		).WithLoadMode(goanalysis.LoadModeTypesInfo),
	}
}

func main() {

}
