package next

import (
	"log"

	"github.com/tv42/cliutil/subcommands"
	"github.com/tv42/sinus/cli"
)

type nextCommand struct {
	subcommands.Description
	Arguments struct{}
}

func (c *nextCommand) Run() error {
	transport, err := cli.App.AVTransport()
	if err != nil {
		return err
	}

	if err := transport.Next(0); err != nil {
		log.Fatalf("next: %v", err)
	}
	return nil
}

var next = nextCommand{
	Description: "go to next song",
}

func init() {
	subcommands.Register(&next)
}
