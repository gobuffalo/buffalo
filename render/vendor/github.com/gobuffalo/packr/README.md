# packr

Packr is a simple solution for bundling static assets inside of Go binaries. Most importantly it does it in a way that is friendly to developers while they are developing.

## Installation

```text
$ go get -u github.com/gobuffalo/packr/...
```

## Usage

### In Code

The first step in using Packr is to create a new box. A box represents a folder on disk. Once you have a box you can get `string` or `[]byte` representations of the file.

```go
// set up a new box by giving it a path to a
// folder on disk:
box := packr.NewBox("./templates")

// Get the string representation of a file:
html := box.String("index.html")
// Get the string representation of a file, or an error if it doesn't exist:
html, err := box.MustString("index.html")

// Get the []byte representation of a file:
html := box.Bytes("index.html")
// Get the []byte representation of a file, or an error if it doesn't exist:
html, err := box.MustBytes("index.html")
```

### Building a Binary (the easy way)

When it comes time to build, or install, your Go binary, simply use `packr build` or `packr install` just as you would `go build` or `go install`. All flags for the `go` tool are supported and everything works the way you expect, the only difference is your static assets are now bundled in the generated binary. If you want more control over how this happens, looking at the following section on building binaries (the hard way).

### Building a Binary (the hard way)

Before you build your Go binary, run the `packr` command first. It will look for all the boxes in your code and then generate `.go` files that pack the static files into bytes that can be bundled into the Go binary.

```
$ packr
--> packing foo/foo-packr.go
--> packing example-packr.go
```

Then run your `go build command` like normal.

#### Cleaning Up

When you're done it is recommended that you run the `packr clean` command. This will remove all of the generated files that Packr created for you.

```
$ packr clean
----> cleaning up example-packr.go
----> cleaning up foo/foo-packr.go
```

Why do you want to do this? Packr first looks to the information stored in these generated files, if the information isn't there it looks to disk. This makes it easy to work with in development.
