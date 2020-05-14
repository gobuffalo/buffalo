package decoders

func TimeDecoderFn() func([]string) (interface{}, error) {
	return func(vals []string) (interface{}, error) {
		return parseTime(vals)
	}
}
