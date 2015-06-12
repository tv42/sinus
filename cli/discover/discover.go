package discover

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/huin/goupnp"
	"github.com/huin/goupnp/dcps/av1"
	"github.com/tv42/cliutil/flagx"
	"github.com/tv42/cliutil/positional"
	"github.com/tv42/cliutil/subcommands"
	"github.com/tv42/sinus/cli"
	"github.com/tv42/sinus/util"
)

type discoverCommand struct {
	subcommands.Description
	Arguments struct {
		positional.Optional
		Matches []flagx.Regexp `positional:",metavar=REGEXP"`
	}
}

func quote(s string) string {
	q := strconv.Quote(s)
	if len(s) > 0 && q[0] == '"' && q[len(q)-1] == '"' {
		q = q[1 : len(q)-1]
	}
	return q
}

func hostname(u *url.URL) string {
	host := u.Host
	if host == "" {
		return u.Opaque
	}
	if h, _, err := net.SplitHostPort(host); err == nil {
		return h
	}
	return host
}

func match(dev *goupnp.RootDevice, matches []flagx.Regexp) bool {
	for _, m := range matches {
		if m.FindStringIndex(dev.Device.FriendlyName) == nil &&
			m.FindStringIndex(dev.Device.Manufacturer) == nil &&
			m.FindStringIndex(dev.Device.ModelName) == nil &&
			m.FindStringIndex(dev.URLBase.String()) == nil &&
			m.FindStringIndex(util.DeviceUUID(dev)) == nil {
			return false
		}
	}
	return true
}

func (c *discoverCommand) Run() error {
	maybes, err := goupnp.DiscoverDevices(av1.URN_AVTransport_1)
	if err != nil {
		return err
	}
	if len(maybes) == 0 {
		return errors.New("no devices discovered")
	}

	devs := make(map[*goupnp.RootDevice]struct{}, len(maybes))
	for _, maybe := range maybes {
		if maybe.Err != nil {
			continue
		}
		devs[maybe.Root] = struct{}{}
	}
	if len(devs) == 0 {
		for _, maybe := range maybes {
			log.Printf("discovery error: %v", maybe.Err)
		}
		return errors.New("discovery failed")
	}

	if len(c.Arguments.Matches) > 0 {
		for dev := range devs {
			if !match(dev, c.Arguments.Matches) {
				delete(devs, dev)
			}
		}
		if len(devs) == 0 {
			return errors.New("no device matches")
		}
	}

	if len(devs) == 1 {
		// we have a match
		for dev := range devs {
			if err := cli.App.WriteConfig("location", []byte(dev.URLBase.String())); err != nil {
				return fmt.Errorf("cannot store location: %v", err)
			}
			log.Printf("remember device: %v", quote(dev.Device.FriendlyName))
		}
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	for dev := range devs {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			quote(dev.Device.FriendlyName),
			quote(dev.Device.Manufacturer),
			quote(dev.Device.ModelName),
			quote(hostname(&dev.URLBase)),
			quote(util.DeviceUUID(dev)),
		)
	}
	w.Flush()

	fmt.Println()
	fmt.Println("run `sinus discover REGEXP..` to select a player")
	return nil
}

var discover = discoverCommand{
	Description: "discover Sonos players in the network",
}

func init() {
	subcommands.Register(&discover)
}
