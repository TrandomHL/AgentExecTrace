package redact

import (
	"strings"
	"testing"
)

func TestTextRedactsNamedAndRecognizedSecrets(t *testing.T) {
	input := "API_TOKEN=sk_test_NOT_A_REAL_SECRET\nAuthorization: Bearer abcdefghijklmnopqrstuvwxyz123456\n-----BEGIN PRIVATE KEY-----\npretend-key-material\n-----END PRIVATE KEY-----\neyJhbGciOiJub25lIn0.eyJub3QtYS1yZWFsLXNlY3JldCI6dHJ1ZX0.signature\nnormal=value"
	got := Text(input)
	for _, secret := range []string{"sk_test_NOT_A_REAL_SECRET", "abcdefghijklmnopqrstuvwxyz123456", "pretend-key-material", "eyJhbGciOiJub25lIn0"} {
		if strings.Contains(got, secret) {
			t.Fatalf("report retained secret %q: %q", secret, got)
		}
	}
	if !strings.Contains(got, "normal=value") || !strings.Contains(got, "<REDACTED") {
		t.Fatalf("unexpected redaction: %q", got)
	}
}
