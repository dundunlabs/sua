package sua

import (
	"testing"

	"github.com/dundunlabs/sua/test"
)

func TestCond(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		c := cond{field: "name", op: "=", value: "John Doe"}
		test.Equal(t, c.string(""), `("name" = 'John Doe')`)
	})
}
