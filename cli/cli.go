package cli

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/Wessie/appdirs"
	"github.com/huin/goupnp"
	"github.com/huin/goupnp/dcps/av1"
	"github.com/tv42/cliutil/subcommands"
)

type app struct {
	flag.FlagSet
	Config struct {
		Verbose bool
	}

	// The media player device.
	//
	// Lazily initialized, access through getDevice.
	device *goupnp.RootDevice
}

func (a *app) ReadConfig(name string) ([]byte, error) {
	return ioutil.ReadFile(path.Join(appdirs.UserConfigDir("sinus", "", "", false), name))
}

func (a *app) WriteConfig(name string, data []byte) error {
	dir := appdirs.UserConfigDir("sinus", "", "", false)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	tmp, err := ioutil.TempFile(dir, "tmp."+name+".")
	if err != nil {
		return err
	}
	closed := false
	removed := false
	defer func() {
		if !closed {
			// silence errcheck
			_ = tmp.Close()
		}
		if !removed {
			// silence errcheck
			_ = os.Remove(tmp.Name())
		}
	}()

	if _, err := tmp.Write(data); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	closed = true

	err = os.Rename(tmp.Name(), path.Join(dir, name))
	if err != nil {
		return err
	}
	removed = true

	return nil
}

func (a *app) getDevice() (*goupnp.RootDevice, error) {
	loc, err := a.ReadConfig("location")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("no device selected, run discover")
		}
		return nil, err
	}
	u, err := url.Parse(string(loc))
	if err != nil {
		log.Printf("url parsing: %v", err)

	}
	root, err := goupnp.DeviceByURL(u)
	if err != nil {
		log.Printf("error opening upnp device: %v", err)
		return nil, err
	}
	a.device = root
	return root, nil
}

func (a *app) RenderingControl() (*av1.RenderingControl1, error) {
	dev, err := a.getDevice()
	if err != nil {
		return nil, err
	}
	r, err := av1.NewRenderingControl1ClientsFromRootDevice(dev, &a.device.URLBase)
	if err != nil {
		return nil, err
	}
	return r[0], nil
}

func (a *app) AVTransport() (*av1.AVTransport1, error) {
	dev, err := a.getDevice()
	if err != nil {
		return nil, err
	}
	r, err := av1.NewAVTransport1ClientsFromRootDevice(dev, &dev.URLBase)
	if err != nil {
		return nil, err
	}
	return r[0], nil
}

// App allows command-line callables access to global flags, such as
// verbosity.
var App = app{}

func init() {
	App.BoolVar(&App.Config.Verbose, "v", false, "verbose output")

	subcommands.Register(&App)
}

func run(result subcommands.Result) (ok bool) {
	cmd := result.Command()
	run := cmd.(subcommands.Runner)
	err := run.Run()
	if err != nil {
		log.Printf("error: %v", err)
		return false
	}
	return true
}

// Main is primary entry point into the command line application.
func Main() (exitstatus int) {
	progName := filepath.Base(os.Args[0])
	log.SetFlags(0)
	log.SetPrefix(progName + ": ")

	result, err := subcommands.Parse(&App, progName, os.Args[1:])
	if err == flag.ErrHelp {
		result.Usage()
		return 0
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", result.Name(), err)
		result.Usage()
		return 2
	}

	ok := run(result)
	if !ok {
		return 1
	}
	return 0
}
