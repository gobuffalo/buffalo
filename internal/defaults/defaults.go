package defaults

func String(s ...string) string {
	var r string
	for _, s := range s {
		if s != "" {
			r = s
			break
		}
	}

	return r
}

func Int(i1, i2 int) int {
	if i1 == 0 {
		return i2
	}
	return i1
}

func Int64(i1, i2 int64) int64 {
	if i1 == 0 {
		return i2
	}
	return i1
}

func Float32(i1, i2 float32) float32 {
	if i1 == 0.0 {
		return i2
	}
	return i1
}

func Float64(i1, i2 float64) float64 {
	if i1 == 0.0 {
		return i2
	}
	return i1
}
