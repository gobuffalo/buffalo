package binding

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseTimeErrorParsing(t *testing.T) {
	r := require.New(t)

	timeCustom := TimeCustomTypeDecoder{
		formats: &[]string{""},
	}

	_, err := timeCustom.parseTime([]string{"this is sparta"})
	r.Error(err)
}

func TestParseTime(t *testing.T) {
	r := require.New(t)

	timeCustom := TimeCustomTypeDecoder{
		formats: &defaultRequestBinder.timeFormats,
	}

	testCases := []struct {
		input     string
		expected  time.Time
		expectErr bool
	}{
		{
			input:     "2017-01-01",
			expected:  time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC),
			expectErr: false,
		},
		{
			input:     "2018-07-13T15:34",
			expected:  time.Date(2018, time.July, 13, 15, 34, 0, 0, time.UTC),
			expectErr: false,
		},
		{
			input:     "2018-20-10T30:15",
			expected:  time.Time{},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tt, err := timeCustom.parseTime([]string{tc.input})
		if !tc.expectErr {
			r.NoError(err)
		}

		r.Equal(tc.expected, tt)
	}
}
