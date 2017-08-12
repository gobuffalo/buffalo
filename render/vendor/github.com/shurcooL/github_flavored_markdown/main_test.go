package github_flavored_markdown_test

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/shurcooL/github_flavored_markdown"
	"github.com/shurcooL/github_flavored_markdown/gfmstyle"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ExampleMarkdown() {
	text := []byte("Hello world github/linguist#1 **cool**, and #1!")

	os.Stdout.Write(github_flavored_markdown.Markdown(text))

	// Output:
	// <p>Hello world github/linguist#1 <strong>cool</strong>, and #1!</p>
}

// An example of how to generate a complete HTML page, including CSS styles.
func ExampleMarkdown_completeHTMLPage() {
	// Serve the "/assets/gfm.css" file.
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(gfmstyle.Assets)))

	var w io.Writer = os.Stdout // It can be an http.ResponseWriter.
	markdown := []byte("# GitHub Flavored Markdown\n\nHello.")

	io.WriteString(w, `<html><head><meta charset="utf-8"><link href="/assets/gfm.css" media="all" rel="stylesheet" type="text/css" /><link href="//cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)
	w.Write(github_flavored_markdown.Markdown(markdown))
	io.WriteString(w, `</article></body></html>`)

	// Output:
	// <html><head><meta charset="utf-8"><link href="/assets/gfm.css" media="all" rel="stylesheet" type="text/css" /><link href="//cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;"><h1><a name="github-flavored-markdown" class="anchor" href="#github-flavored-markdown" rel="nofollow" aria-hidden="true"><span class="octicon octicon-link"></span></a>GitHub Flavored Markdown</h1>
	//
	// <p>Hello.</p>
	// </article></body></html>
}

func TestComponents(t *testing.T) {
	tests := []struct {
		text string
		want string
	}{
		{
			// Heading.
			text: "## git diff",
			want: `<h2><a name="git-diff" class="anchor" href="#git-diff" rel="nofollow" aria-hidden="true"><span class="octicon octicon-link"></span></a>git diff</h2>` + "\n",
		},
		{
			// Heading Link.
			text: "### [Some **bold** _italic_ link](http://www.example.com)",
			want: `<h3><a name="some-bold-italic-link" class="anchor" href="#some-bold-italic-link" rel="nofollow" aria-hidden="true"><span class="octicon octicon-link"></span></a><a href="http://www.example.com" rel="nofollow">Some <strong>bold</strong> <em>italic</em> link</a></h3>` + "\n",
		},
		{
			// Task List.
			text: `- [ ] This is an incomplete task.
- [x] This is done.
`,
			want: `<ul>
<li><input type="checkbox" disabled=""> This is an incomplete task.</li>
<li><input type="checkbox" checked="" disabled=""> This is done.</li>
</ul>
`,
		},
		{
			// No need to insert an empty line to start a (code, quote, ordered list, unordered list) block.
			// See issue https://github.com/shurcooL/github_flavored_markdown/issues/9.
			text: `Some text
- A
- B
- C`,
			want: `<p>Some text</p>

<ul>
<li>A</li>
<li>B</li>
<li>C</li>
</ul>
`,
		},
	}

	for _, test := range tests {
		if got := string(github_flavored_markdown.Markdown([]byte(test.text))); got != test.want {
			t.Errorf("\ngot %q\nwant %q", got, test.want)
		}
	}
}

func ExampleHeading() {
	heading := github_flavored_markdown.Heading(atom.H2, "Hello > Goodbye")
	html.Render(os.Stdout, heading)

	// Output:
	// <h2><a name="hello-goodbye" class="anchor" href="#hello-goodbye" rel="nofollow" aria-hidden="true"><span class="octicon-link"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16" style="fill: currentColor; vertical-align: top;"><path d="M4 9h1v1H4c-1.5 0-3-1.69-3-3.5S2.55 3 4 3h4c1.45 0 3 1.69 3 3.5 0 1.41-.91 2.72-2 3.25V8.59c.58-.45 1-1.27 1-2.09C10 5.22 8.98 4 8 4H4c-.98 0-2 1.22-2 2.5S3 9 4 9zm9-3h-1v1h1c1 0 2 1.22 2 2.5S13.98 12 13 12H9c-.98 0-2-1.22-2-2.5 0-.83.42-1.64 1-2.09V6.25c-1.09.53-2 1.84-2 3.25C6 11.31 7.55 13 9 13h4c1.45 0 3-1.69 3-3.5S14.5 6 13 6z"></path></svg></span></a>Hello &gt; Goodbye</h2>
}
