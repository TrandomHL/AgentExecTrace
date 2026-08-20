package diff

import (
	"testing"

	"github.com/agentexectrace/agentexectrace/internal/model"
)

func TestCompareReportsSemanticPathAndNamespaceChanges(t *testing.T) {
	left := model.Snapshot{SchemaVersion: 1, CWD: `C:\repo`, PathNamespace: "windows-drive", Path: []string{`C:\Git\bin`, `C:\Windows\System32`}}
	right := model.Snapshot{SchemaVersion: 1, CWD: "/home/me/repo", PathNamespace: "posix", Path: []string{"/usr/bin", "/mnt/c/Git/bin"}}
	changes := Compare(left, right)
	if !HasField(changes, "cwd") || !HasField(changes, "path_namespace") || !HasField(changes, "path") {
		t.Fatalf("missing semantic changes: %#v", changes)
	}
}
