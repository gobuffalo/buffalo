package mail

// Sender interface for any upcoming mailers.
type Sender interface {
	Send(Message) error
}

// BatchSender interface for sending batch or single mail
type BatchSender interface {
	Sender
	SendBatch(messages ...Message) (generalError error, errorsByMessages []error)
}
