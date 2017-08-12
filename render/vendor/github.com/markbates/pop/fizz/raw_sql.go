package fizz

func (f fizzer) RawSql() interface{} {
	return func(sql string) {
		f.add(sql, nil)
	}
}
