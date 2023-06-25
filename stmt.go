package sua

func newStmt(db *DB) *stmt {
	s := &stmt{DB: db}
	s.selectStmt = newSelectStmt(s)
	return s
}

type stmt struct {
	*DB
	table table
	*selectStmt
}

func (s *stmt) Table(name string, alias string) *stmt {
	ns := s.clone()
	ns.table = table{name: name, alias: alias}
	return ns
}

func (s *stmt) SQL() string {
	selectSql := s.selectStmt.sql()
	return selectSql
}

func (s *stmt) clone() *stmt {
	selectStmt := *s.selectStmt
	ns := &stmt{
		DB:         s.DB,
		table:      s.table,
		selectStmt: &selectStmt,
	}
	selectStmt.stmt = ns
	return ns
}
