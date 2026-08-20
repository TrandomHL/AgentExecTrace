package resolve

import (
	"os"
	"path/filepath"
	"strings"
)

type Platform string

const Windows Platform = "windows"

const POSIX Platform = "posix"

const MaxCandidates = 128

type Candidate struct {
	Path       string `json:"path"`
	Exists     bool   `json:"exists"`
	Provenance string `json:"provenance,omitempty"`
}

type Result struct {
	Name       string      `json:"name"`
	Candidates []Candidate `json:"candidates"`
	Selected   string      `json:"selected,omitempty"`
	Truncated  bool        `json:"candidates_truncated"`
	Reason     string      `json:"reason"`
}

func Resolve(name string, dirs, extensions []string, platform Platform) Result {
	result := Result{Name: name}
	explicitExtension := filepath.Ext(name) != ""
	if explicitExtension || platform != Windows {
		extensions = []string{""}
	}
	for _, dir := range dirs {
		for _, extension := range extensions {
			candidate := filepath.Join(dir, name)
			if extension != "" && !strings.HasSuffix(strings.ToLower(name), strings.ToLower(extension)) {
				candidate += extension
			}
			exists, provenance := candidateInfo(candidate, platform)
			if len(result.Candidates) < MaxCandidates {
				result.Candidates = append(result.Candidates, Candidate{Path: candidate, Exists: exists, Provenance: provenance})
			} else {
				result.Truncated = true
			}
			if exists && result.Selected == "" {
				result.Selected = candidate
			}
		}
	}
	if result.Selected == "" {
		result.Reason = "no candidate exists"
	} else {
		result.Reason = "first existing candidate in PATH/PATHEXT order"
	}
	if result.Truncated {
		result.Reason += "; candidate list truncated"
	}
	return result
}

func candidateInfo(path string, platform Platform) (bool, string) {
	info, err := os.Stat(path)
	if err == nil {
		return usable(info, path, platform)
	}
	if platform != Windows {
		return false, ""
	}
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		return false, ""
	}
	for _, entry := range entries {
		if strings.EqualFold(entry.Name(), filepath.Base(path)) {
			info, err := os.Stat(filepath.Join(filepath.Dir(path), entry.Name()))
			if err != nil {
				return false, ""
			}
			return usable(info, filepath.Join(filepath.Dir(path), entry.Name()), platform)
		}
	}
	return false, ""
}

func usable(info os.FileInfo, path string, platform Platform) (bool, string) {
	if info.IsDir() || (platform != Windows && info.Mode().Perm()&0o111 == 0) {
		return false, ""
	}
	return true, provenance(path)
}

func provenance(path string) string {
	lower := strings.ToLower(path)
	base := filepath.Base(lower)
	switch {
	case base == "wsl.exe":
		return "WSL launcher"
	case base == "git-bash.exe" || strings.Contains(lower, "\\git\\bin\\bash.exe") || strings.Contains(lower, "/git/bin/bash") || strings.Contains(lower, "msys"):
		return "Git Bash / MSYS candidate"
	case strings.HasSuffix(lower, ".cmd"):
		return "cmd shim"
	case strings.HasSuffix(lower, ".bat"):
		return "bat shim"
	case strings.HasSuffix(lower, ".ps1"):
		return "PowerShell script"
	case strings.HasSuffix(lower, ".exe"):
		return "native executable"
	default:
		return "unknown"
	}
}
