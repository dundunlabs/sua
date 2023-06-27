package sua

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type stmtOp string

const (
	opSelect stmtOp = "SELECT"
	opInsert stmtOp = "INSERT"
	opUpdate stmtOp = "UPDATE"
	opDelete stmtOp = "DELETE"
)

func newStmt(db *DB) *stmt {
	return &stmt{DB: db}
}

type stmt struct {
	*DB
	op     stmtOp
	table  table
	cols   columns
	wc     *whereClause
	models []map[string]any
	limit  int
	offset int
}

func (s *stmt) Table(name string, alias string) *stmt {
	ns := s.clone()
	ns.table = table{name: name, alias: alias}
	return ns
}

func (s *stmt) Select(cols ...string) *stmt {
	ns := s.clone()
	ns.cols = append(ns.cols, cols...)
	return ns
}

func (s *stmt) Model(models ...map[string]any) *stmt {
	ns := s.clone()
	ns.models = append(ns.models, models...)
	return ns
}

func (s *stmt) Limit(v int) *stmt {
	ns := s.clone()
	ns.limit = v
	return ns
}

func (s *stmt) Offset(v int) *stmt {
	ns := s.clone()
	ns.offset = v
	return ns
}

func (s *stmt) Where(values ...map[string]any) *stmt {
	return s.where(values...)
}

func (s *stmt) WhereGroup(g *stmt) *stmt {
	return s.whereGroup(g)
}

func (s *stmt) WhereNot(values ...map[string]any) *stmt {
	ns := s.Where(values...)
	ns.wc.not = true
	return ns
}

func (s *stmt) WhereNotGroup(g *stmt) *stmt {
	ns := s.whereGroup(g)
	ns.wc.not = true
	return ns
}

func (s *stmt) And(values ...map[string]any) *stmt {
	return s.where(values...)
}

func (s *stmt) AndGroup(g *stmt) *stmt {
	return s.whereGroup(g)
}

func (s *stmt) AndNot(values ...map[string]any) *stmt {
	ns := s.And(values...)
	ns.wc.not = true
	return ns
}

func (s *stmt) AndNotGroup(g *stmt) *stmt {
	ns := s.whereGroup(g)
	ns.wc.not = true
	return ns
}

func (s *stmt) Or(values ...map[string]any) *stmt {
	ns := s.where(values...)
	ns.wc.or = true
	return ns
}

func (s *stmt) OrGroup(g *stmt) *stmt {
	ns := s.whereGroup(g)
	ns.wc.or = true
	return ns
}

func (s *stmt) OrNot(values ...map[string]any) *stmt {
	ns := s.Or(values...)
	ns.wc.not = true
	return ns
}

func (s *stmt) OrNotGroup(g *stmt) *stmt {
	ns := s.whereGroup(g)
	ns.wc.or = true
	ns.wc.not = true
	return ns
}

func (s *stmt) SQL() string {
	switch s.op {
	case opSelect:
		return s.selectSql()
	case opInsert:
		return s.insertSql()
	case opUpdate:
		return s.updateSql()
	case opDelete:
		return s.deleteSql()
	default:
		return ""
	}
}

func (s *stmt) selectSql() string {
	ns := s.table.aliasOrName()

	r := fmt.Sprintf(
		"%s %s FROM %q AS %q",
		s.op,
		s.cols.string(ns),
		s.table.name,
		ns,
	)

	if s.wc != nil {
		if ws := s.wc.string(ns); ws != "" {
			r += " WHERE " + ws
		}
	}

	r += s.limitOffsetString()

	return r
}

func (s *stmt) insertSql() string {
	fields := s.fields()
	cols := []string{}
	for _, f := range fields {
		cols = append(cols, fmt.Sprintf("%q", f))
	}

	values := []string{}
	for _, m := range s.models {
		value := []string{}
		for _, f := range fields {
			if v, ok := m[f]; ok {
				switch v := v.(type) {
				case string:
					value = append(value, fmt.Sprintf("'%s'", v))
				default:
					value = append(value, fmt.Sprintf("%v", v))
				}
			} else {
				value = append(value, "DEFAULT")
			}
		}
		values = append(values, fmt.Sprintf("(%s)", strings.Join(value, ", ")))
	}

	return fmt.Sprintf(
		"%s INTO %q (%s) VALUES %s",
		opInsert,
		s.table.name,
		strings.Join(cols, ", "),
		strings.Join(values, ", "),
	)
}

func (s *stmt) updateSql() string {
	updates := []string{}
	ws := ""
	if s.wc != nil {
		if w := s.wc.string(""); w != "" {
			ws = " WHERE " + w
		}
	}
	for _, m := range s.models {
		values := []string{}
		cols := maps.Keys(m)
		sort.Strings(cols)

		for _, k := range cols {
			v := m[k]
			switch v := v.(type) {
			case string:
				values = append(values, fmt.Sprintf("%q='%s'", k, v))
			default:
				values = append(values, fmt.Sprintf("%q=%v", k, v))
			}
		}
		u := fmt.Sprintf(
			"%s %q SET %s%s",
			opUpdate,
			s.table.name,
			strings.Join(values, ", "),
			ws,
		)

		updates = append(updates, u)
	}
	return strings.Join(updates, "; ")
}

func (s *stmt) deleteSql() string {
	ns := s.table.aliasOrName()

	r := fmt.Sprintf(
		"%s FROM %q AS %q",
		s.op,
		s.table.name,
		ns,
	)

	if s.wc != nil {
		if ws := s.wc.string(ns); ws != "" {
			r += " WHERE " + ws
		}
	}

	return r
}

func (s *stmt) limitOffsetString() string {
	r := ""
	if s.limit > 0 {
		r += fmt.Sprintf(" LIMIT %d", s.limit)
	}
	if s.offset > 0 {
		r += fmt.Sprintf(" OFFSET %d", s.offset)
	}
	return r
}

func (s *stmt) where(values ...map[string]any) *stmt {
	return s.whereFunc(func(w *whereClause) *whereClause {
		g := whereClause{}

		for _, value := range values {
			fields := maps.Keys(value)
			sort.Strings(fields)

			for _, k := range fields {
				v := value[k]
				c := cond{
					field: k,
					op:    "=",
					value: v,
				}

				g = g.append(func(w *whereClause) {
					w.cond = c
				})
			}
		}

		w.cond = condGroup(g)
		return w
	})
}

func (s *stmt) whereGroup(g *stmt) *stmt {
	return s.whereFunc(func(w *whereClause) *whereClause {
		w.cond = condGroup(*g.wc)
		return w
	})
}

func (s *stmt) whereFunc(fn func(w *whereClause) *whereClause) *stmt {
	ns := s.clone()
	ns.wc = fn(&whereClause{before: s.wc})
	return ns
}

func (s *stmt) fields() []string {
	fields := []string{}
	for _, m := range s.models {
		for k, _ := range m {
			if !slices.Contains(fields, k) {
				fields = append(fields, k)
			}
		}
	}
	sort.Strings(fields)
	return fields
}

func (s *stmt) clone() *stmt {
	ns := *s
	return &ns
}
