package sua

type table struct {
	name  string
	alias string
}

func (t table) aliasOrName() string {
	if t.alias != "" {
		return t.alias
	}
	return t.name
}
