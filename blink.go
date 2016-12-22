// +build linux
package main

import (
	"os"
	"syscall"
	"time"

	flag "github.com/spf13/pflag"
)

const device = "/dev/console"

// Thanks Dave Cheney, what a guy!:
//     https://github.com/davecheney/pcap/blob/10760a170da6335ec1a48be06a86f494b0ef74ab/bpf.go#L45
func ioctl(fd int, request, argp uintptr) error {
	_, _, errorp := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), request, argp)
	return os.NewSyscallError("ioctl", errorp)
}

// blink will turn on the keyboard lights for the given amount of time. Yes ALL
// the keyboard lights.
func blink(onLen time.Duration) {
	// ya this is probably not safe, cause I ported this to Go from Python
	// using four year old go code about how to make ioctl calls in go (btw the
	// below code is probably SUPER unsafe).
	console_fd, e := syscall.Open(device, os.O_RDONLY|syscall.O_CLOEXEC, 0666)
	if e != nil {
		panic(e)
	}

	KDSETLED := 0x4B32

	SCR_LED := 0x01
	NUM_LED := 0x02
	CAP_LED := 0x04

	all_on := SCR_LED | NUM_LED | CAP_LED
	all_off := 0
	ioctl(console_fd, uintptr(KDSETLED), uintptr(all_on))
	time.Sleep(onLen)
	ioctl(console_fd, uintptr(KDSETLED), uintptr(all_off))
}

func main() {
	onlen := flag.Int64("onlen", 1000, "Length of time to turn on the lights, in milliseconds")
	flag.Parse()

	dura := time.Duration(*onlen) * time.Millisecond
	blink(dura)
}
