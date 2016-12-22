package cli

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	flag "github.com/spf13/pflag"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/lelandbatey/blink"
)

func Run() int {
	onlen := flag.Int64("onlen", 1000, "Length of time to turn on the lights, in milliseconds")
	delim := flag.String("delim", "\\n", "String to blink on")

	flag.Parse()

	d, err := strconv.Unquote(fmt.Sprintf("\"%s\"", *delim))
	if err == nil {
		*delim = d
	}

	duration := time.Duration(*onlen) * time.Millisecond
	if terminal.IsTerminal(syscall.Stdin) {
		err = blink.Do(duration)
	} else {
		err = blink.DoOnDelim(duration, os.Stdin, *delim)
	}

	if err != nil {
		fmt.Println(err)
		return 1
	}

	return 0
}
