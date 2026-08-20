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

	"github.com/agentexectrace/agentexectrace/internal/contextinfo"
	"github.com/agentexectrace/agentexectrace/internal/diff"
	"github.com/agentexectrace/agentexectrace/internal/model"
	"github.com/agentexectrace/agentexectrace/internal/probe"
	"github.com/agentexectrace/agentexectrace/internal/redact"
	"github.com/agentexectrace/agentexectrace/internal/resolve"
)

var version = "0.1.0-dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
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
	if err := flags.Parse(args); err != nil || flags.NArg() == 0 || *limit < 0 {
		return errors.New("usage: probe [--max-bytes n] [--output file] -- <command> [args...]")
	}
	return writeJSON(*output, probe.Run(context.Background(), flags.Args(), *limit, nil), stdout)
}

func diffCommand(args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("diff", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	output := flags.String("output", "", "write JSON to file")
	if err := flags.Parse(args); err != nil || flags.NArg() != 2 {
		return errors.New("usage: diff [--output file] <left.json> <right.json>")
	}
	left, err := readSnapshot(flags.Arg(0))
	if err != nil {
		return err
	}
	right, err := readSnapshot(flags.Arg(1))
	if err != nil {
		return err
	}
	return writeJSON(*output, diff.Compare(left, right), stdout)
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
	sanitized := redact.Text(string(data))
	if *output == "" {
		_, err = io.WriteString(stdout, sanitized)
		return err
	}
	return os.WriteFile(*output, []byte(sanitized), 0o600)
}

func readSnapshot(path string) (model.Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return model.Snapshot{}, err
	}
	var snapshot model.Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return model.Snapshot{}, fmt.Errorf("decode %s: %w", path, err)
	}
	if snapshot.SchemaVersion != 1 {
		return model.Snapshot{}, fmt.Errorf("unsupported snapshot schema %d", snapshot.SchemaVersion)
	}
	return snapshot, nil
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
	fmt.Fprintln(writer, "commands: snapshot, resolve, probe, diff, report --redact")
}
