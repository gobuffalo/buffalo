go
==

[![Build Status](https://travis-ci.org/shurcooL/go.svg?branch=master)](https://travis-ci.org/shurcooL/go) [![GoDoc](https://godoc.org/github.com/shurcooL/go?status.svg)](https://godoc.org/github.com/shurcooL/go)

Common Go code.

Installation
------------

```bash
go get -u github.com/shurcooL/go/...
```

Directories
-----------

| Path                                                                                                  | Synopsis                                                                                                                                                          |
|-------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [analysis](https://godoc.org/github.com/shurcooL/go/analysis)                                         | Package analysis provides a routine that determines if a file is generated or handcrafted.                                                                        |
| [browser](https://godoc.org/github.com/shurcooL/go/browser)                                           | Package browser provides utilities for interacting with users' browsers.                                                                                          |
| [ctxhttp](https://godoc.org/github.com/shurcooL/go/ctxhttp)                                           | Package ctxhttp provides helper functions for performing context-aware HTTP requests.                                                                             |
| [gddo](https://godoc.org/github.com/shurcooL/go/gddo)                                                 | Package gddo is a simple client library for accessing the godoc.org API.                                                                                          |
| [generated](https://godoc.org/github.com/shurcooL/go/generated)                                       | Package generated provides a function that parses a Go file and reports whether it contains a "// Code generated ...                                              |
| [gfmutil](https://godoc.org/github.com/shurcooL/go/gfmutil)                                           | Package gfmutil offers functionality to render GitHub Flavored Markdown to io.Writer.                                                                             |
| [gopathutil](https://godoc.org/github.com/shurcooL/go/gopathutil)                                     | Package gopathutil provides tools to operate on GOPATH workspace.                                                                                                 |
| [gopherjs_http](https://godoc.org/github.com/shurcooL/go/gopherjs_http)                               | Package gopherjs_http provides helpers for compiling Go using GopherJS and serving it over HTTP.                                                                  |
| [gopherjs_http/jsutil](https://godoc.org/github.com/shurcooL/go/gopherjs_http/jsutil)                 | Package jsutil provides utility functions for interacting with native JavaScript APIs.                                                                            |
| [importgraphutil](https://godoc.org/github.com/shurcooL/go/importgraphutil)                           | Package importgraphutil augments "golang.org/x/tools/refactor/importgraph" with a way to build graphs ignoring tests.                                             |
| [indentwriter](https://godoc.org/github.com/shurcooL/go/indentwriter)                                 | Package indentwriter implements an io.Writer wrapper that indents every non-empty line with specified number of tabs.                                             |
| [ioutil](https://godoc.org/github.com/shurcooL/go/ioutil)                                             | Package ioutil provides a WriteFile func with an io.Reader as input.                                                                                              |
| [open](https://godoc.org/github.com/shurcooL/go/open)                                                 | Package open offers ability to open files or URLs as if user double-clicked it in their OS.                                                                       |
| [openutil](https://godoc.org/github.com/shurcooL/go/openutil)                                         | Package openutil displays Markdown or HTML in a new browser tab.                                                                                                  |
| [ospath](https://godoc.org/github.com/shurcooL/go/ospath)                                             | Package ospath provides utilities to get OS-specific directories.                                                                                                 |
| [osutil](https://godoc.org/github.com/shurcooL/go/osutil)                                             | Package osutil offers a utility for manipulating a set of environment variables.                                                                                  |
| [parserutil](https://godoc.org/github.com/shurcooL/go/parserutil)                                     | Package parserutil offers convenience functions for parsing Go code to AST.                                                                                       |
| [pipeutil](https://godoc.org/github.com/shurcooL/go/pipeutil)                                         | Package pipeutil provides additional functionality for gopkg.in/pipe.v2 package.                                                                                  |
| [printerutil](https://godoc.org/github.com/shurcooL/go/printerutil)                                   | Package printerutil provides formatted printing of AST nodes.                                                                                                     |
| [reflectfind](https://godoc.org/github.com/shurcooL/go/reflectfind)                                   | Package reflectfind offers funcs to perform deep-search via reflect to find instances that satisfy given query.                                                   |
| [reflectsource](https://godoc.org/github.com/shurcooL/go/reflectsource)                               | Package sourcereflect implements run-time source reflection, allowing a program to look up string representation of objects from the underlying .go source files. |
| [timeutil](https://godoc.org/github.com/shurcooL/go/timeutil)                                         | Package timeutil provides a func for getting start of week of given time.                                                                                         |
| [trash](https://godoc.org/github.com/shurcooL/go/trash)                                               | Package trash implements functionality to move files into trash.                                                                                                  |
| [trim](https://godoc.org/github.com/shurcooL/go/trim)                                                 | Package trim contains helpers for trimming strings.                                                                                                               |
| [vfs/godocfs/godocfs](https://godoc.org/github.com/shurcooL/go/vfs/godocfs/godocfs)                   | Package godocfs implements vfs.FileSystem using a http.FileSystem.                                                                                                |
| [vfs/godocfs/html/vfstemplate](https://godoc.org/github.com/shurcooL/go/vfs/godocfs/html/vfstemplate) | Package vfstemplate offers html/template helpers that use vfs.FileSystem.                                                                                         |
| [vfs/godocfs/path/vfspath](https://godoc.org/github.com/shurcooL/go/vfs/godocfs/path/vfspath)         | Package vfspath implements utility routines for manipulating virtual file system paths.                                                                           |
| [vfs/godocfs/vfsutil](https://godoc.org/github.com/shurcooL/go/vfs/godocfs/vfsutil)                   | Package vfsutil implements some I/O utility functions for vfs.FileSystem.                                                                                         |

License
-------

-	[MIT License](https://opensource.org/licenses/mit-license.php)
