package mail

// Sender interface for any upcoming mailers.
type Sender interface {
	Send(Message) error
}

// Sender interface for any upcoming mailers.
type BatchSender interface {
	Sender
	SendBatch(messages ...Message) (generalError error, errorsByMessages []error)
}