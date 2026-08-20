package diff

import (
	"encoding/json"
	"fmt"

	"github.com/TrandomHL/AgentExecTrace/internal/model"
	"github.com/TrandomHL/AgentExecTrace/internal/probe"
	"github.com/TrandomHL/AgentExecTrace/internal/resolve"
)

type Priority string

const (
	HighSignal    Priority = "high_signal"
	Informational Priority = "informational"
)

type Change struct {
	Field    string   `json:"field"`
	Finding  string   `json:"finding"`
	Priority Priority `json:"priority"`
	Left     any      `json:"left"`
	Right    any      `json:"right"`
}

func Compare(left, right model.Snapshot) []Change {
	var changes []Change
	if left.Platform.OS != right.Platform.OS || left.Platform.Arch != right.Platform.Arch || left.Platform.IsWSL != right.Platform.IsWSL {
		changes = append(changes, change("platform", "execution_namespace_changed", HighSignal, left.Platform, right.Platform))
	}
	if left.CWD != right.CWD {
		changes = append(changes, change("cwd", "cwd_changed", HighSignal, left.CWD, right.CWD))
	}
	if left.PathNamespace != right.PathNamespace {
		changes = append(changes, change("path_namespace", "path_namespace_changed", HighSignal, left.PathNamespace, right.PathNamespace))
	}
	changes = append(changes, comparePath(left.Path, right.Path)...)
	if !sameStrings(left.PathExt, right.PathExt) {
		changes = append(changes, change("path_ext", "pathext_changed", Informational, left.PathExt, right.PathExt))
	}
	return changes
}

func CompareResolve(left, right resolve.Result) []Change {
	if left.Selected == "" && right.Selected == "" {
		return nil
	}
	if left.Selected == "" || right.Selected == "" {
		return []Change{change("selected", "command_missing", HighSignal, left.Selected, right.Selected)}
	}
	if left.Selected != right.Selected {
		return []Change{change("selected", "command_target_changed", HighSignal, left.Selected, right.Selected)}
	}
	return nil
}

func CompareProbe(left, right probe.Result) []Change {
	if left.ExitCode != right.ExitCode || left.LaunchError != right.LaunchError || left.Stdout != right.Stdout || left.Stderr != right.Stderr {
		return []Change{change("probe", "probe_result_changed", HighSignal, left, right)}
	}
	return nil
}

func CompareJSON(left, right []byte) ([]Change, error) {
	var leftType, rightType map[string]json.RawMessage
	if err := json.Unmarshal(left, &leftType); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(right, &rightType); err != nil {
		return nil, err
	}
	if leftType["schema_version"] != nil && rightType["schema_version"] != nil {
		var a, b model.Snapshot
		if err := json.Unmarshal(left, &a); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(right, &b); err != nil {
			return nil, err
		}
		if a.SchemaVersion != 1 || b.SchemaVersion != 1 {
			return nil, fmt.Errorf("unsupported snapshot schema")
		}
		return Compare(a, b), nil
	}
	if leftType["name"] != nil && rightType["name"] != nil {
		var a, b resolve.Result
		if err := json.Unmarshal(left, &a); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(right, &b); err != nil {
			return nil, err
		}
		return CompareResolve(a, b), nil
	}
	if leftType["argv"] != nil && rightType["argv"] != nil {
		var a, b probe.Result
		if err := json.Unmarshal(left, &a); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(right, &b); err != nil {
			return nil, err
		}
		return CompareProbe(a, b), nil
	}
	return nil, fmt.Errorf("diff inputs must both be snapshot, resolve, or probe JSON")
}

func comparePath(left, right []string) []Change {
	leftCounts, rightCounts := counts(left), counts(right)
	var changes []Change
	for _, entry := range left {
		if leftCounts[entry] > rightCounts[entry] {
			changes = append(changes, change("path", "path_entry_removed", Informational, entry, nil))
		}
		leftCounts[entry]--
	}
	for _, entry := range right {
		if rightCounts[entry] > counts(left)[entry] {
			changes = append(changes, change("path", "path_entry_added", Informational, nil, entry))
		}
		rightCounts[entry]--
	}
	if !sameStrings(commonEntries(left, right), commonEntries(right, left)) {
		changes = append(changes, change("path", "path_order_changed", Informational, left, right))
	}
	return changes
}

func change(field, finding string, priority Priority, left, right any) Change {
	return Change{Field: field, Finding: finding, Priority: priority, Left: left, Right: right}
}
func counts(entries []string) map[string]int {
	result := make(map[string]int)
	for _, entry := range entries {
		result[entry]++
	}
	return result
}
func sameStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func commonEntries(entries, other []string) []string {
	remaining := counts(other)
	var result []string
	for _, entry := range entries {
		if remaining[entry] > 0 {
			result = append(result, entry)
			remaining[entry]--
		}
	}
	return result
}

func HasField(changes []Change, field string) bool {
	for _, change := range changes {
		if change.Field == field {
			return true
		}
	}
	return false
}

func HasFinding(changes []Change, finding string) bool {
	for _, change := range changes {
		if change.Finding == finding {
			return true
		}
	}
	return false
}
func HasPriority(changes []Change, priority Priority) bool {
	for _, change := range changes {
		if change.Priority == priority {
			return true
		}
	}
	return false
}
func decodeSnapshot(data []byte, snapshot *model.Snapshot) error {
	return json.Unmarshal(data, snapshot)
}
