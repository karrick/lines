package gows

import (
	"golang.org/x/sys/windows"
)

func getwinsize() (int, int, error) {
	var csbi windows.ConsoleScreenBufferInfo
	err := windows.GetConsoleScreenBufferInfo(windows.Stdout, &csbi)
	if err != nil {
		return 0, 0, err
	}
	return int(csbi.MaximumWindowSize.X), int(csbi.MaximumWindowSize.Y), nil
}
