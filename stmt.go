package sua

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type stmtOp string

const (
	opSelect stmtOp = "SELECT"
	opInsert stmtOp = "INSERT"
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
	models []map[string]any
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

func (s *stmt) Insert(models ...map[string]any) *stmt {
	ns := s.clone()
	ns.models = append(ns.models, models...)
	return ns
}

func (s *stmt) SQL() string {
	switch s.op {
	case opSelect:
		return s.selectSql()
	case opInsert:
		return s.insertSql()
	case opDelete:
		return s.deleteSql()
	default:
		return ""
	}
}

func (s *stmt) selectSql() string {
	ns := s.table.aliasOrName()

	return fmt.Sprintf(
		"%s %s FROM %q AS %q",
		s.op,
		s.cols.string(ns),
		s.table.name,
		ns,
	)
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

func (s *stmt) deleteSql() string {
	ns := s.table.aliasOrName()

	return fmt.Sprintf(
		"%s FROM %q AS %q",
		s.op,
		s.table.name,
		ns,
	)
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
	return fields
}

func (s *stmt) clone() *stmt {
	ns := *s
	return &ns
}
