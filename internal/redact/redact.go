package redact

import "regexp"

var patterns = []*regexp.Regexp{
	regexp.MustCompile(`(?im)^(\s*[A-Za-z_][A-Za-z0-9_]*(?:TOKEN|SECRET|PASSWORD|API_KEY|PRIVATE_KEY|ACCESS_KEY)[A-Za-z0-9_]*\s*[:=]\s*)([^\r\n]+)`),
	regexp.MustCompile(`(?i)(authorization\s*:\s*bearer\s+)([A-Za-z0-9._~-]{12,})`),
	regexp.MustCompile(`\bsk-[A-Za-z0-9_-]{12,}\b`),
	regexp.MustCompile(`(?s)-----BEGIN [A-Z ]*PRIVATE KEY-----.*?-----END [A-Z ]*PRIVATE KEY-----`),
	regexp.MustCompile(`\beyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{6,}\b`),
}

// Text returns a shareable text copy with common credential shapes replaced.
func Text(input string) string {
	output := patterns[0].ReplaceAllString(input, "${1}<REDACTED:secret>")
	output = patterns[1].ReplaceAllString(output, "${1}<REDACTED:bearer-token>")
	output = patterns[2].ReplaceAllString(output, "<REDACTED:api-key>")
	output = patterns[3].ReplaceAllString(output, "<REDACTED:private-key>")
	return patterns[4].ReplaceAllString(output, "<REDACTED:jwt>")
}
