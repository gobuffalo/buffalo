// Package openutil displays Markdown or HTML in a new browser tab.
package openutil

import (
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/shurcooL/github_flavored_markdown/gfmstyle"
	"github.com/shurcooL/go/gfmutil"
	"github.com/shurcooL/go/open"
)

// DisplayMarkdownInBrowser displays given Markdown in a new browser window/tab.
func DisplayMarkdownInBrowser(markdown []byte) {
	stopServerChan := make(chan struct{})

	handler := func(w http.ResponseWriter, req *http.Request) {
		gfmutil.WriteGitHubFlavoredMarkdownViaLocal(w, markdown)

		// TODO: A better way to fix: /assets/gfm/gfm.css Failed to load resource: net::ERR_CONNECTION_REFUSED.
		// HACK: Give some time for other assets to finish loading.
		go func() {
			time.Sleep(1 * time.Second)
			stopServerChan <- struct{}{}
		}()
	}

	http.HandleFunc("/index", handler)
	http.Handle("/assets/gfm/", http.StripPrefix("/assets/gfm", http.FileServer(gfmstyle.Assets))) // Serve the "/assets/gfm/gfm.css" file.
	http.Handle("/favicon.ico", http.NotFoundHandler())

	// TODO: Aquire a free port similarly to using ioutil.TempFile() for files.
	// TODO: Consider using httptest.NewServer.
	open.Open("http://localhost:7044/index")

	err := httpstoppable۰ListenAndServe("localhost:7044", nil, stopServerChan)
	if err != nil {
		panic(err)
	}
}

// DisplayHTMLInBrowser displays given html page in a new browser window/tab.
// query can be empty, otherwise it should begin with "?" like "?key=value".
func DisplayHTMLInBrowser(mux *http.ServeMux, stopServerChan <-chan struct{}, query string) {
	// TODO: Aquire a free port similarly to using ioutil.TempFile() for files.
	open.Open("http://localhost:7044/index" + query)

	err := httpstoppable۰ListenAndServe("localhost:7044", mux, stopServerChan)
	if err != nil {
		panic(err)
	}
}

// ListenAndServe listens on the TCP network address addr
// and then calls Serve with handler to handle requests
// on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
// Handler is typically nil, in which case the http.DefaultServeMux is
// used.
//
// When receiving from stop unblocks (because it's closed or a value is sent),
// listener is closed and ListenAndServe returns with nil error.
// Otherise, it always returns a non-nil error.
//
// Deprecated: Go 1.8 added native support for stopping a server in net/http.
// net/http should be used instead. This copied function will be removed soon.
func httpstoppable۰ListenAndServe(addr string, handler http.Handler, stop <-chan struct{}) error {
	srv := &http.Server{Addr: addr, Handler: handler}
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	go func() {
		<-stop
		err := ln.Close()
		if err != nil {
			log.Println("httpstoppable.ListenAndServe: error closing listener:", err)
		}
	}()
	err = srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
	switch { // Serve always returns a non-nil error.
	case strings.Contains(err.Error(), "use of closed network connection"):
		return nil
	default:
		return err
	}
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe so dead TCP connections
// (e.g. closing laptop mid-download) eventually go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
