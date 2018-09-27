package golf

import (
	"fmt"
	"os"
	"path/filepath"
	"unicode/utf8"
)

// Arg returns the i'th command-line argument. Arg(0) is the first remaining
// argument after flags have been processed. Arg returns an empty string if the
// requested element does not exist.
func Arg(i int) string {
	if i < len(remainingArguments) {
		return remainingArguments[i]
	}
	return ""
}

// Args returns the non-flag command-line arguments.
func Args() []string {
	return remainingArguments
}

// NArg returns the number of arguments remaining after flags have been
// processed.
func NArg() int {
	return len(remainingArguments)
}

// NFlag returns the number of command-line flags that have been set.
func NFlag() int {
	return argsProcessed
}

// PrintDefaults prints to standard error, a usage message showing the default
// settings of all defined command-line flags.
func PrintDefaults() {
	for _, opt := range flags {
		var typeName string
		description := opt.Description()
		value := opt.Default()

		switch value.(type) {
		case bool:
			// do not want to add a default blob when boolean
		default:
			typeName = fmt.Sprintf(" %T", value)
		}

		short := opt.Short()
		long := opt.Long()

		if short != utf8.RuneError {
			if long != "" {
				fmt.Fprintf(os.Stderr, "  -%c, --%s%s\n", short, long, typeName)
			} else {
				fmt.Fprintf(os.Stderr, "  -%c%s\n", short, typeName)
			}
		} else {
			fmt.Fprintf(os.Stderr, "  --%s%s\n", long, typeName)
		}

		d := "%v"
		switch value.(type) {
		case string, rune:
			d = "%q"
		}

		if description != "" {
			fmt.Fprintf(os.Stderr, "\t%s (default: "+d+")\n", description, value)
		} else {
			fmt.Fprintf(os.Stderr, "\t(default: "+d+")\n", value)
		}
	}
}

// Usage prints command line usage to stderr.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", filepath.Base(os.Args[0]))
	PrintDefaults()
}
