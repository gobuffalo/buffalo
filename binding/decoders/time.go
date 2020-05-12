package decoders

func TimeDecoderFn(formats []string) func([]string) (interface{}, error) {
	return func(vals []string) (interface{}, error) {
		return parseTime(vals, formats)
	}
}
