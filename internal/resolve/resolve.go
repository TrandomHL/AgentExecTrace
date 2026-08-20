package resolve

import (
	"os"
	"path/filepath"
	"strings"
)

type Platform string

const Windows Platform = "windows"

const MaxCandidates = 128

type Candidate struct {
	Path   string `json:"path"`
	Exists bool   `json:"exists"`
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
			exists := fileExists(candidate, platform == Windows)
			if len(result.Candidates) < MaxCandidates {
				result.Candidates = append(result.Candidates, Candidate{Path: candidate, Exists: exists})
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

func fileExists(path string, caseInsensitive bool) bool {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir()
	}
	if !caseInsensitive {
		return false
	}
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if strings.EqualFold(entry.Name(), filepath.Base(path)) {
			info, err := os.Stat(filepath.Join(filepath.Dir(path), entry.Name()))
			return err == nil && !info.IsDir()
		}
	}
	return false
}
