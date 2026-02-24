package mysqlrepo

import (
	"fmt"
	"strings"
)

func MapForCreate(pairs map[string]any) (string, string, []any) {
	var columns []string
	var vals []any
	var slots []string
	for k, v := range pairs {
		columns = append(columns, k)
		vals = append(vals, v)

		if _, ok := v.(interface{ String() string }); ok && strings.Contains(fmt.Sprintf("%T", v), "UUID") {
			slots = append(slots, "UUID_TO_BIN(?)")
			continue
		}
		slots = append(slots, "?")
	}
	return strings.Join(columns, ", "), strings.Join(slots, ", "), vals
}

func MapForUpdate(pairs map[string]any) (string, []any) {
	var sets []string
	var vals []any
	for k, v := range pairs {
		// If the value is empty, don't include it in the update statement
		if v == nil || (fmt.Sprintf("%v", v) == "") {
			continue
		}

		vals = append(vals, v)
		if _, ok := v.(interface{ String() string }); ok && strings.Contains(fmt.Sprintf("%T", v), "UUID") {
			sets = append(sets, k+" = UUID_TO_BIN(?)")
			continue
		}
		sets = append(sets, k+" = ?")
	}
	return strings.Join(sets, ", "), vals
}
