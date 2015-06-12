// Command sinus is a command-line client for Sonos audio products.
package main

import (
	"os"

	"github.com/tv42/sinus/cli"
)

//go:generate go run task/gen-imports.go -o commands.gen.go github.com/tv42/sinus/cli/...

func main() {
	code := cli.Main()
	os.Exit(code)
}
