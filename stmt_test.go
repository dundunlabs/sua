package sua

import (
	"testing"

	"github.com/dundunlabs/sua/test"
)

var (
	mockSelectStmt = &stmt{op: opSelect}
	mockInsertStmt = &stmt{op: opInsert}
	mockUpdateStmt = &stmt{op: opUpdate}
	mockDeleteStmt = &stmt{op: opDelete}
)

func TestStmt(t *testing.T) {
	t.Run("Select", func(t *testing.T) {
		s := mockSelectStmt.Table("users", "u")
		s1 := s.Select("id", "name")
		s2 := s1.Select("gender")

		test.Equal(t, s.SQL(), `SELECT "u".* FROM "users" AS "u"`)
		test.Equal(t, s1.SQL(), `SELECT "u"."id", "u"."name" FROM "users" AS "u"`)
		test.Equal(t, s2.SQL(), `SELECT "u"."id", "u"."name", "u"."gender" FROM "users" AS "u"`)
	})

	t.Run("Insert", func(t *testing.T) {
		mockInsertUser := mockInsertStmt.Table("users", "u")

		t.Run("Single", func(t *testing.T) {
			s := mockInsertUser.Model(map[string]any{"name": "foo", "gender": "male"})
			test.Equal(t, s.SQL(), `INSERT INTO "users" ("gender", "name") VALUES ('male', 'foo')`)
		})

		t.Run("Multiple", func(t *testing.T) {
			s := mockInsertUser.Model(
				map[string]any{"name": "foo"},
				map[string]any{"gender": "male"},
			)
			test.Equal(t, s.SQL(), `INSERT INTO "users" ("gender", "name") VALUES (DEFAULT, 'foo'), ('male', DEFAULT)`)
		})
	})

	t.Run("Update", func(t *testing.T) {
		mockUpdateUser := mockUpdateStmt.Table("users", "u")

		t.Run("Single", func(t *testing.T) {
			s := mockUpdateUser.Model(map[string]any{"name": "foo", "gender": "male"})
			test.Equal(t, s.SQL(), `UPDATE "users" SET "gender"='male', "name"='foo'`)
		})

		t.Run("Multiple", func(t *testing.T) {
			s := mockUpdateUser.Model(
				map[string]any{"name": "foo"},
				map[string]any{"gender": "male"},
			)
			test.Equal(t, s.SQL(), `UPDATE "users" SET "name"='foo'; UPDATE "users" SET "gender"='male'`)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		s := mockDeleteStmt.Table("users", "u")
		test.Equal(t, s.SQL(), `DELETE FROM "users" AS "u"`)
	})

	t.Run("Where", func(t *testing.T) {
		s1 := mockSelectStmt.Table("users", "u").
			Where(map[string]any{"name": "John Doe", "gender": "male"})
		s2 := s1.And(map[string]any{"age": 18}).
			Or(map[string]any{"age": 20})
		s3 := mockDeleteStmt.Table("users", "u").
			WhereNot(map[string]any{"age": 18}).
			OrNot(map[string]any{"age": 20}).
			AndNot(map[string]any{"gender": "male"})
		s4 := mockUpdateStmt.Table("users", "u").
			Model(map[string]any{"name": "Foo Bar"}).
			WhereGroup(new(stmt).
				Where(map[string]any{"id": 1}).
				Or(map[string]any{"id": 2})).
			AndNotGroup(new(stmt).
				Where(map[string]any{"id": 3}).
				Or(map[string]any{"id": 4}))

		test.Equal(t, s1.SQL(), `SELECT "u".* FROM "users" AS "u" WHERE (("u"."gender" = 'male') AND ("u"."name" = 'John Doe'))`)
		test.Equal(t, s2.SQL(), `SELECT "u".* FROM "users" AS "u" WHERE (("u"."gender" = 'male') AND ("u"."name" = 'John Doe')) AND (("u"."age" = 18)) OR (("u"."age" = 20))`)
		test.Equal(t, s3.SQL(), `DELETE FROM "users" AS "u" WHERE NOT (("u"."age" = 18)) OR NOT (("u"."age" = 20)) AND NOT (("u"."gender" = 'male'))`)
		test.Equal(t, s4.SQL(), `UPDATE "users" SET "name"='Foo Bar' WHERE ((("id" = 1)) OR (("id" = 2))) AND NOT ((("id" = 3)) OR (("id" = 4)))`)
	})

	t.Run("LimitOffset", func(t *testing.T) {
		s := mockSelectStmt.Table("users", "u").Limit(25).Offset(100)
		test.Equal(t, s.SQL(), `SELECT "u".* FROM "users" AS "u" LIMIT 25 OFFSET 100`)
	})
}
