package binding

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseTimeErrorParsing(t *testing.T) {
	r := require.New(t)
	_, err := parseTime([]string{"this is sparta"})
	r.Error(err)
}

func TestParseTime(t *testing.T) {
	r := require.New(t)
	tt, err := parseTime([]string{"2017-01-01"})
	r.NoError(err)
	expected := time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)
	r.Equal(expected, tt)
}
