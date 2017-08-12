# github.com/markbates/pop/nulls

This package should be used in place of the built-in null types in the `sql` package.

The real benefit of this packages comes in its implementation of `MarshalJSON` and `UnmarshalJSON` to properly encode/decode `null` values.

## Installation

``` bash
$ go get github.com/markbates/pop/nulls
```

## Supported Datatypes

* `string` (`nulls.String`) - Replaces `sql.NullString`
* `int64` (`nulls.Int64`) - Replaces `sql.NullInt64`
* `float64` (`nulls.Float64`) - Replaces `sql.NullFloat64`
* `bool` (`nulls.Bool`) - Replaces `sql.NullBool`
* `[]byte` (`nulls.ByteSlice`)
* `float32` (`nulls.Float32`)
* `int` (`nulls.Int`)
* `int32` (`nulls.Int32`)
* `uint32` (`nulls.UInt32`)
* `time.Time` (`nulls.Time`)
