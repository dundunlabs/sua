package sua

import (
	"testing"

	"github.com/dundunlabs/sua/test"
)

var mockSelectStmt = newSelectStmt(mockStmt)

func TestTable(t *testing.T) {
	s := mockSelectStmt.Table("users", "u")
	test.Equal(t, s.SQL(), `SELECT "u".* FROM "users" AS "u"`)
}

func TestSelect(t *testing.T) {
	s := mockSelectStmt.Table("users", "u")
	s1 := s.Select("id", "name")
	s2 := s1.Select("gender")

	test.Equal(t, s1.SQL(), `SELECT "u"."id", "u"."name" FROM "users" AS "u"`)
	test.Equal(t, s2.SQL(), `SELECT "u"."id", "u"."name", "u"."gender" FROM "users" AS "u"`)
}
