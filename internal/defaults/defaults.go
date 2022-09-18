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

func Int(i ...int) int {
	var r int
	for _, i := range i {
		if i != 0 {
			r = i
			break
		}
	}

	return r
}

func Int64(i ...int64) int64 {
	var r int64
	for _, i := range i {
		if i != 0 {
			r = i
			break
		}
	}

	return r
}

func Float32(i ...float32) float32 {
	var r float32
	for _, i := range i {
		if i != 0.0 {
			r = i
			break
		}
	}

	return r
}

func Float64(i ...float64) float64 {
	var r float64
	for _, i := range i {
		if i != 0.0 {
			r = i
			break
		}
	}

	return r
}
