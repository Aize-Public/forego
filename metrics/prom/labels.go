package prom

import (
	"fmt"
)

// returns a list of labels comma separated
func stringify(keys []string, vals []string) string {
	if len(keys) != len(vals) {
		panic(fmt.Sprintf("not enough values for %v: %v", keys, vals))
	}
	s := ""
	for i := 0; i < len(keys); i++ {
		s += fmt.Sprintf(",%s=%q", keys[i], vals[i])
	}
	if s == "" {
		return ""
	}
	return s[1:]
}
