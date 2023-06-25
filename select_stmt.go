package sua

import (
	"fmt"
)

const (
	selectOp = "SELECT"
)

func newSelectStmt(stmt *stmt) *selectStmt {
	return &selectStmt{stmt: stmt}
}

type selectStmt struct {
	*stmt
	columns columns
}

func (s *selectStmt) Select(columns ...string) *selectStmt {
	ns := s.clone()
	ns.selectStmt.columns = append(ns.selectStmt.columns, columns...)
	return ns.selectStmt
}

func (s *selectStmt) sql() string {
	ns := s.table.aliasOrName()
	return fmt.Sprintf(
		"%s %s FROM %q AS %q",
		selectOp,
		s.columns.string(ns),
		s.table.name,
		ns,
	)
}
