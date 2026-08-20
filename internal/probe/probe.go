package probe

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"
)

type Capture struct {
	Text      string `json:"text"`
	UTF8      bool   `json:"utf8"`
	Truncated bool   `json:"truncated"`
}

type Result struct {
	Argv        []string `json:"argv"`
	Stdout      Capture  `json:"stdout"`
	Stderr      Capture  `json:"stderr"`
	ExitCode    int      `json:"exit_code"`
	LaunchError string   `json:"launch_error,omitempty"`
}

func Run(ctx context.Context, argv []string, limit int, overrides map[string]string) Result {
	result := Result{Argv: append([]string(nil), argv...), ExitCode: -1}
	if len(argv) == 0 {
		result.LaunchError = "no command supplied"
		return result
	}
	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Env = environment(overrides)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		result.LaunchError = err.Error()
		return result
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		result.LaunchError = err.Error()
		return result
	}
	if err := cmd.Start(); err != nil {
		result.LaunchError = err.Error()
		return result
	}
	stdoutResult := make(chan Capture, 1)
	stderrResult := make(chan Capture, 1)
	go func() { stdoutResult <- capture(stdout, limit) }()
	go func() { stderrResult <- capture(stderr, limit) }()
	result.Stdout = <-stdoutResult
	result.Stderr = <-stderrResult
	err = cmd.Wait()
	if err == nil {
		result.ExitCode = 0
		return result
	}
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		result.ExitCode = exitError.ExitCode()
		return result
	}
	result.LaunchError = err.Error()
	return result
}

func capture(reader io.Reader, limit int) Capture {
	if limit < 0 {
		limit = 0
	}
	var saved []byte
	buffer := make([]byte, 32*1024)
	truncated := false
	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			remaining := limit - len(saved)
			if remaining > 0 {
				keep := n
				if keep > remaining {
					keep = remaining
				}
				saved = append(saved, buffer[:keep]...)
			}
			if n > remaining {
				truncated = true
			}
		}
		if err != nil {
			break
		}
	}
	valid := utf8.Valid(saved)
	return Capture{Text: strings.ToValidUTF8(string(saved), ""), UTF8: valid, Truncated: truncated}
}

func environment(overrides map[string]string) []string {
	env := os.Environ()
	for key, value := range overrides {
		prefix := key + "="
		for index := range env {
			if strings.EqualFold(env[index][:strings.IndexByte(env[index], '=')+1], prefix) {
				env[index] = prefix + value
				prefix = ""
				break
			}
		}
		if prefix != "" {
			env = append(env, prefix+value)
		}
	}
	return env
}
