package plugins

// Command that the plugin supplies
type Command struct {
	// Name "foo"
	Name string `json:"name"`
	// BuffaloCommand "generate"
	BuffaloCommand string `json:"buffalo_command"`
	// Description "generates a foo"
	Description string `json:"description"`
	Binary      string `json:"-"`
}

// Commands is a slice of Command
type Commands []Command
