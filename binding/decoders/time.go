package decoders

// TimeDecoderFn is a custom type decoder func for Time fields
func TimeDecoderFn() func([]string) (any, error) {
	return func(vals []string) (any, error) {
		return parseTime(vals)
	}
}
