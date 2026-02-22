package analyzer

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"strings"
	"unicode"
)

var loggerConfigs = []LoggerConfig{
	{
		PkgPath:    "log/slog",
		ExtractPkg: getSlogPkgName,
	},
	{
		PkgPath:    "go.uber.org/zap",
		ExtractPkg: getZapPkgName,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			n, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			if len(n.Args) == 0 {
				return true
			}

			argType := pass.TypesInfo.TypeOf(n.Args[0])
			if argType == nil || argType.String() != "string" {
				return true
			}

			found := false

			for _, cfg := range loggerConfigs {
				if pkgName := cfg.ExtractPkg(n, pass); pkgName != nil {
					pkgPath := pkgName.Imported().Path()
					if strings.HasSuffix(pkgPath, cfg.PkgPath) {
						found = true
						break
					}
				}
			}

			if !found {
				return true
			}

			message := types.ExprString(n.Args[0])
			message = strings.Trim(message, "\"")

			if message == "" {
				return true
			}

			//проверки
			if unicode.IsUpper([]rune(message)[0]) {
				pass.Reportf(n.Args[0].Pos(), "log starting with capital letter")
			}

			if !isEnglishLetter(message) {
				pass.Reportf(n.Args[0].Pos(), "log contain non-English letter")
			}

			if hasSpecialCharsOrEmoji(message) {
				pass.Reportf(n.Args[0].Pos(), "log contain special char or emoji")
			}

			if hasSensitiveData(message) {
				pass.Reportf(n.Args[0].Pos(), "log contain sensitive data")
			}

			return true
		})
	}

	return nil, nil
}

func isEnglishLetter(str string) bool {
	for _, letter := range str {
		if unicode.IsLetter(letter) && !((letter >= 'a' && letter <= 'z') || (letter >= 'A' && letter <= 'Z')) {
			return false
		}
	}
	return true
}

func hasSpecialCharsOrEmoji(str string) bool {
	for _, letter := range str {
		if unicode.IsPunct(letter) || unicode.IsSymbol(letter) {
			return true
		}
	}
	return false
}

func hasSensitiveData(str string) bool {
	var sensitiveWords = []string{"password", "api_key", "api key", "apiKey",
		"token:", "jwt", "session", "refresh"}

	lowStr := strings.ToLower(str)

	for _, word := range sensitiveWords {
		if strings.Contains(lowStr, word) {
			return true
		}
	}
	return false
}
