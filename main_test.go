package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/TrandomHL/AgentExecTrace/internal/model"
)

func init() {
	if len(os.Args) > 1 && os.Args[1] == selfProbeFlag {
		os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
	}
}

func TestSnapshotDoesNotDumpEnvironmentValues(t *testing.T) {
	t.Setenv("API_TOKEN", "sk_test_NEVER_IN_SNAPSHOT")
	var out, errOut bytes.Buffer
	if code := run([]string{"snapshot"}, &out, &errOut); code != 0 {
		t.Fatalf("snapshot failed: %s", errOut.String())
	}
	if strings.Contains(out.String(), "NEVER_IN_SNAPSHOT") {
		t.Fatalf("snapshot disclosed environment value: %s", out.String())
	}
	var snapshot model.Snapshot
	if err := json.Unmarshal(out.Bytes(), &snapshot); err != nil || snapshot.SchemaVersion != 1 {
		t.Fatalf("invalid snapshot: %v %#v", err, snapshot)
	}
}

func TestReportRedactWritesNewSafeCopy(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.txt")
	output := filepath.Join(dir, "output.txt")
	if err := os.WriteFile(input, []byte("DB_PASSWORD=not-a-real-password\nplain=text\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	if code := run([]string{"report", "--redact", "--output", output, input}, &out, &errOut); code != 0 {
		t.Fatalf("report failed: %s", errOut.String())
	}
	sanitized, err := os.ReadFile(output)
	if err != nil || strings.Contains(string(sanitized), "not-a-real-password") || !strings.Contains(string(sanitized), "plain=text") {
		t.Fatalf("unsafe report output: %v %q", err, sanitized)
	}
}

func TestResolveAndDiffCommandsProduceJSON(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"resolve", "agentexectrace-command-that-does-not-exist"}, &out, &errOut); code != 0 {
		t.Fatalf("resolve failed: %s", errOut.String())
	}
	if !strings.Contains(out.String(), `"reason"`) {
		t.Fatalf("resolve did not produce JSON: %s", out.String())
	}
	out.Reset()
	if code := run([]string{"diff", "testdata/snapshots/windows-agent.json", "testdata/snapshots/wsl-terminal.json"}, &out, &errOut); code != 0 {
		t.Fatalf("diff failed: %s", errOut.String())
	}
	if !strings.Contains(out.String(), `"cwd"`) || !strings.Contains(out.String(), `"path_namespace"`) {
		t.Fatalf("diff did not report semantic changes: %s", out.String())
	}
}

func TestProbeCommandCapturesResult(t *testing.T) {
	if os.Getenv("AGENTEXECTRACE_CLI_HELPER") == "1" {
		_, _ = os.Stdout.WriteString("out")
		_, _ = os.Stderr.WriteString("err")
		os.Exit(3)
	}
	t.Setenv("AGENTEXECTRACE_CLI_HELPER", "1")
	var out, errOut bytes.Buffer
	args := []string{"probe", "--max-bytes", "64", "--", os.Args[0], "-test.run=TestProbeCommandCapturesResult"}
	if code := run(args, &out, &errOut); code != 0 {
		t.Fatalf("probe failed: %s", errOut.String())
	}
	if !strings.Contains(out.String(), `"exit_code": 3`) || !strings.Contains(out.String(), `"text": "out"`) || !strings.Contains(out.String(), `"text": "err"`) {
		t.Fatalf("probe did not preserve result: %s", out.String())
	}
}

func TestProbeWithoutArgumentsRunsSelfProbe(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"probe"}, &out, &errOut); code != 0 {
		t.Fatalf("self-probe failed: %s", errOut.String())
	}
	for _, want := range []string{`"exit_code": 7`, "agentexectrace-self-probe-stdout", "agentexectrace-self-probe-stderr", "中文"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("self-probe missing %q: %s", want, out.String())
		}
	}
}
