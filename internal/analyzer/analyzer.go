package analyzer

import (
	"flag"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "MyLinter",
	Doc: "Check that logs message follow next rules: start with low char, " +
		"only english, no special symbols or emoji, no sensitive data",
	Run:   run,
	Flags: flags(),
}

func flags() flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.StringVar(&configPath, "config", "", "path to config file")
	return *fs
}
