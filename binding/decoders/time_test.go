package decoders

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseTimeErrorParsing(t *testing.T) {
	r := require.New(t)

	_, err := parseTime([]string{"this is sparta"}, []string{})
	r.Error(err)
}

func TestParseTime(t *testing.T) {
	r := require.New(t)

	formats := []string{
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
		tt, err := parseTime([]string{tc.input}, formats)
		if !tc.expectErr {
			r.NoError(err)
		}

		r.Equal(tc.expected, tt)
	}
}

func TestParseTimeConflicting(t *testing.T) {
	// RegisterTimeFormats()

	r := require.New(t)

	formats := []string{
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
	}

	tt, err := parseTime([]string{"2017-01-10"}, formats)

	r.NoError(err)
	expected := time.Date(2017, time.October, 1, 0, 0, 0, 0, time.UTC)
	r.Equal(expected, tt)
}
