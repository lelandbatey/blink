// +build linux
package blink

import (
	"io"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

const device = "/dev/console"

// ioctl is a helper function for making an ioctl call using Go's syscall package.
// Thanks Dave Cheney, what a guy!:
//     https://github.com/davecheney/pcap/blob/10760a170da6335ec1a48be06a86f494b0ef74ab/bpf.go#L45
func ioctl(fd int, request, argp uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), request, argp)
	if errno != 0 {
		return os.NewSyscallError("ioctl", errno)
	}
	return nil
}

// Do will turn on the keyboard lights for the given amount of time. Yes ALL
// the keyboard lights.
func Do(onLen time.Duration) error {
	// This is probably not safe. I ported this to Go from Python using four
	// year old Go code about how to make ioctl calls in Go
	console_fd, err := syscall.Open(device, os.O_RDONLY|syscall.O_CLOEXEC, 0666)
	defer func() {
		if err := syscall.Close(console_fd); err != nil {
			log.Printf("Failed to close file descriptor for /dev/console, fd %v", console_fd)
		}
	}()
	if err != nil {
		return errors.Wrapf(err, "cannot open %q using syscall \"O_RDONLY|O_CLOEXEC 0666\"", device)
	}

	// KDSETLED is an ioctl argument for manually changing the state of
	// keyboard LEDs. You can find an excellent example of how it's used, with
	// further references, here:
	//     http://www.tldp.org/LDP/lkmpg/2.6/html/x1194.html
	KDSETLED := 0x4B32

	// These values are defined in 'include/uapi/linux/kd.h' of the Linux
	// kernel source.
	SCR_LED := 0x01
	NUM_LED := 0x02
	CAP_LED := 0x04

	all_on := SCR_LED | NUM_LED | CAP_LED
	// restore will restore the previous value of the keyboard lights. Must be
	// a value higher than 7, so we choose 0xFF.
	restore := 0xFF
	if err := ioctl(console_fd, uintptr(KDSETLED), uintptr(all_on)); err != nil {
		return err
	}
	time.Sleep(onLen)
	if err = ioctl(console_fd, uintptr(KDSETLED), uintptr(restore)); err != nil {
		return err
	}

	return nil
}

// DoOnDelim will call blink for duration every time a delimiter is read on
// the reader and will not blink for at least that duration.
func DoOnDelim(duration time.Duration, r io.Reader, delimiter string) error {
	delim := []byte(delimiter)
	dpos := 0
	buf := make([]byte, 1)
	for {
		_, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "cannot continue reading input")
		}
		if buf[0] == delim[dpos] {
			// We found the delimiter guys, do the blink!
			if dpos == len(delim)-1 {
				dpos = 0
				if err := Do(duration); err != nil {
					return err
				}
				time.Sleep(duration)
			} else {
				dpos += 1
			}
		} else {
			dpos = 0
		}
	}

	return nil
}
