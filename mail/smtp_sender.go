package mail

import (
	"io"
	"strconv"

	"github.com/pkg/errors"
	gomail "gopkg.in/gomail.v2"
)

//SMTPSender allows to send Emails by connecting to a SMTP server.
type SMTPSender struct {
	Dialer *gomail.Dialer
}

//Send a message using SMTP configuration or returns an error if something goes wrong.
func (sm SMTPSender) Send(message Message) error {
	m := gomail.NewMessage()

	m.SetHeader("From", message.From)
	m.SetHeader("To", message.To...)
	m.SetHeader("Subject", message.Subject)
	m.SetHeader("Cc", message.CC...)
	m.SetHeader("Bcc", message.Bcc...)

	if len(message.Bodies) > 0 {
		mainBody := message.Bodies[0]
		m.SetBody(mainBody.ContentType, mainBody.Content, gomail.SetPartEncoding(gomail.Unencoded))
	}

	if len(message.Bodies) > 1 {
		for i := 1; i < len(message.Bodies); i++ {
			alt := message.Bodies[i]
			m.AddAlternative(alt.ContentType, alt.Content, gomail.SetPartEncoding(gomail.Unencoded))
		}
	}

	for _, at := range message.Attachments {
		settings := gomail.SetCopyFunc(func(w io.Writer) error {
			if _, err := io.Copy(w, at.Reader); err != nil {
				return err
			}

			return nil
		})

		m.Attach(at.Name, settings)
	}

	for field, value := range message.Headers {
		m.SetHeader(field, value)
	}

	err := sm.Dialer.DialAndSend(m)

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

//NewSMTPSender builds a SMTP mail based in passed config.
func NewSMTPSender(host string, port string, user string, password string) (SMTPSender, error) {
	iport, err := strconv.Atoi(port)

	if err != nil {
		return SMTPSender{}, errors.New("invalid port for the SMTP mail")
	}

	dialer := &gomail.Dialer{
		Host: host,
		Port: iport,
	}

	if user != "" {
		dialer.Username = user
		dialer.Password = password
	}

	return SMTPSender{
		Dialer: dialer,
	}, nil
}
