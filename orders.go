package sua

import (
	"fmt"
	"strings"
)

type Sort string

const (
	SortASC  Sort = "ASC"
	SortDESC Sort = "DESC"
)

type order struct {
	col  string
	sort Sort
}

func (o order) string(ns string) string {
	r := fmt.Sprintf("%q %s", o.col, o.sort)
	if ns == "" {
		return r
	}
	return fmt.Sprintf("%q.%s", ns, r)
}

type orders []order

func (os orders) string(ns string) string {
	r := make([]string, len(os))
	for i, o := range os {
		r[i] = o.string(ns)
	}
	return strings.Join(r, ", ")
}
