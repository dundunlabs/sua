package sua

import (
	"testing"

	"github.com/dundunlabs/sua/test"
)

var (
	mockSelectStmt = &stmt{op: opSelect}
	mockInsertStmt = &stmt{op: opInsert}
	mockDeleteStmt = &stmt{op: opDelete}
)

func TestSelect(t *testing.T) {
	s := mockSelectStmt.Table("users", "u")
	s1 := s.Select("id", "name")
	s2 := s1.Select("gender")

	test.Equal(t, s.SQL(), `SELECT "u".* FROM "users" AS "u"`)
	test.Equal(t, s1.SQL(), `SELECT "u"."id", "u"."name" FROM "users" AS "u"`)
	test.Equal(t, s2.SQL(), `SELECT "u"."id", "u"."name", "u"."gender" FROM "users" AS "u"`)
}

func TestInsert(t *testing.T) {
	mockInsertUser := mockInsertStmt.Table("users", "u")

	t.Run("Single", func(t *testing.T) {
		s := mockInsertUser.Insert(map[string]any{"name": "foo", "gender": "male"})
		test.Equal(t, s.SQL(), `INSERT INTO "users" ("name", "gender") VALUES ('foo', 'male')`)
	})

	t.Run("Multiple", func(t *testing.T) {
		s := mockInsertUser.Insert(
			map[string]any{"name": "foo"},
			map[string]any{"gender": "male"},
		)
		test.Equal(t, s.SQL(), `INSERT INTO "users" ("name", "gender") VALUES ('foo', DEFAULT), (DEFAULT, 'male')`)
	})
}

func TestDelete(t *testing.T) {
	s := mockDeleteStmt.Table("users", "u")
	test.Equal(t, s.SQL(), `DELETE FROM "users" AS "u"`)
}
