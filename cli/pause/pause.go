package pause

import (
	"log"

	"github.com/tv42/cliutil/subcommands"
	"github.com/tv42/sinus/cli"
)

type pauseCommand struct {
	subcommands.Description
	Arguments struct{}
}

func (c *pauseCommand) Run() error {
	transport, err := cli.App.AVTransport()
	if err != nil {
		return err
	}

	if err := transport.Pause(0); err != nil {
		log.Fatalf("pause: %v", err)
	}
	return nil
}

var pause = pauseCommand{
	Description: "pause music",
}

func init() {
	subcommands.Register(&pause)
}
