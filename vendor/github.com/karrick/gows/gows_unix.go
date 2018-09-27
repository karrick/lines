// +build darwin dragonfly freebsd netbsd openbsd linux

package gows

import (
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

func getwinsize() (int, int, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return 0, 0, err
	}
	ws, err := unix.IoctlGetWinsize(int(tty.Fd()), syscall.TIOCGWINSZ)
	err2 := tty.Close()
	if err != nil {
		return 0, 0, err
	}
	return int(ws.Col), int(ws.Row), err2
}
