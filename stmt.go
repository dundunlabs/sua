package sua

import "fmt"

type stmtOp string

const (
	opSelect stmtOp = "SELECT"
	opDelete stmtOp = "DELETE"
)

func newStmt(db *DB) *stmt {
	return &stmt{DB: db}
}

type stmt struct {
	*DB
	op    stmtOp
	table table
	cols  columns
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

func (s *stmt) SQL() string {
	ns := s.table.aliasOrName()
	switch s.op {
	case opSelect:
		return fmt.Sprintf(
			"%s %s FROM %q AS %q",
			s.op,
			s.cols.string(ns),
			s.table.name,
			ns,
		)
	case opDelete:
		return fmt.Sprintf(
			"%s FROM %q AS %q",
			s.op,
			s.table.name,
			ns,
		)
	default:
		return ""
	}
}

func (s *stmt) clone() *stmt {
	ns := *s
	return &ns
}
