package diff

import (
	"reflect"

	"github.com/agentexectrace/agentexectrace/internal/model"
)

type Change struct {
	Field string `json:"field"`
	Left  any    `json:"left"`
	Right any    `json:"right"`
}

func Compare(left, right model.Snapshot) []Change {
	var changes []Change
	appendChange := func(field string, a, b any) {
		if !reflect.DeepEqual(a, b) {
			changes = append(changes, Change{Field: field, Left: a, Right: b})
		}
	}
	appendChange("cwd", left.CWD, right.CWD)
	appendChange("path_namespace", left.PathNamespace, right.PathNamespace)
	appendChange("path", left.Path, right.Path)
	appendChange("path_ext", left.PathExt, right.PathExt)
	return changes
}

func HasField(changes []Change, field string) bool {
	for _, change := range changes {
		if change.Field == field {
			return true
		}
	}
	return false
}
