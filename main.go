package main

import (
	"os"

	"github.com/mick-roper/rdfox-cli/cmd"
)

var Version string

func main() {
	os.Exit(cmd.Execute(Version))
}
