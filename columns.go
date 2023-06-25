package sua

import (
	"fmt"
	"strings"
)

type columns []string

func (cs columns) string(ns string) string {
	if len(cs) == 0 {
		if ns == "" {
			return "*"
		}
		return fmt.Sprintf("%q.*", ns)
	}

	cols := make([]string, len(cs))
	for i, c := range cs {
		if ns == "" {
			cols[i] = fmt.Sprintf("%q", c)
		} else {
			cols[i] = fmt.Sprintf("%q.%q", ns, c)
		}
	}
	return strings.Join(cols, ", ")
}
