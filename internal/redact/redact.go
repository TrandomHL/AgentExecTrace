package redact

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var patterns = []*regexp.Regexp{
	regexp.MustCompile(`(?im)^(\s*[A-Za-z_][A-Za-z0-9_]*(?:TOKEN|SECRET|PASSWORD|KEY)[A-Za-z0-9_]*\s*[:=]\s*)([^\r\n]+)`),
	regexp.MustCompile(`(?i)(authorization\s*:\s*bearer\s+)([A-Za-z0-9._~-]{12,})`),
	regexp.MustCompile(`\bsk-[A-Za-z0-9_-]{12,}\b`),
	regexp.MustCompile(`(?s)-----BEGIN [A-Z ]*PRIVATE KEY-----.*?-----END [A-Z ]*PRIVATE KEY-----`),
	regexp.MustCompile(`\beyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{6,}\b`),
}

// Text returns a shareable text copy with common credential shapes replaced.
func Text(input string) string {
	output, _ := sanitizeText(input)
	return output
}

type Summary struct{ Secret, Token, HomePath int }

var windowsHome = regexp.MustCompile(`(?i)[a-z]:\\users\\[^\\/\s]+`)
var posixHome = regexp.MustCompile(`/home/[^/\s]+`)

func Report(input string) string {
	var value any
	var summary Summary
	if json.Unmarshal([]byte(input), &value) == nil {
		value = sanitizeJSON(value, &summary)
		var data bytes.Buffer
		encoder := json.NewEncoder(&data)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "  ")
		_ = encoder.Encode(value)
		return markdown("json", data.String(), summary)
	}
	output, summary := sanitizeText(input)
	return markdown("text", output, summary)
}

func sanitizeJSON(value any, summary *Summary) any {
	switch typed := value.(type) {
	case map[string]any:
		for key, item := range typed {
			if isSecretName(key) {
				typed[key] = "<REDACTED:secret>"
				summary.Secret++
				continue
			}
			typed[key] = sanitizeJSON(item, summary)
		}
	case []any:
		for index, item := range typed {
			typed[index] = sanitizeJSON(item, summary)
		}
	case string:
		output, itemSummary := sanitizeText(typed)
		*summary = Summary{Secret: summary.Secret + itemSummary.Secret, Token: summary.Token + itemSummary.Token, HomePath: summary.HomePath + itemSummary.HomePath}
		return output
	}
	return value
}

func isSecretName(key string) bool {
	upper := strings.ToUpper(key)
	return upper == "AUTHORIZATION" || strings.HasSuffix(upper, "_TOKEN") || strings.HasSuffix(upper, "_SECRET") || strings.HasSuffix(upper, "_PASSWORD") || strings.HasSuffix(upper, "_KEY")
}

func sanitizeText(input string) (string, Summary) {
	output, summary := input, Summary{}
	replace := func(pattern *regexp.Regexp, replacement string, count *int) {
		*count += len(pattern.FindAllStringIndex(output, -1))
		output = pattern.ReplaceAllString(output, replacement)
	}
	replace(patterns[0], "${1}<REDACTED:secret>", &summary.Secret)
	replace(patterns[1], "${1}<REDACTED:bearer-token>", &summary.Token)
	replace(patterns[2], "<REDACTED:api-key>", &summary.Token)
	replace(patterns[3], "<REDACTED:private-key>", &summary.Secret)
	replace(patterns[4], "<REDACTED:jwt>", &summary.Token)
	replace(windowsHome, "<HOME>", &summary.HomePath)
	replace(posixHome, "<HOME>", &summary.HomePath)
	return output, summary
}

func markdown(kind, body string, summary Summary) string {
	return fmt.Sprintf("# AgentExecTrace Shareable Report\n\n```%s\n%s\n```\n\n## Redactions:\n\n- secret: %d\n- token: %d\n- home_path: %d\n", kind, strings.TrimSuffix(body, "\n"), summary.Secret, summary.Token, summary.HomePath)
}
