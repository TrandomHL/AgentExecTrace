package contextinfo

import (
	"os"
	"runtime"
	"strings"

	"github.com/TrandomHL/AgentExecTrace/internal/model"
	"github.com/TrandomHL/AgentExecTrace/internal/pathinfo"
)

func Snapshot(version string) (model.Snapshot, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return model.Snapshot{}, err
	}
	path := os.Getenv("PATH")
	pathExt := os.Getenv("PATHEXT")
	return model.Snapshot{
		SchemaVersion: 1,
		ToolVersion:   version,
		Platform:      model.Platform{OS: runtime.GOOS, Arch: runtime.GOARCH, IsWSL: isWSL()},
		CWD:           cwd,
		PathNamespace: string(pathinfo.Classify(cwd)),
		Path:          split(path, os.PathListSeparator),
		PathExt:       split(pathExt, ';'),
	}, nil
}

func split(value string, separator rune) []string {
	if value == "" {
		return nil
	}
	return strings.Split(value, string(separator))
}

func isWSL() bool {
	if os.Getenv("WSL_DISTRO_NAME") != "" {
		return true
	}
	if runtime.GOOS != "linux" {
		return false
	}
	data, err := os.ReadFile("/proc/sys/kernel/osrelease")
	return err == nil && strings.Contains(strings.ToLower(string(data)), "microsoft")
}
