package sua

import (
	"testing"

	"github.com/dundunlabs/sua/test"
)

func TestCondString(t *testing.T) {
	c := cond{field: "name", op: "=", value: "John Doe"}
	test.Equal(t, c.string(""), `("name" = 'John Doe')`)
}
