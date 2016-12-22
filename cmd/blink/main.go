package main

import (
	"os"

	"github.com/lelandbatey/blink/cmd/cli"
)

func main() {
	os.Exit(cli.Run())
}
