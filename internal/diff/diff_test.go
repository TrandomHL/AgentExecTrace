package diff

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/TrandomHL/AgentExecTrace/internal/model"
	"github.com/TrandomHL/AgentExecTrace/internal/probe"
	"github.com/TrandomHL/AgentExecTrace/internal/resolve"
)

func TestCompareReportsSemanticSnapshotChanges(t *testing.T) {
	left := model.Snapshot{SchemaVersion: 1, Platform: model.Platform{OS: "windows"}, CWD: `C:\repo`, PathNamespace: "windows-drive", Path: []string{"a", "b", "c"}, PathExt: []string{".EXE"}}
	right := model.Snapshot{SchemaVersion: 1, Platform: model.Platform{OS: "linux", IsWSL: true}, CWD: "/mnt/c/repo", PathNamespace: "wsl-mount", Path: []string{"b", "a", "d"}, PathExt: []string{".EXE", ".CMD"}}
	changes := Compare(left, right)
	for _, finding := range []string{"execution_namespace_changed", "cwd_changed", "path_namespace_changed", "path_entry_added", "path_entry_removed", "path_order_changed", "pathext_changed"} {
		if !HasFinding(changes, finding) {
			t.Fatalf("missing %s: %#v", finding, changes)
		}
	}
	if !HasPriority(changes, HighSignal) || !HasPriority(changes, Informational) {
		t.Fatalf("missing semantic changes: %#v", changes)
	}
}

func TestCompareResolveAndProbeEvidence(t *testing.T) {
	leftResolve := resolve.Result{Name: "bash", Selected: `C:\Program Files\Git\bin\bash.exe`}
	rightResolve := resolve.Result{Name: "bash", Selected: `C:\Windows\System32\bash.exe`}
	if !HasFinding(CompareResolve(leftResolve, rightResolve), "command_target_changed") {
		t.Fatal("missing command target finding")
	}
	if !HasFinding(CompareResolve(leftResolve, resolve.Result{Name: "bash"}), "command_missing") {
		t.Fatal("missing command missing finding")
	}
	leftProbe := probe.Result{ExitCode: 0, Stdout: probe.Capture{Text: "marker"}}
	rightProbe := probe.Result{ExitCode: 7, Stdout: probe.Capture{Text: "marker"}}
	if !HasFinding(CompareProbe(leftProbe, rightProbe), "probe_result_changed") {
		t.Fatal("missing probe finding")
	}
}

func TestGoldenSnapshotPairsDecodeAndCompare(t *testing.T) {
	var snapshots []model.Snapshot
	for _, name := range []string{"semantic-left.json", "semantic-right.json"} {
		data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "diff", name))
		if err != nil {
			t.Fatal(err)
		}
		var snapshot model.Snapshot
		if err := decodeSnapshot(data, &snapshot); err != nil {
			t.Fatal(err)
		}
		snapshots = append(snapshots, snapshot)
	}
	if !HasFinding(Compare(snapshots[0], snapshots[1]), "path_order_changed") {
		t.Fatal("golden snapshots did not report PATH reordering")
	}
}
