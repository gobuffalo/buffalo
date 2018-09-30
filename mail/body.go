package mail

// Body represents one of the bodies in the Message could be main or alternative
type Body struct {
	Content     string
	ContentType string
}
