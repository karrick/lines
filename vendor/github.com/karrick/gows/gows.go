package gows

// GetWinSize returns the size of the console window. It returns the number of
// columns, number of rows, or whatever error was encountered while attempting
// to get the requested window size. On POSIX it attempts to open the associated
// controlling terminal from /dev/tty. On Windows it invokes
// GetConsoleScreenBufferInfo for the process' standard output stream.
func GetWinSize() (int, int, error) {
	return getwinsize()
}
