package plugins

import (
	"context"

	"github.com/gobuffalo/buffalo/genny/plugins/install"
	"github.com/gobuffalo/buffalo/plugins"
	"github.com/gobuffalo/events"
	"github.com/gobuffalo/genny/v2"
)

// Listen is listener for plugin events pipeline
func Listen(e events.Event) error {
	if e.Kind != "buffalo:setup:started" {
		return nil
	}

	run := genny.WetRunner(context.Background())

	opts := &install.Options{}
	gg, err := install.New(opts)
	if err != nil {
		return err
	}
	run.WithGroup(gg)
	payload := e.Payload
	payload["plugins"] = opts.Plugins
	events.EmitPayload(plugins.EvtSetupStarted, payload)
	if err := run.Run(); err != nil {
		events.EmitError(plugins.EvtSetupErr, err, payload)
		return err
	}
	events.EmitPayload(plugins.EvtSetupFinished, payload)
	return nil
}
