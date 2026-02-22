package analyzer

import (
	"encoding/json"
	"errors"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"io"
	"net/http"
	"net/url"
	"strings"
	"unicode"
)

var sensitiveWords = []string{"password:", "api_key=", "token:"}

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
			message = strings.Trim(message, `"`)
			if message == "" {
				return true
			}

			if message == "" {
				return true
			}

			newLog, problems := processLogMessage(message)

			if len(problems) > 0 {
				pass.Report(analysis.Diagnostic{
					Pos:     n.Args[0].Pos(),
					End:     n.Args[0].End(),
					Message: "log has issues: " + strings.Join(problems, ", "),
					SuggestedFixes: []analysis.SuggestedFix{
						{
							TextEdits: []analysis.TextEdit{
								{
									Pos:     n.Args[0].Pos(),
									End:     n.Args[0].End(),
									NewText: []byte(newLog),
								},
							},
						},
					},
				})
			}

			return true
		})
	}

	return nil, nil
}

func processLogMessage(message string) (string, []string) {
	var problems []string
	newLog := message
	if hasSensitiveData(message) {
		newLog = blurSensitiveData(newLog)
		problems = append(problems, "sensitive data")
	}
	if hasSpecialCharsOrEmoji(message) {
		newLog = removeSpecialCharsOrEmoji(newLog)
		problems = append(problems, "special char or emoji")
	}
	if !isEnglishLetter(message) {
		tmpLog, err := translateToEnglish(newLog)
		if err == nil {
			newLog = tmpLog
		}
		problems = append(problems, "non-English")
	}
	if unicode.IsUpper([]rune(message)[0]) {
		newLog = toLowerCase(newLog)
		problems = append(problems, "starting with capital letter")
	}
	return `"` + newLog + `"`, problems
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
	lowStr := strings.ToLower(str)

	for _, word := range sensitiveWords {
		if strings.Contains(lowStr, word) {
			return true
		}
	}
	return false
}

func toLowerCase(str string) string {
	runes := []rune(str)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func removeSpecialCharsOrEmoji(str string) string {
	runes := []rune(str)
	res := strings.Builder{}
	res.Grow(len(str))

	for _, v := range runes {
		if !(unicode.IsPunct(v) || unicode.IsSymbol(v)) {
			res.WriteRune(v)
		}
	}
	return res.String()
}

func translateToEnglish(text string) (string, error) {
	resp, err := http.Get("https://translate.googleapis.com/translate_a/single?client=gtx&sl=auto&tl=en&dt=t&q=" + url.QueryEscape(text))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result []interface{}
	err = json.Unmarshal(body, &result)
	if err != nil || len(result) == 0 {
		return "", err
	}

	translations, ok := result[0].([]interface{})
	if !ok || len(translations) == 0 {
		return "", errors.New("incorrect answer format")
	}

	firstTranslation, ok := translations[0].([]interface{})
	if !ok || len(firstTranslation) == 0 {
		return "", errors.New("incorrect translation format")
	}

	answ, ok := firstTranslation[0].(string)
	if !ok {
		return "", errors.New("translation text is not string")
	}

	return answ, nil
}

func blurSensitiveData(str string) string {
	answ := str
	lowStr := strings.ToLower(str)
	for _, word := range sensitiveWords {
		if idx := strings.Index(lowStr, word); idx != -1 {
			indexValStart := idx + len(word) + 1
			indexValEnd := strings.IndexAny(answ[indexValStart:], " :,;.-=")
			if indexValEnd == -1 {
				indexValEnd = len(answ)
			} else {
				indexValEnd += indexValStart
			}
			answ = answ[:indexValStart] + "MUST BE EDITED" + answ[indexValEnd+1:]
		}
	}
	return answ
}
