package playQueue

import (
	"log"
	"net/url"
	"strings"

	"github.com/tv42/cliutil/subcommands"
	"github.com/tv42/sinus/cli"
	"github.com/tv42/sinus/util"
)

type playQueueCommand struct {
	subcommands.Description
	Arguments struct{}
}

func (c *playQueueCommand) Run() error {
	transport, err := cli.App.AVTransport()
	if err != nil {
		return err
	}

	_, _, currentURI, _, _, _, _, _, _, err := transport.GetMediaInfo(0)
	if err != nil {
		log.Fatalf("get media info: %v", err)
	}
	if !strings.HasPrefix(currentURI, "x-rincon-queue:") {
		// not playing the queue, switch to it
		//
		// BUG: rewinds to start of queue item #0
		u := &url.URL{
			Scheme:   "x-rincon-queue",
			Opaque:   util.DeviceUUID(transport.RootDevice),
			Fragment: "0",
		}
		if err := transport.SetAVTransportURI(0, u.String(), ""); err != nil {
			log.Fatalf("play from queue: %v", err)
		}
	}

	// start playback to be sure
	if err := transport.Play(0, "1"); err != nil {
		log.Fatalf("play: %v", err)
	}
	return nil
}

var playQueue = playQueueCommand{
	Description: "play music from the queue",
}

func init() {
	subcommands.Register(&playQueue)
}
