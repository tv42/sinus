package lineIn

import (
	"log"
	"net/url"

	"github.com/tv42/cliutil/subcommands"
	"github.com/tv42/sinus/cli"
	"github.com/tv42/sinus/util"
)

type lineInCommand struct {
	subcommands.Description
	Arguments struct{}
}

func (c *lineInCommand) Run() error {
	transport, err := cli.App.AVTransport()
	if err != nil {
		return err
	}

	u := &url.URL{
		Scheme: "x-rincon-stream",
		Opaque: string(util.DeviceUUID(transport.RootDevice)),
	}
	if err := transport.SetAVTransportURI(0, u.String(), ""); err != nil {
		log.Fatalf("line in: %v", err)
	}

	// start playback to be sure
	if err := transport.Play(0, "1"); err != nil {
		log.Fatalf("play: %v", err)
	}
	return nil
}

var lineIn = lineInCommand{
	Description: "switch to line-in",
}

func init() {
	subcommands.Register(&lineIn)
}
