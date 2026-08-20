package resolve

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveWindowsUsesPathAndPATHEXTOrder(t *testing.T) {
	first := t.TempDir()
	second := t.TempDir()
	if err := os.WriteFile(filepath.Join(second, "tool.cmd"), []byte("fixture"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(first, "tool.exe"), []byte("fixture"), 0o600); err != nil {
		t.Fatal(err)
	}

	result := Resolve("tool", []string{first, second}, []string{".EXE", ".CMD"}, Windows)
	if !strings.EqualFold(result.Selected, filepath.Join(first, "tool.exe")) {
		t.Fatalf("selected %q, want first PATH candidate", result.Selected)
	}
	if len(result.Candidates) != 4 {
		t.Fatalf("candidate count = %d, want 4", len(result.Candidates))
	}
}

func TestResolveWindowsMatchesPATHEXTCandidateCaseInsensitively(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tool.cmd")
	if err := os.WriteFile(path, []byte("fixture"), 0o600); err != nil {
		t.Fatal(err)
	}

	result := Resolve("tool", []string{dir}, []string{".CMD"}, Windows)
	if result.Selected == "" || !result.Candidates[0].Exists {
		t.Fatalf("case-insensitive candidate was not found: %#v", result)
	}
}

func TestResolveKeepsExplicitExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tool.cmd")
	if err := os.WriteFile(path, []byte("fixture"), 0o600); err != nil {
		t.Fatal(err)
	}
	result := Resolve("tool.cmd", []string{dir}, []string{".EXE", ".CMD"}, Windows)
	if result.Selected != path || len(result.Candidates) != 1 {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestResolveCapsReportedCandidatesButKeepsScanning(t *testing.T) {
	result := Resolve("missing", make([]string, MaxCandidates+1), []string{".EXE"}, Windows)
	if len(result.Candidates) != MaxCandidates || !result.Truncated {
		t.Fatalf("candidate cap not reported: %#v", result)
	}
}
