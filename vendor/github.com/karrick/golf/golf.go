package golf

import (
	"fmt"
	"io"
	"os"
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
	PrintDefaultsTo(os.Stderr)
}

// PrintDefaultsTo prints to w, a usage message showing the default settings of
// all defined command-line flags.
func PrintDefaultsTo(w io.Writer) {
	for _, opt := range flags {
		var def, typeName string
		description := opt.Description()
		value := opt.Default()

		switch value.(type) {
		case bool:
			typeName = "" // do not want to add a default blob when boolean
		case string, rune:
			def = fmt.Sprintf(" (default: %q)", value)
			typeName = fmt.Sprintf(" %T", value)
		default:
			def = fmt.Sprintf(" (default: %v)", value)
			typeName = fmt.Sprintf(" %T", value)
		}

		short := opt.Short()
		long := opt.Long()

		if short != utf8.RuneError {
			if long != "" {
				fmt.Fprintf(w, "  -%c, --%s%s%s\n", short, long, typeName, def)
			} else {
				fmt.Fprintf(w, "  -%c%s%s\n", short, typeName, def)
			}
		} else {
			fmt.Fprintf(w, "  --%s%s%s\n", long, typeName, def)
		}

		if description != "" {
			fmt.Fprint(w, dw.Wrap(fmt.Sprintf("%s", description)))
		}
	}
}

// dw is default wrapper, used for printing default options.
var dw = LineWrapper{Max: 80, Prefix: "    "}

// Usage prints command line usage to stderr, but may be overridden by programs
// that need to customize the usage information.
var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	PrintDefaults()
}
