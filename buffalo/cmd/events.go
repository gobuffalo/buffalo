package cmd

const (
	// EvtSetupStarted is emitted when `buffalo setup` starts
	EvtSetupStarted = "buffalo:setup:started"
	// EvtSetupErr is emitted if `buffalo setup` fails
	EvtSetupErr = "buffalo:setup:err"
	// EvtSetupFinished is emitted when `buffalo setup` finishes
	EvtSetupFinished = "buffalo:setup:finished"
)
