package sua

import "context"

func NewModel[M any](db *DB, model M) *Model[M] {
	return &Model[M]{db, model}
}

type Model[M any] struct {
	*DB
	model M
}

func (m *Model[M]) All(ctx context.Context) []M {
	ms := []M{}
	return ms
}
