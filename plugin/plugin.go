package plugin

import (
	"github.com/jst-Frenzy/myLinter/internal/analyzer"
	"golang.org/x/tools/go/analysis"
)

func New(settings any) ([]*analysis.Analyzer, error) {
	if conf, ok := settings.(map[string]interface{}); ok {
		if cfgPath, ok := conf["config"].(string); ok && cfgPath != "" {
			analyzer.SetConfigPath(cfgPath)
		}
	}

	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}
