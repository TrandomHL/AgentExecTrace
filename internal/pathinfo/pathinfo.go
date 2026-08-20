package pathinfo

import "strings"

type Kind string

const (
	WindowsDrive Kind = "windows-drive"
	UNC          Kind = "unc"
	POSIX        Kind = "posix"
	WSLMount     Kind = "wsl-mount"
	Unknown      Kind = "unknown"
)

func Classify(path string) Kind {
	if strings.HasPrefix(path, `\\`) || strings.HasPrefix(path, "//") {
		return UNC
	}
	if len(path) >= 3 && ((path[0] >= 'a' && path[0] <= 'z') || (path[0] >= 'A' && path[0] <= 'Z')) && path[1] == ':' && (path[2] == '\\' || path[2] == '/') {
		return WindowsDrive
	}
	if strings.HasPrefix(path, "/mnt/") && len(path) >= 7 && path[6] == '/' {
		return WSLMount
	}
	if strings.HasPrefix(path, "/") {
		return POSIX
	}
	return Unknown
}
