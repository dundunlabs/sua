package sua

import "fmt"

type icond interface {
	string(ns string) string
}

type cond struct {
	field string
	op    string
	value any
}

func (c cond) string(ns string) string {
	f := fmt.Sprintf("%q", c.field)
	if ns != "" {
		f = fmt.Sprintf("%q.%s", ns, f)
	}
	value := fmt.Sprintf("%v", c.value)
	switch v := c.value.(type) {
	case string:
		value = fmt.Sprintf("'%s'", v)
	}
	return fmt.Sprintf("(%s %s %v)", f, c.op, value)
}

type condGroup whereClause

func (g condGroup) string(ns string) string {
	return fmt.Sprintf("(%s)", whereClause(g).string(ns))
}
