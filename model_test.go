package sua

import (
	"context"
	"testing"

	"github.com/dundunlabs/sua/test"
)

type mockUser struct {
}

var mockUserModel = NewModel(mockDB, mockUser{})

func TestModel(t *testing.T) {
	t.Run("All", func(t *testing.T) {
		mockUserModel.Where()
		got := mockUserModel.All(context.Background())
		want := []mockUser{}
		test.DeepEqual(t, got, want)
	})
}
