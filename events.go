package buffalo

const (
	// EvtAppStart is emitted when buffalo.App#Serve is called
	EvtAppStart = "buffalo:app:start"
	// EvtAppStartErr is emitted when an error occurs calling buffalo.App#Serve
	EvtAppStartErr = "buffalo:app:start:err"

	// EvtAppStop is emitted when buffalo.App#Stop is called
	EvtAppStop = "buffalo:app:stop"
	// EvtAppStopErr is emitted when an error occurs calling buffalo.App#Stop
	EvtAppStopErr = "buffalo:app:stop:err"

	// EvtRouteStarted is emitted when a requested route is being processed
	EvtRouteStarted = "buffalo:route:started"
	// EvtRouteFinished is emitted when a requested route is completed
	EvtRouteFinished = "buffalo:route:finished"
	// EvtRouteErr is emitted when there is a problem handling processing a route
	EvtRouteErr = "buffalo:route:err"

	// EvtWorkerStart is emitted when buffalo.App#Serve is called and workers are started
	EvtWorkerStart = "buffalo:worker:start"
	// EvtWorkerStartErr is emitted when an error occurs when starting workers
	EvtWorkerStartErr = "buffalo:worker:start:err"

	// EvtWorkerStop is emitted when buffalo.App#Stop is called and workers are stopped
	EvtWorkerStop = "buffalo:worker:stop"
	// EvtWorkerStopErr is emitted when an error occurs when stopping workers
	EvtWorkerStopErr = "buffalo:worker:stop:err"

	// EvtFailureErr is emitted when something can't be processed at all. it is a bad thing
	EvtFailureErr = "buffalo:failure:err"
)
