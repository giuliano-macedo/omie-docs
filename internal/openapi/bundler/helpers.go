package openapi_bundler

import (
	"regexp"
	"strings"

	"golang.org/x/text/unicode/norm"
)

func formatFeatureName(featureName string) string {
	normalized := norm.NFD.Bytes([]byte(featureName))

	var stringBuilder strings.Builder
	for _, c := range normalized {
		isAlphaNum := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
		switch {
		case isAlphaNum:
			stringBuilder.WriteByte(c)
		case c == byte(' '):
			stringBuilder.WriteRune('_')
		}
	}
	return strings.ToLower(stringBuilder.String())
}

func removeJsComments(sourceCode string) string {
	singleLineRegex := regexp.MustCompile(`//.*`)
	sourceCode = singleLineRegex.ReplaceAllString(sourceCode, "")

	multiLineRegex := regexp.MustCompile(`\/\*(.|\n)*?\*\/`)
	sourceCode = multiLineRegex.ReplaceAllString(sourceCode, "")

	return sourceCode
}
