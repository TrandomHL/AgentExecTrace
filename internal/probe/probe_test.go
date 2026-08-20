package probe

import (
	"context"
	"os"
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
