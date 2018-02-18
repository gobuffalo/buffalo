package mail

// Sender interface for any upcoming mailers.
type Sender interface {
	Send(Message) error
}
