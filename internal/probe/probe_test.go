package probe

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestRunCapturesSeparateStreamsExitAndInvalidUTF8(t *testing.T) {
	if os.Getenv("AGENTEXECTRACE_HELPER") == "1" {
		_, _ = os.Stdout.Write([]byte("ok\xff"))
		_, _ = os.Stderr.WriteString("problem")
		os.Exit(7)
	}
	result := Run(context.Background(), []string{os.Args[0], "-test.run=TestRunCapturesSeparateStreamsExitAndInvalidUTF8"}, 1024, map[string]string{"AGENTEXECTRACE_HELPER": "1"})
	if result.ExitCode != 7 || result.Stdout.Text != "ok" || result.Stderr.Text != "problem" || result.Stdout.UTF8 {
		t.Fatalf("unexpected probe result: %#v", result)
	}
}

func TestRunRecordsLaunchFailure(t *testing.T) {
	result := Run(context.Background(), []string{filepath.Join(t.TempDir(), "does-not-exist")}, 1024, nil)
	if result.ExitCode != -1 || result.LaunchError == "" {
		t.Fatalf("unexpected launch failure result: %#v", result)
	}
}

func TestRunMarksTruncation(t *testing.T) {
	if os.Getenv("AGENTEXECTRACE_HELPER") == "truncation" {
		_, _ = os.Stdout.WriteString(strings.Repeat("x", 4096))
		return
	}
	result := Run(context.Background(), []string{os.Args[0], "-test.run=TestRunMarksTruncation"}, 16, map[string]string{"AGENTEXECTRACE_HELPER": "truncation"})
	if result.ExitCode != 0 || !result.Stdout.Truncated || len(result.Stdout.Text) != 16 {
		t.Fatalf("unexpected truncation result: %#v", result)
	}
}

func TestRunPreservesArgv(t *testing.T) {
	if os.Getenv("AGENTEXECTRACE_HELPER") == "argv" {
		return
	}
	argv := []string{os.Args[0], "-test.run=TestRunPreservesArgv", "two words", `quote"like`, "中文"}
	result := Run(context.Background(), argv, 1024, map[string]string{"AGENTEXECTRACE_HELPER": "argv"})
	if result.ExitCode != 0 || !reflect.DeepEqual(result.Argv, argv) {
		t.Fatalf("unexpected argv result: %#v", result)
	}
}
