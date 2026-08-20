package model

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSnapshotFixturesDecode(t *testing.T) {
	for _, name := range []string{"windows-agent.json", "wsl-terminal.json"} {
		data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "snapshots", name))
		if err != nil {
			t.Fatal(err)
		}
		var snapshot Snapshot
		if err := json.Unmarshal(data, &snapshot); err != nil {
			t.Fatalf("decode %s: %v", name, err)
		}
		if snapshot.SchemaVersion != 1 || snapshot.CWD == "" {
			t.Fatalf("unexpected fixture %s: %#v", name, snapshot)
		}
	}
}
