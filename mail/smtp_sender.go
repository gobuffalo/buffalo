package mail

import (
	"io"
	"strconv"

	"github.com/pkg/errors"
	gomail "gopkg.in/gomail.v2"
)

//SMTPSender allows to send Emails by connecting to a SMTP server.
type SMTPSender struct {
	Dialer  *gomail.Dialer
	message *gomail.Message
}

//Send a message using SMTP configuration or returns an error if something goes wrong.
func (sm SMTPSender) Send(message Message) error {
	sm.message = gomail.NewMessage()

	sm.message.SetHeader("From", message.From)
	sm.message.SetHeader("To", message.To...)
	sm.message.SetHeader("Subject", message.Subject)
	sm.message.SetHeader("Cc", message.CC...)
	sm.message.SetHeader("Bcc", message.Bcc...)

	sm.addBodies(message)
	sm.addAttachments(message)

	for field, value := range message.Headers {
		sm.message.SetHeader(field, value)
	}

	err := sm.Dialer.DialAndSend(sm.message)

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (sm SMTPSender) addBodies(message Message) {
	if len(message.Bodies) == 0 {
		return
	}

	mainBody := message.Bodies[0]
	sm.message.SetBody(mainBody.ContentType, mainBody.Content, gomail.SetPartEncoding(gomail.Unencoded))

	for i := 1; i < len(message.Bodies); i++ {
		alt := message.Bodies[i]
		sm.message.AddAlternative(alt.ContentType, alt.Content, gomail.SetPartEncoding(gomail.Unencoded))
	}
}

func (sm SMTPSender) addAttachments(message Message) {
	for _, at := range message.Attachments {
		settings := gomail.SetCopyFunc(func(w io.Writer) error {
			if _, err := io.Copy(w, at.Reader); err != nil {
				return err
			}

			return nil
		})

		sm.message.Attach(at.Name, settings)
	}
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
