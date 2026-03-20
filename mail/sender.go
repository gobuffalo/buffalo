package mail

// Sender defines the interface for sending individual email messages.
type Sender interface {
	// Send delivers a single email message.
	Send(Message) error
}

// BatchSender defines the interface for sending multiple email messages.
type BatchSender interface {
	Sender
	// SendBatch delivers multiple messages. It returns per-message errors
	// and any general error that prevented sending entirely.
	SendBatch(messages ...Message) (errorsByMessages []error, generalError error)
}
