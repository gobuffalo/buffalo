package worker

const (
	// EvtWorkerStart is emitted when buffalo.App#Serve is called and workers are started
	EvtWorkerStart = "buffalo:worker:start"
	// EvtWorkerStartErr is emitted when an error occurs when starting workers
	EvtWorkerStartErr = "buffalo:worker:start:err"

	// EvtWorkerStop is emitted when buffalo.App#Stop is called and workers are stopped
	EvtWorkerStop = "buffalo:worker:stop"
	// EvtWorkerStopErr is emitted when an error occurs when stopping workers
	EvtWorkerStopErr = "buffalo:worker:stop:err"
)
