package main

import (
	"golang.org/x/tools/go/analysis"
	"myLinter/internal/analyzer"
)

func New(settings any) ([]*analysis.Analyzer, error) {
	if conf, ok := settings.(map[string]interface{}); ok {
		if cfgPath, ok := conf["config"].(string); ok && cfgPath != "" {
			analyzer.SetConfigPath(cfgPath)
		}
	}

	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}
