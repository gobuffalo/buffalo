package build

import (
	"github.com/gobuffalo/buffalo-cli/genny/build"
)

const (
	// EvtBuildStart is emitted when building starts
	EvtBuildStart = build.EvtBuildStart
	// EvtBuildStop is emitted when building stops
	EvtBuildStop = build.EvtBuildStop
	// EvtBuildStopErr is emitted when building is stopped due to an error
	EvtBuildStopErr = build.EvtBuildStopErr
)

var New = build.New
var Cleanup = build.Cleanup

type Options = build.Options
