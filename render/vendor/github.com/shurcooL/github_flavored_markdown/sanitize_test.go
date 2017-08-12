package github_flavored_markdown

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"testing"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// In this test, nothing should be sanitized away.
func TestSanitizeLarge(t *testing.T) {
	var text = []byte(`### GitHub Flavored Markdown rendered locally using go gettable native Go code

` + "```Go" + `
package main

import "fmt"

func main() {
	// This is a comment!
	/* so is this */
	fmt.Println("Hello, playground", 123, 1.336)
}
` + "```" + `

` + "```diff" + `
diff --git a/main.go b/main.go
index dc83bf7..5260a7d 100644
--- a/main.go
+++ b/main.go
@@ -1323,10 +1323,10 @@ func (this *GoPackageSelecterAdapter) GetSelectedGoPackage() *GoPackage {
 }

 // TODO: Move to the right place.
-var goPackages = &exp14.GoPackages{SkipGoroot: false}
+var goPackages = &exp14.GoPackages{SkipGoroot: true}

 func NewGoPackageListingWidget(pos, size mathgl.Vec2d) *SearchableListWidget {
 	goPackagesSliceStringer := &goPackagesSliceStringer{goPackages}
` + "```" + `
`)

	htmlFlags := 0
	renderer := &renderer{Html: blackfriday.HtmlRenderer(htmlFlags, "", "").(*blackfriday.Html)}

	unsanitized := blackfriday.Markdown(text, renderer, extensions)

	// GitHub Flavored Markdown-like sanitization policy.
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").Matching(bluemonday.SpaceSeparatedTokens).OnElements("div", "span")
	p.AllowAttrs("class", "name").Matching(bluemonday.SpaceSeparatedTokens).OnElements("a")
	p.AllowAttrs("rel").Matching(regexp.MustCompile(`^nofollow$`)).OnElements("a")
	p.AllowAttrs("aria-hidden").Matching(regexp.MustCompile(`^true$`)).OnElements("a")
	p.AllowDataURIImages()

	output := p.SanitizeBytes(unsanitized)

	diff, err := diff(unsanitized, output)
	if err != nil {
		log.Fatalln(err)
	}

	if len(diff) != 0 {
		t.Errorf("Difference of %d lines:\n%s", bytes.Count(diff, []byte("\n")), string(diff))
	}
}

func TestSanitize(t *testing.T) {
	tests := []struct {
		text string
		want string
	}{
		{
			// Make sure that <script> tag is sanitized away.
			text: "Hello <script>alert();</script> world.",
			want: "<p>Hello  world.</p>\n",
		},
		{
			// Make sure that "class" attribute values that are not sane get sanitized away.
			// Just a normal class name, should be preserved.
			text: `Hello <span class="foo bar bash">there</span> world.`,
			want: `<p>Hello <span class="foo bar bash">there</span> world.</p>` + "\n",
		},
		{
			// JavaScript in class name, should be sanitized away.
			text: `Hello <span class="javascript:alert('XSS')">there</span> world.`,
			want: "<p>Hello <span>there</span> world.</p>" + "\n",
		},
		{
			// Script injection attempt, should be sanitized away.
			text: `Hello <span class="><script src='http://hackers.org/XSS.js'></script>">there</span> world.`,
			want: "<p>Hello ",
		},
	}

	for _, test := range tests {
		if got := string(Markdown([]byte(test.text))); got != test.want {
			t.Errorf("\ngot %q\nwant %q", got, test.want)
		}
	}
}

func TestSanitizeAnchorName(t *testing.T) {
	tests := []struct {
		text string
		want string
	}{
		{
			text: "## Did you just steal this template from Tom's TOML?",
			want: `<h2><a name="did-you-just-steal-this-template-from-tom-s-toml" class="anchor" href="#did-you-just-steal-this-template-from-tom-s-toml" rel="nofollow" aria-hidden="true"><span class="octicon octicon-link"></span></a>Did you just steal this template from Tom&#39;s TOML?</h2>` + "\n",
		},
		{
			text: `## What about "quotes" & things?`,
			want: `<h2><a name="what-about-quotes-things" class="anchor" href="#what-about-quotes-things" rel="nofollow" aria-hidden="true"><span class="octicon octicon-link"></span></a>What about &#34;quotes&#34; &amp; things?</h2>` + "\n",
		},
	}

	for _, test := range tests {
		if got := string(Markdown([]byte(test.text))); got != test.want {
			t.Errorf("\ngot %q\nwant %q", got, test.want)
		}
	}
}

// TODO: Factor out.
func diff(b1, b2 []byte) (data []byte, err error) {
	f1, err := ioutil.TempFile("", "")
	if err != nil {
		return
	}
	defer os.Remove(f1.Name())
	defer f1.Close()

	f2, err := ioutil.TempFile("", "")
	if err != nil {
		return
	}
	defer os.Remove(f2.Name())
	defer f2.Close()

	f1.Write(b1)
	f2.Write(b2)

	data, err = exec.Command("diff", "-u", f1.Name(), f2.Name()).CombinedOutput()
	if len(data) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil
	}
	return
}
