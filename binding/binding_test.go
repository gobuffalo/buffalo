package binding

import (
	"testing"
	"time"
)

func TestParseTimeErrorParsing(t *testing.T) {
	_, err := parseTime([]string{"this is sparta"})
	if err == nil {
		t.Fatal("expected an error, got nothing")
	}
}

func TestParseTime(t *testing.T) {
	tt, err := parseTime([]string{"2017-01-01"})
	if err != nil {
		t.Fatal(err)
	}
	expected := time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)
	if tt != expected {
		t.Fatal("expected %v, got %v", expected, tt)
	}
}