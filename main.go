package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/TrandomHL/AgentExecTrace/internal/contextinfo"
	"github.com/TrandomHL/AgentExecTrace/internal/diff"
	"github.com/TrandomHL/AgentExecTrace/internal/probe"
	"github.com/TrandomHL/AgentExecTrace/internal/redact"
	"github.com/TrandomHL/AgentExecTrace/internal/resolve"
)

var version = "0.1.0-dev"

const selfProbeFlag = "--agentexectrace-self-probe"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 && args[0] == selfProbeFlag {
		fmt.Fprintln(stdout, "agentexectrace-self-probe-stdout spaces \"quotes\" 中文")
		fmt.Fprintln(stderr, "agentexectrace-self-probe-stderr")
		return 7
	}
	if len(args) == 0 {
		usage(stderr)
		return 2
	}
	var err error
	switch args[0] {
	case "snapshot":
		err = snapshotCommand(args[1:], stdout)
	case "resolve":
		err = resolveCommand(args[1:], stdout)
	case "probe":
		err = probeCommand(args[1:], stdout)
	case "diff":
		err = diffCommand(args[1:], stdout)
	case "report":
		err = reportCommand(args[1:], stdout)
	case "help", "-h", "--help":
		usage(stdout)
		return 0
	default:
		err = fmt.Errorf("unknown command %q", args[0])
	}
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	return 0
}

func snapshotCommand(args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("snapshot", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	output := flags.String("output", "", "write JSON to file")
	if err := flags.Parse(args); err != nil || flags.NArg() != 0 {
		return errors.New("usage: snapshot [--output file]")
	}
	snapshot, err := contextinfo.Snapshot(version)
	if err != nil {
		return err
	}
	return writeJSON(*output, snapshot, stdout)
}

func resolveCommand(args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("resolve", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	output := flags.String("output", "", "write JSON to file")
	if err := flags.Parse(args); err != nil || flags.NArg() != 1 {
		return errors.New("usage: resolve [--output file] <name>")
	}
	platform := resolve.Platform("")
	if os.PathSeparator == '\\' {
		platform = resolve.Windows
	}
	path := split(os.Getenv("PATH"), os.PathListSeparator)
	pathExt := split(os.Getenv("PATHEXT"), ';')
	return writeJSON(*output, resolve.Resolve(flags.Arg(0), path, pathExt, platform), stdout)
}

func probeCommand(args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("probe", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	limit := flags.Int("max-bytes", 64*1024, "maximum stored bytes per stream")
	output := flags.String("output", "", "write JSON to file")
	if err := flags.Parse(args); err != nil || *limit < 0 {
		return errors.New("usage: probe [--max-bytes n] [--output file] [-- <command> [args...]]")
	}
	argv := flags.Args()
	if len(argv) == 0 {
		executable, err := os.Executable()
		if err != nil {
			return err
		}
		argv = []string{executable, selfProbeFlag, "spaces value", `"quoted"`, "中文"}
	}
	return writeJSON(*output, probe.Run(context.Background(), argv, *limit, nil), stdout)
}

func diffCommand(args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("diff", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	output := flags.String("output", "", "write JSON to file")
	if err := flags.Parse(args); err != nil || flags.NArg() != 2 {
		return errors.New("usage: diff [--output file] <left.json> <right.json>")
	}
	left, err := os.ReadFile(flags.Arg(0))
	if err != nil {
		return err
	}
	right, err := os.ReadFile(flags.Arg(1))
	if err != nil {
		return err
	}
	changes, err := diff.CompareJSON(left, right)
	if err != nil {
		return err
	}
	return writeJSON(*output, changes, stdout)
}

func reportCommand(args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("report", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	redactOutput := flags.Bool("redact", false, "redact common secrets")
	output := flags.String("output", "", "write sanitized report to file")
	if err := flags.Parse(args); err != nil || !*redactOutput || flags.NArg() != 1 {
		return errors.New("usage: report --redact [--output file] <input>")
	}
	data, err := os.ReadFile(flags.Arg(0))
	if err != nil {
		return err
	}
	sanitized := redact.Report(string(data))
	if *output == "" {
		_, err = io.WriteString(stdout, sanitized)
		return err
	}
	return os.WriteFile(*output, []byte(sanitized), 0o600)
}

func writeJSON(path string, value any, stdout io.Writer) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if path == "" {
		_, err = stdout.Write(data)
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func split(value string, separator rune) []string {
	if value == "" {
		return nil
	}
	return strings.Split(value, string(separator))
}

func usage(writer io.Writer) {
	fmt.Fprintln(writer, "AgentExecTrace: local execution-context evidence")
	fmt.Fprintln(writer, "commands: snapshot, resolve, probe [-- <command> [args...]], diff, report --redact")
}
