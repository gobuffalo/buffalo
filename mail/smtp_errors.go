// Portions of this code are derived from the go-mail/mail project.
// https://github.com/go-mail/mail (MIT License)

package mail

import "fmt"

// sendError represents the failure to transmit a Message, detailing the cause
// of the failure and index of the Message within a batch.
type sendError struct {
	// Index specifies the index of the Message within a batch.
	Index uint
	Cause error
}

func (err *sendError) Error() string {
	return fmt.Sprintf("gomail: could not send email %d: %v",
		err.Index+1, err.Cause)
}
