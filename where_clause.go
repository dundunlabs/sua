package sua

type whereClause struct {
	not    bool
	or     bool
	cond   icond
	before *whereClause
}

func (w whereClause) append(fn func(w *whereClause)) whereClause {
	nw := w
	if w.cond != nil {
		nw = whereClause{before: &w}
	}
	fn(&nw)
	return nw
}

func (w whereClause) string(ns string) string {
	cond := w.cond.string(ns)
	r := cond
	if w.not {
		r = "NOT " + r
	}
	if w.before != nil {
		op := " AND "
		if w.or {
			op = " OR "
		}
		r = w.before.string(ns) + op + r
	}
	return r
}
