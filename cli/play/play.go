package play

import (
	"log"

	"github.com/tv42/cliutil/subcommands"
	"github.com/tv42/sinus/cli"
)

type playCommand struct {
	subcommands.Description
}

func (c *playCommand) Run() error {
	transport, err := cli.App.AVTransport()
	if err != nil {
		return err
	}

	if err := transport.Play(0, "1"); err != nil {
		log.Fatalf("play: %v", err)
	}
	return nil
}

var play = playCommand{
	Description: "start playing current audio source",
}

func init() {
	subcommands.Register(&play)
}
