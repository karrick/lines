package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/karrick/golf"
)

var ProgramName string

func init() {
	var err error
	if ProgramName, err = os.Executable(); err != nil {
		ProgramName = os.Args[0]
	}
	ProgramName = filepath.Base(ProgramName)
}

func main() {
	if err := cmd(); err != nil {
		stderr("%s\n", err)
		if _, ok := err.(ErrUsage); ok {
			golf.Usage()
			os.Exit(2)
		}
		os.Exit(1)
	}
}

// newline returns a string with exactly one terminating newline character.
// More simple than strings.TrimRight.  When input string has multiple newline
// characters, it will strip off all but first one, reusing the same underlying
// string bytes.  When string does not end in a newline character, it returns
// the original string with a newline character appended.
func newline(s string) string {
	l := len(s)
	if l == 0 {
		return "\n"
	}

	// While this is O(length s), it stops as soon as it finds the first non
	// newline character in the string starting from the right hand side of the
	// input string.  Generally this only scans one or two characters and
	// returns.
	for i := l - 1; i >= 0; i-- {
		if s[i] != '\n' {
			if i+1 < l && s[i+1] == '\n' {
				return s[:i+2]
			}
			return s[:i+1] + "\n"
		}
	}

	return s[:1] // all newline characters, so just return the first one
}

// stderr formats and prints its arguments to standard error after prefixing
// them with the program name.
func stderr(f string, args ...interface{}) {
	os.Stderr.Write([]byte(ProgramName + ": " + newline(fmt.Sprintf(f, args...))))
}

// verbose formats and prints its arguments to standard error after prefixing
// them with the program name.  This skips printing when optVerbose is false.
func verbose(f string, args ...interface{}) {
	if *optVerbose {
		stderr(f, args...)
	}
}

// warning formats and prints its arguments to standard error after prefixing
// them with the program name.  This skips printing when optQuiet is true.
func warning(f string, args ...interface{}) {
	if !*optQuiet {
		stderr(f, args...)
	}
}

// ErrUsage is an error that a function may return when it is not invoked
// correctly.
type ErrUsage struct {
	f string
	a []interface{}
}

func NewErrUsage(f string, a ...interface{}) ErrUsage {
	return ErrUsage{f: f, a: a}
}

func (e ErrUsage) Error() string { return fmt.Sprintf(e.f, e.a...) }
