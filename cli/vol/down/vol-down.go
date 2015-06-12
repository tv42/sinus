package volDown

import (
	"log"

	"github.com/tv42/cliutil/subcommands"
	"github.com/tv42/sinus/cli"
)

type volDownCommand struct {
	subcommands.Description
	Arguments struct{}
}

func (c *volDownCommand) Run() error {
	control, err := cli.App.RenderingControl()
	if err != nil {
		return err
	}

	old, err := control.GetVolume(0, "Master")
	if err != nil {
		log.Fatalf("get volume: %v", err)
	}
	// TODO need a relative step, but sonos doesn't support GetVolumeDBRange
	vol := old - 5
	if vol > old {
		// wraparound, go back to min value
		vol = 0
	}
	if err := control.SetVolume(0, "Master", vol); err != nil {
		log.Fatalf("get volume: %v", err)
	}
	return nil
}

var volDown = volDownCommand{
	Description: "decrease volume",
}

func init() {
	subcommands.Register(&volDown)
}
