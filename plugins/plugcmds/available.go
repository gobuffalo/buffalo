package plugcmds

import (
	"encoding/json"
	"io"
	"os"

	"github.com/gobuffalo/buffalo/plugins"
	"github.com/gobuffalo/events"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// NewAvailable returns a fully formed Available type
func NewAvailable() *Available {
	return &Available{
		plugs: plugMap{},
	}
}

// Available used to manage all of the available commands
// for the plugin
type Available struct {
	plugs plugMap
}

type plug struct {
	BuffaloCommand string
	Cmd            *cobra.Command
	Plugin         plugins.Command
}

func (p plug) String() string {
	b, _ := json.Marshal(p.Plugin)
	return string(b)
}

// Cmd returns the "available" command
func (a *Available) Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "available",
		Short: "a list of available buffalo plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Encode(os.Stdout)
		},
	}
}

// Commands returns all of the commands that are available
func (a *Available) Commands() []*cobra.Command {
	cmds := []*cobra.Command{a.Cmd()}
	a.plugs.Range(func(_ string, p plug) bool {
		cmds = append(cmds, p.Cmd)
		return true
	})
	return cmds
}

// Add a new command to this list of available ones.
// The bufCmd should corresponding buffalo command that
// command should live below.
//
// Special "commands":
//	"root" - is the `buffalo` command
//	"events" - listens for emitted events
func (a *Available) Add(bufCmd string, cmd *cobra.Command) error {
	if len(cmd.Aliases) == 0 {
		cmd.Aliases = []string{}
	}
	p := plug{
		BuffaloCommand: bufCmd,
		Cmd:            cmd,
		Plugin: plugins.Command{
			Name:           cmd.Use,
			BuffaloCommand: bufCmd,
			Description:    cmd.Short,
			Aliases:        cmd.Aliases,
			UseCommand:     cmd.Use,
		},
	}
	a.plugs.Store(p.String(), p)
	return nil
}

// Mount all of the commands that are available
// on to the other command. This is the recommended
// approach for using Available.
//	a.Mount(rootCmd)
func (a *Available) Mount(cmd *cobra.Command) {
	// mount all the cmds on to the cobra command
	cmd.AddCommand(a.Cmd())
	a.plugs.Range(func(_ string, p plug) bool {
		cmd.AddCommand(p.Cmd)
		return true
	})
}

// Encode into the required Buffalo plugins available
// format
func (a *Available) Encode(w io.Writer) error {
	var plugs plugins.Commands
	a.plugs.Range(func(_ string, p plug) bool {
		plugs = append(plugs, p.Plugin)
		return true
	})
	return json.NewEncoder(w).Encode(plugs)
}

// Listen adds a command for github.com/gobuffalo/events.
// This will listen for ALL events. Use ListenFor to
// listen to a regex of events.
func (a *Available) Listen(fn func(e events.Event) error) error {
	return a.Add("events", buildListen(fn))
}

// ListenFor adds a command for github.com/gobuffalo/events.
// This will only listen for events that match the regex provided.
func (a *Available) ListenFor(rx string, fn func(e events.Event) error) error {
	cmd := buildListen(fn)
	p := plug{
		BuffaloCommand: "events",
		Cmd:            cmd,
		Plugin: plugins.Command{
			Name:           cmd.Use,
			BuffaloCommand: "events",
			Description:    cmd.Short,
			Aliases:        cmd.Aliases,
			UseCommand:     cmd.Use,
			ListenFor:      rx,
		},
	}
	a.plugs.Store(p.String(), p)
	return nil
}

func buildListen(fn func(e events.Event) error) *cobra.Command {
	listenCmd := &cobra.Command{
		Use:     "listen",
		Short:   "listens to github.com/gobuffalo/events",
		Aliases: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("must pass a payload")
			}

			e := events.Event{}
			err := json.Unmarshal([]byte(args[0]), &e)
			if err != nil {
				return errors.WithStack(err)
			}

			return fn(e)
		},
	}
	return listenCmd
}
