package decoders

import (
	"testing"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/stretchr/testify/require"
)

func Test_NullTimeCustomDecoder_Decode(t *testing.T) {
	r := require.New(t)

	decoderFn := NullTimeDecoderFn([]string{
		"2006-02-01",
		time.RFC3339,
		"01/02/2006",
		"2006-01-02",
		"2006-01-02T15:04",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	})

	testCases := []struct {
		input     string
		expected  time.Time
		expectErr bool
	}{
		{
			input:    "2017-01-01",
			expected: time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    "2018-07-13T15:34",
			expected: time.Date(2018, time.July, 13, 15, 34, 0, 0, time.UTC),
		},
		{
			input:    "2018-20-10T30:15",
			expected: time.Time{},
		},
		{
			input:    "",
			expected: time.Time{},
		},
	}

	for _, testCase := range testCases {

		tt, err := decoderFn([]string{testCase.input})
		r.IsType(tt, nulls.Time{})
		nt := tt.(nulls.Time)

		if testCase.expectErr {
			r.Error(err)
			r.Equal(nt.Valid, false)
			continue
		}

		r.Equal(testCase.expected, nt.Time)
	}
}
