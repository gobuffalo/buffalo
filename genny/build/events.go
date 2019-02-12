package build

const (
	// EvtBuildStart is emitted when building starts
	EvtBuildStart = "buffalo:build:start"
	// EvtBuildStop is emitted when building stops
	EvtBuildStop = "buffalo:build:stop"
	// EvtBuildStopErr is emitted when building is stopped due to an error
	EvtBuildStopErr = "buffalo:build:stop:err"
)
