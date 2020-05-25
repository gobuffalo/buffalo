package decoders

// TimeDecoderFn is a custom type decoder func for Time fields
func TimeDecoderFn() func([]string) (interface{}, error) {
	return func(vals []string) (interface{}, error) {
		return parseTime(vals)
	}
}
