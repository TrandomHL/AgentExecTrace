package model

// Snapshot is the versioned, portable description of one execution context.
type Snapshot struct {
	SchemaVersion int      `json:"schema_version"`
	ToolVersion   string   `json:"tool_version,omitempty"`
	Platform      Platform `json:"platform"`
	CWD           string   `json:"cwd"`
	PathNamespace string   `json:"path_namespace"`
	Path          []string `json:"path,omitempty"`
	PathExt       []string `json:"path_ext,omitempty"`
}

type Platform struct {
	OS    string `json:"os"`
	Arch  string `json:"arch"`
	IsWSL bool   `json:"is_wsl"`
}
