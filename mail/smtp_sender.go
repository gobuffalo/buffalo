package mail

import (
	"fmt"
	"io"
	"strconv"

	gomail "github.com/gobuffalo/buffalo/mail/internal/mail"
)

//SMTPSender allows to send Emails by connecting to a SMTP server.
type SMTPSender struct {
	Dialer *gomail.Dialer
}

//Send a message using SMTP configuration or returns an error if something goes wrong.
func (sm SMTPSender) Send(message Message) error {
	gm := gomail.NewMessage()

	gm.SetHeader("From", message.From)
	gm.SetHeader("To", message.To...)
	gm.SetHeader("Subject", message.Subject)
	gm.SetHeader("Cc", message.CC...)
	gm.SetHeader("Bcc", message.Bcc...)

	sm.addBodies(message, gm)
	sm.addAttachments(message, gm)

	for field, value := range message.Headers {
		gm.SetHeader(field, value)
	}

	err := sm.Dialer.DialAndSend(gm)

	if err != nil {
		return err
	}

	return nil
}

func (sm SMTPSender) addBodies(message Message, gm *gomail.Message) {
	if len(message.Bodies) == 0 {
		return
	}

	mainBody := message.Bodies[0]
	gm.SetBody(mainBody.ContentType, mainBody.Content, gomail.SetPartEncoding(gomail.Unencoded))

	for i := 1; i < len(message.Bodies); i++ {
		alt := message.Bodies[i]
		gm.AddAlternative(alt.ContentType, alt.Content, gomail.SetPartEncoding(gomail.Unencoded))
	}
}

func (sm SMTPSender) addAttachments(message Message, gm *gomail.Message) {

	for _, at := range message.Attachments {
		currentAttachement := at
		settings := gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := io.Copy(w, currentAttachement.Reader)
			return err
		})

		if currentAttachement.Embedded {
			gm.Embed(currentAttachement.Name, settings)
		} else {
			gm.Attach(currentAttachement.Name, settings)
		}

	}
}

//NewSMTPSender builds a SMTP mail based in passed config.
func NewSMTPSender(host string, port string, user string, password string) (SMTPSender, error) {
	iport, err := strconv.Atoi(port)

	if err != nil {
		return SMTPSender{}, fmt.Errorf("invalid port for the SMTP mail")
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
