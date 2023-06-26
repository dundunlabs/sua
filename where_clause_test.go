package sua

import (
	"testing"

	"github.com/dundunlabs/sua/test"
)

func TestWhere(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		w := whereClause{
			cond: cond{field: "name", op: "=", value: "John Doe"},
		}
		test.Equal(t, w.string("u"), `("u"."name" = 'John Doe')`)
	})

	t.Run("AndOr", func(t *testing.T) {
		w := whereClause{
			or:   true,
			not:  true,
			cond: cond{field: "name", op: "=", value: "John Doe"},
			before: &whereClause{
				cond: cond{field: "gender", op: "!=", value: "male"},
				before: &whereClause{
					cond: cond{field: "age", op: "<", value: 18},
				},
			},
		}
		test.Equal(t, w.string(""), `("age" < 18) AND ("gender" != 'male') OR NOT ("name" = 'John Doe')`)
	})

	t.Run("Nested", func(t *testing.T) {
		w := whereClause{
			or:  true,
			not: true,
			cond: condGroup{
				cond: cond{field: "name", op: "=", value: "John Doe"},
				before: &whereClause{
					cond: cond{field: "age", op: "<", value: 18},
				},
			},
			before: &whereClause{
				cond: cond{field: "gender", op: "!=", value: "male"},
			},
		}
		test.Equal(t, w.string(""), `("gender" != 'male') OR NOT (("age" < 18) AND ("name" = 'John Doe'))`)
	})
}
