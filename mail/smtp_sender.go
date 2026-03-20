package mail

import (
	"fmt"
	"io"
	"strconv"
)

// SMTPSender delivers emails via an SMTP server.
type SMTPSender struct {
	// Dialer configures the connection to the SMTP server.
	Dialer *Dialer
}

// Send delivers a single message via SMTP.
func (sm SMTPSender) Send(message Message) error {
	return sm.Dialer.dialAndSend(sm.prepareMessage(message))
}

// SendBatch delivers multiple messages using a single SMTP connection.
// Returns per-message errors and any general connection error.
func (sm SMTPSender) SendBatch(messages ...Message) (errorsByMessages []error, generalError error) {
	preparedMessages := make([]*smtpMessage, len(messages))
	for i, message := range messages {
		preparedMessages[i] = sm.prepareMessage(message)
	}

	s, err := sm.Dialer.dial()
	if err != nil {
		return nil, err
	}
	defer s.Close()

	return sendSMTP(s, preparedMessages...), nil
}
func (sm SMTPSender) prepareMessage(message Message) *smtpMessage {
	gm := newSMTPMessage()

	gm.setHeader("From", message.From)
	gm.setHeader("To", message.To...)
	gm.setHeader("Subject", message.Subject)
	gm.setHeader("Cc", message.CC...)
	gm.setHeader("Bcc", message.Bcc...)

	sm.addBodies(message, gm)
	sm.addAttachments(message, gm)

	for field, value := range message.Headers {
		gm.setHeader(field, value)
	}

	return gm
}

func (sm SMTPSender) addBodies(message Message, gm *smtpMessage) {
	if len(message.Bodies) == 0 {
		return
	}

	mainBody := message.Bodies[0]
	gm.setBody(mainBody.ContentType, mainBody.Content, setPartEncoding(encodingUnencoded))

	for i := 1; i < len(message.Bodies); i++ {
		alt := message.Bodies[i]
		gm.addAlternative(alt.ContentType, alt.Content, setPartEncoding(encodingUnencoded))
	}
}

func (sm SMTPSender) addAttachments(message Message, gm *smtpMessage) {

	for _, at := range message.Attachments {
		currentAttachement := at
		settings := setCopyFunc(func(w io.Writer) error {
			_, err := io.Copy(w, currentAttachement.Reader)
			return err
		})

		if currentAttachement.Embedded {
			gm.embed(currentAttachement.Name, settings)
		} else {
			gm.attach(currentAttachement.Name, settings)
		}

	}
}

// NewSMTPSender builds a SMTP mail based in passed config.
func NewSMTPSender(host string, port string, user string, password string) (SMTPSender, error) {
	iport, err := strconv.Atoi(port)

	if err != nil {
		return SMTPSender{}, fmt.Errorf("invalid port for the SMTP mail")
	}

	dialer := &Dialer{
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
