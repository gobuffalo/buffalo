package buffalo

// TODO: TODO-v1 check if they are really need to be exported.
/* The event id should be unique across packages as the format of
   "<package-name>:<additional-names>:<optional-error>" as documented. They
   should not be used by another packages to keep it informational. To make
   it sure, they need to be internal.
   Especially for plugable conponents like servers or workers, they can have
   their own event definition if they need but the buffalo runtime can emit
   generalize events when e.g. the runtime calls configured worker.
*/
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

	// EvtServerStart is emitted when buffalo is about to start servers
	EvtServerStart = "buffalo:server:start"
	// EvtServerStartErr is emitted when an error occurs when starting servers
	EvtServerStartErr = "buffalo:server:start:err"
	// EvtServerStop is emitted when buffalo is about to stop servers
	EvtServerStop = "buffalo:server:stop"
	// EvtServerStopErr is emitted when an error occurs when stopping servers
	EvtServerStopErr = "buffalo:server:stop:err"

	// EvtWorkerStart is emitted when buffalo is about to start workers
	EvtWorkerStart = "buffalo:worker:start"
	// EvtWorkerStartErr is emitted when an error occurs when starting workers
	EvtWorkerStartErr = "buffalo:worker:start:err"
	// EvtWorkerStop is emitted when buffalo is about to stop workers
	EvtWorkerStop = "buffalo:worker:stop"
	// EvtWorkerStopErr is emitted when an error occurs when stopping workers
	EvtWorkerStopErr = "buffalo:worker:stop:err"

	// EvtFailureErr is emitted when something can't be processed at all. it is a bad thing
	EvtFailureErr = "buffalo:failure:err"
)
