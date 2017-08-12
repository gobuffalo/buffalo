package plush

import (
	"fmt"
	"html/template"
	"log"
)

// ExampleRender using `if`, `for`, `else`, functions, etc...
func ExampleRender() {
	html := `<html>
<%= if (names && len(names) > 0) { %>
<ul>
<%= for (n) in names { %>
	<li><%= capitalize(n) %></li>
<% } %>
</ul>
<% } else { %>
	<h1>Sorry, no names. :(</h1>
<% } %>
</html>`

	ctx := NewContext()
	ctx.Set("names", []string{"john", "paul", "george", "ringo"})

	s, err := Render(html, ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(s)
	// output: <html>
	//
	// <ul>
	//
	//	<li>John</li>
	//
	//	<li>Paul</li>
	//
	//	<li>George</li>
	//
	//	<li>Ringo</li>
	//
	// </ul>
	//
	// </html>
}

func ExampleRender_scripletTags() {
	html := `<%
let h = {name: "mark"}
let greet = fn(n) {
  return "hi " + n
}
%>
<h1><%= greet(h["name"]) %></h1>`

	s, err := Render(html, NewContext())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(s)
	// output:<h1>hi mark</h1>
}

func ExampleRender_customHelperFunctions() {
	html := `<p><%= one() %></p>
<p><%= greet("mark")%></p>
<%= can("update") { %>
<p>i can update</p>
<% } %>
<%= can("destroy") { %>
<p>i can destroy</p>
<% } %>
`

	ctx := NewContext()
	ctx.Set("one", func() int {
		return 1
	})
	ctx.Set("greet", func(s string) string {
		return fmt.Sprintf("Hi %s", s)
	})
	ctx.Set("can", func(s string, help HelperContext) (template.HTML, error) {
		if s == "update" {
			h, err := help.Block()
			return template.HTML(h), err
		}
		return "", nil
	})

	s, err := Render(html, ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(s)
	// output: <p>1</p>
	// <p>Hi mark</p>
	//
	// <p>i can update</p>
}

func ExampleRender_nilValue() {
	html := `<html>
<%= if (names && len(names) > 0) { %>
	<ul>
		<%= for (n) in names { %>
			<li><%= capitalize(n) %></li>
		<% } %>
	</ul>
<% } else { %>
	<h1>Sorry, no names. :(</h1>
<% } %>
</html>`

	s, err := Render(html, NewContext())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(s)
	// output: <html>
	//
	// 	<h1>Sorry, no names. :(</h1>
	//
	// </html>
}

func ExampleRender_forIterator() {
	html := `<%= for (v) in between(3,6) { %><%=v%><% } %>`

	s, err := Render(html, NewContext())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(s)
	// output: 45
}
