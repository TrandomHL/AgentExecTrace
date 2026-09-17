package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/TrandomHL/AgentExecTrace/internal/diff"
	"github.com/TrandomHL/AgentExecTrace/internal/probe"
	"github.com/TrandomHL/AgentExecTrace/internal/resolve"
)

// This fixture is entered only by explicitly launched test subprocesses.
func TestControlledFixture(t *testing.T) {
	if os.Getenv("AET_CONTROLLED_FIXTURE") != "1" {
		return
	}
	if os.Args[len(os.Args)-1] == "read-relative" {
		data, err := os.ReadFile("fixture.txt")
		if err != nil {
			fmt.Fprint(os.Stderr, "fixture-missing")
			os.Exit(3)
		}
		fmt.Fprint(os.Stdout, string(data))
	} else {
		executable, err := os.Executable()
		if err != nil {
			os.Exit(4)
		}
		fmt.Fprint(os.Stdout, filepath.Base(filepath.Dir(executable)))
	}
	os.Exit(0)
}

// Controlled examples use the built production CLI and real child processes.
// They never change the parent environment or publish raw diagnostic output.
func TestControlledExamples(t *testing.T) {
	root := t.TempDir()
	suffix := ""
	if runtime.GOOS == "windows" {
		suffix = ".exe"
	}
	binary := filepath.Join(root, "agentexectrace"+suffix)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	build := exec.CommandContext(ctx, "go", "build", "-o", binary, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	fixtureBytes, err := os.ReadFile(self)
	if err != nil {
		t.Fatal(err)
	}
	fixture := func(t *testing.T, name string) string {
		t.Helper()
		dir := filepath.Join(root, name)
		if err := os.Mkdir(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(dir, "aetfixture"+suffix)
		if err := os.WriteFile(path, fixtureBytes, 0o700); err != nil {
			t.Fatal(err)
		}
		return path
	}
	first := fixture(t, "tool-A")
	second := fixture(t, "tool-B")
	firstDir, secondDir := filepath.Dir(first), filepath.Dir(second)
	sameFile := func(a, b string) bool {
		left, err := os.Stat(a)
		if err != nil {
			return false
		}
		right, err := os.Stat(b)
		return err == nil && os.SameFile(left, right)
	}
	environment := func(path, extensions string) []string {
		var env []string
		for _, entry := range os.Environ() {
			key, _, _ := strings.Cut(entry, "=")
			if !strings.EqualFold(key, "PATH") && !strings.EqualFold(key, "PATHEXT") && key != "AET_CONTROLLED_FIXTURE" {
				env = append(env, entry)
			}
		}
		return append(env, "PATH="+path, "PATHEXT="+extensions, "AET_CONTROLLED_FIXTURE=1")
	}
	normal := environment(firstDir, ".EXE;.CMD;.BAT")
	call := func(t *testing.T, cwd string, env []string, args ...string) []byte {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, binary, args...)
		cmd.Dir, cmd.Env = cwd, env
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("CLI %q: %v\n%s", args, err, out)
		}
		return out
	}
	decode := func(t *testing.T, data []byte, value any) {
		t.Helper()
		if err := json.Unmarshal(data, value); err != nil {
			t.Fatalf("decode: %v", err)
		}
	}
	compare := func(t *testing.T, left, right []byte, want string) {
		t.Helper()
		dir := t.TempDir()
		a, b := filepath.Join(dir, "left.json"), filepath.Join(dir, "right.json")
		for path, data := range map[string][]byte{a: left, b: right} {
			if err := os.WriteFile(path, data, 0o600); err != nil {
				t.Fatal(err)
			}
		}
		var changes []diff.Change
		decode(t, call(t, root, normal, "diff", a, b), &changes)
		if want == "" {
			if len(changes) != 0 {
				t.Fatalf("restored control has %d changes", len(changes))
			}
		} else if !diff.HasFinding(changes, want) {
			t.Fatalf("missing finding %s", want)
		}
	}
	probeTool := func(t *testing.T, cwd string, env []string, tool string, extra ...string) probe.Result {
		t.Helper()
		args := append([]string{"probe", "--", tool, "-test.run=^TestControlledFixture$", "--"}, extra...)
		var result probe.Result
		decode(t, call(t, cwd, env, args...), &result)
		return result
	}

	t.Run("cwd", func(t *testing.T) {
		project, wrong := filepath.Join(root, "project"), filepath.Join(root, "other")
		for _, dir := range []string{project, wrong} {
			if err := os.Mkdir(dir, 0o700); err != nil {
				t.Fatal(err)
			}
		}
		if err := os.WriteFile(filepath.Join(project, "fixture.txt"), []byte("fixture-found"), 0o600); err != nil {
			t.Fatal(err)
		}
		before := call(t, project, normal, "snapshot")
		after := call(t, wrong, normal, "snapshot")
		compare(t, before, after, "cwd_changed")
		good := probeTool(t, project, normal, first, "read-relative")
		bad := probeTool(t, wrong, normal, first, "read-relative")
		restored := probeTool(t, project, normal, first, "read-relative")
		if good.ExitCode != 0 || good.Stdout.Text != "fixture-found" || bad.ExitCode != 3 || bad.Stderr.Text != "fixture-missing" || restored.ExitCode != 0 || restored.Stdout != good.Stdout {
			t.Fatal("relative-file symptom/control failed")
		}
		compare(t, before, call(t, project, normal, "snapshot"), "")
		t.Log("CONTROLLED: relative file succeeds -> exit 3 in other CWD; finding=cwd_changed; restored CWD succeeds, no snapshot differences")
	})

	t.Run("pathext", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			t.Skip("Windows executable lookup experiment; not evidence for POSIX/WSL PATHEXT behavior")
		}
		broken := environment(firstDir, ".CPL")
		compare(t, call(t, root, normal, "snapshot"), call(t, root, broken, "snapshot"), "pathext_changed")
		goodJSON := call(t, root, normal, "resolve", "aetfixture")
		badJSON := call(t, root, broken, "resolve", "aetfixture")
		var good, bad resolve.Result
		decode(t, goodJSON, &good)
		decode(t, badJSON, &bad)
		if !sameFile(good.Selected, first) || bad.Selected != "" {
			t.Fatal("PATHEXT candidate selection mismatch")
		}
		compare(t, goodJSON, badJSON, "command_missing")
		launched := probeTool(t, root, normal, "aetfixture")
		missing := probeTool(t, root, broken, "aetfixture")
		explicit := probeTool(t, root, broken, first)
		if launched.ExitCode != 0 || launched.Stdout.Text != "tool-A" || missing.LaunchError == "" || explicit.ExitCode != 0 || explicit.Stdout.Text != "tool-A" {
			t.Fatal("PATHEXT process symptom/control failed")
		}
		compare(t, goodJSON, call(t, root, normal, "resolve", "aetfixture"), "")
		t.Log("CONTROLLED: .EXE present -> launch succeeds; .CPL only -> launch fails; findings=pathext_changed,command_missing; absolute executable still works; restored lookup matches")
	})

	t.Run("path_order", func(t *testing.T) {
		left := environment(firstDir+string(os.PathListSeparator)+secondDir, ".EXE;.CMD;.BAT")
		right := environment(secondDir+string(os.PathListSeparator)+firstDir, ".EXE;.CMD;.BAT")
		compare(t, call(t, root, left, "snapshot"), call(t, root, right, "snapshot"), "path_order_changed")
		a := call(t, root, left, "resolve", "aetfixture")
		b := call(t, root, right, "resolve", "aetfixture")
		var ra, rb resolve.Result
		decode(t, a, &ra)
		decode(t, b, &rb)
		if !sameFile(ra.Selected, first) || !sameFile(rb.Selected, second) {
			t.Fatal("PATH selection mismatch")
		}
		compare(t, a, b, "command_target_changed")
		pa, pb := probeTool(t, root, left, "aetfixture"), probeTool(t, root, right, "aetfixture")
		if pa.ExitCode != 0 || pb.ExitCode != 0 || pa.Stdout.Text != "tool-A" || pb.Stdout.Text != "tool-B" {
			t.Fatal("selected process output mismatch")
		}
		restored := probeTool(t, root, left, "aetfixture")
		if restored.ExitCode != 0 || restored.Stdout != pa.Stdout {
			t.Fatal("restored PATH did not restore output")
		}
		compare(t, a, call(t, root, left, "resolve", "aetfixture"), "")
		t.Log("CONTROLLED: identical command name prints tool-A -> tool-B after PATH reorder; findings=path_order_changed,command_target_changed; restored PATH prints tool-A")
	})

	t.Run("documented_resolve_commands", func(t *testing.T) {
		for _, file := range []string{"README.md", "AI_AGENT_SETUP.md"} {
			text, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			count := 0
			for _, line := range strings.Split(string(text), "\n") {
				if !strings.HasPrefix(line, `.\agentexectrace.exe resolve `) && !strings.HasPrefix(line, `& $tool resolve `) {
					continue
				}
				count++
				dir := t.TempDir()
				args := strings.Fields(line)[1:]
				if args[0] == "$tool" {
					args = args[1:]
				}
				data := call(t, dir, normal, args...)
				for i, arg := range args {
					if arg == "--output" {
						data, err = os.ReadFile(filepath.Join(dir, args[i+1]))
						if err != nil {
							t.Fatal(err)
						}
					}
				}
				var result resolve.Result
				decode(t, data, &result)
				if result.Name != "git" {
					t.Fatalf("unexpected documented command in %s", file)
				}
			}
			if count == 0 {
				t.Fatalf("no resolve examples tested in %s", file)
			}
		}
	})

	t.Run("first_use_files", func(t *testing.T) {
		dir := t.TempDir()
		a, b := filepath.Join(dir, "terminal.json"), filepath.Join(dir, "agent.json")
		changes, report := filepath.Join(dir, "diff.json"), filepath.Join(dir, "report.md")
		call(t, firstDir, normal, "snapshot", "--output", a)
		call(t, secondDir, normal, "snapshot", "--output", b)
		call(t, root, normal, "diff", "--output", changes, a, b)
		original, err := os.ReadFile(changes)
		if err != nil {
			t.Fatal(err)
		}
		call(t, root, normal, "report", "--redact", "--output", report, changes)
		text, err := os.ReadFile(report)
		if err != nil || !bytes.Contains(text, []byte("cwd_changed")) {
			t.Fatal("first-use report missing finding")
		}
		unchanged, err := os.ReadFile(changes)
		if err != nil || !bytes.Equal(original, unchanged) {
			t.Fatal("report changed source evidence")
		}
	})
}
