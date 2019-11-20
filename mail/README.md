# github.com/gobuffalo/buffalo/mail

This package is intended to allow easy Email sending with Buffalo, it allows you to define your custom `mail.Sender` for the provider you would like to use.

## Generator

```bash
$ buffalo generate mailer welcome_email
```

## Example Usage

```go
//actions/mail.go
package x

import (
	"log"
	"net/http"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/buffalo/mail"
	"errors"
	"gitlab.com/wawandco/app/models"
)

var smtp mail.Sender
var r *render.Engine

func init() {

	//Pulling config from the env.
	port := envy.Get("SMTP_PORT", "1025")
	host := envy.Get("SMTP_HOST", "localhost")
	user := envy.Get("SMTP_USER", "")
	password := envy.Get("SMTP_PASSWORD", "")

	var err error
	smtp, err = mail.NewSMTPSender(host, port, user, password)

	if err != nil {
		log.Fatal(err)
	}

	//The rendering engine, this is usually generated inside actions/render.go in your buffalo app.
	r = render.New(render.Options{
		TemplatesBox:   packr.New("app:mail", "../templates"),
	})
}

//SendContactMessage Sends contact message to contact@myapp.com
func SendContactMessage(c *models.Contact) error {

	//Creates a new message
	m := mail.NewMessage()
	m.From = "sender@myapp.com"
	m.Subject = "New Contact"
	m.To = []string{"contact@myapp.com"}

	// Data that will be used inside the templates when rendering.
	data := map[string]interface{}{
		"contact": c,
	}

	// You can add multiple bodies to the message you're creating to have content-types alternatives.
	err := m.AddBodies(data, r.HTML("mail/contact.html"), r.Plain("mail/contact.txt"))

	if err != nil {
		return err
	}

	err = smtp.Send(m)
	if err != nil {
		return err
	}

	return nil
}

```

This `SendContactMessage` could be called by one of your actions, p.e. the action that handles your contact form submission.

```go
//actions/contact.go
...

func ContactFormHandler(c buffalo.Context) error {
    contact := &models.Contact{}
    c.Bind(contact)

    //Calling to send the message
    SendContactMessage(contact)
    return c.Redirect(http.StatusFound, "contact/thanks")
}
...
```

If you're using Gmail or need to configure your SMTP connection you can use the Dialer property on the SMTPSender, p.e: (for Gmail)

```go
...
var smtp mail.Sender

func init() {
    port := envy.Get("SMTP_PORT", "465")
    // or 587 with TLS

	host := envy.Get("SMTP_HOST", "smtp.gmail.com")
	user := envy.Get("SMTP_USER", "your@email.com")
	password := envy.Get("SMTP_PASSWORD", "yourp4ssw0rd")

	var err error
	sender, err := mail.NewSMTPSender(host, port, user, password)
	sender.Dialer.SSL = true

    //or if TLS
    sender.Dialer.TLSConfig = &tls.Config{...}

    smtp = sender
}
...
```
