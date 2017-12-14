package mail

// Sender interface for any upcomming mailers.
type Sender interface {
	Send(Message) error
}
