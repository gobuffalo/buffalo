package errx

// go2 errors
type wrapper interface {
	Unwrap() error
}

// pkg/errors
type causer interface {
	Cause() error
}

func Unwrap(err error) error {
	switch e := err.(type) {
	case wrapper:
		return e.Unwrap()
	case causer:
		return e.Cause()
	}
	return err
}
