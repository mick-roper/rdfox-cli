package main

import (
	"os"

	"github.com/mick-roper/rdfox-cli/cmd"
)

func main() {
	os.Exit(cmd.Execute())
}
