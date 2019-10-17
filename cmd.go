package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/karrick/gobls"
	"github.com/karrick/golf"
	"github.com/karrick/gotb"
)

func init() {
	// Rather than display the entire usage information for a parsing error,
	// merely allow golf library to display the error message, then print the
	// command the user may use to show command line usage information.
	golf.Usage = func() {
		fmt.Fprintf(os.Stderr, "Use `%s --help` for more information.\n", ProgramName)
	}
}

var (
	optHelp    = golf.BoolP('h', "help", false, "Print command line help then exit.")
	optQuiet   = golf.BoolP('q', "quiet", false, "Do not print intermediate errors to stderr.")
	optVerbose = golf.BoolP('v', "verbose", false, "Print verbose output to stderr.")
	optForce   = golf.Bool("force", false, "Print error messages but continue processing.")

	optRange      = golf.StringP('r', "range", "", "Only print lines START-END.")
	optSkipTop    = golf.Uint("skip-top", 0, "Skip printing the top N header lines.")
	optSkipBottom = golf.Uint("skip-bottom", 0, "Skip printing the bottom N footer lines.")
	optTop        = golf.UintP('t', "top", 0, "Only print the top N lines.")
	optBottom     = golf.UintP('b', "bottom", 0, "Only print the bottom N lines.")
)

func cmd() error {
	golf.Parse()

	if *optHelp {
		fmt.Println(golf.Wrap("SUMMARY:  lines [options] [file1 [file2]] [options]"))
		fmt.Println(golf.Wrap("Without command line arguments, reads from standard input and writes to standard output. With command line arguments, reads from each file in sequence, and applies the below logic independently for each file."))
		fmt.Println(golf.Wrap("When given the '--range N' command line argument, prints the line number corresponding to N. When given the '--range START-END' command line argument, prints lines 'START' thru 'END', inclusively. START must not be greater than the value of END. When START is omitted, the first line printed will be the first line of the input. When END is omitted, the final line printed will be the final line of the input."))
		fmt.Println(golf.Wrap("When given the '--skip-top N' command line argument, omits printing the initial N lines, handy for removing a possibly multiline header from some text."))
		fmt.Println(golf.Wrap("When given the '--skip-bottom N' command line argument, omits printing the final N lines, handy for removing a possibly multiline footer from some text."))
		fmt.Println(golf.Wrap("When given the '--top N' command line argument, prints only the initial N lines, similar to the behavior of 'head -n N', but included in this tool for completeness."))
		fmt.Println(golf.Wrap("When given the '--bottom N' command line argument, prints only the final N lines, similar to the behavior of 'tail -n N', but included in this tool for completeness."))
		fmt.Println(golf.Wrap("USAGE:    Not all options may be used with all other options. See below synopsis for reference."))

		fmt.Println("\tlines [\t--top N | --bottom N |\n\t\t" + strings.Join([]string{
			"\t--range M-N | --range M- | --range -N | --range N |",
			"\t--skip-top N | --skip-bottom N ]",
			"\t[file1 [file2...]]",
		}, "\n\t\t") + "\n")
		fmt.Println("EXAMPLES:")
		fmt.Println("\tlines < sample.txt")
		fmt.Println("\tlines sample.txt")
		fmt.Println("\tlines sample.txt --range 4-7")
		fmt.Println("\tlines sample.txt --range -3")
		fmt.Println("\tlines sample.txt --range 7-")
		fmt.Println("\tlines sample.txt --range 3")
		fmt.Println("\tlines sample.txt --skip-top 2")
		fmt.Println("\tlines sample.txt --skip-bottom 2")
		fmt.Println("\tlines sample.txt --skip-top 3 --skip-bottom 2")
		fmt.Println("\tlines sample.txt --top 3")
		fmt.Println("\tlines sample.txt --bottom 3")
		fmt.Println("\nCommand line options:")
		golf.PrintDefaults()
		return nil
	}

	if *optQuiet {
		if *optForce {
			return NewErrUsage("cannot use both --quiet and --force")
		}
		if *optVerbose {
			return NewErrUsage("cannot use both --quiet and --verbose")
		}
	}

	if *optTop != 0 {
		if *optBottom != 0 {
			return NewErrUsage("cannot print only the top, and only the bottom.")
		}
		if *optRange != "" {
			return NewErrUsage("cannot print only the top, and only a range.")
		}
		if *optSkipBottom != 0 {
			return NewErrUsage("cannot print only the top, and skip the bottom.")
		}
		if *optSkipTop != 0 {
			return NewErrUsage("cannot print only the top, and skip the top.")
		}
		return (filter(golf.Args(), func(r io.Reader, w io.Writer) error {
			return top(int(*optTop), r, w)
		}))
	}

	if *optBottom != 0 {
		if *optRange != "" {
			return NewErrUsage("cannot print only the bottom, and only a range.")
		}
		if *optSkipBottom != 0 {
			return NewErrUsage("cannot print only the bottom, and skip the bottom.")
		}
		if *optSkipTop != 0 {
			return NewErrUsage("cannot print only the bottom, and skip the top.")
		}
		return (filter(golf.Args(), func(r io.Reader, w io.Writer) error {
			return bottom(int(*optBottom), r, w)
		}))
	}

	if *optRange != "" {
		if *optSkipBottom != 0 {
			return NewErrUsage("cannot print only a range, and skip the bottom.")
		}
		if *optSkipTop != 0 {
			return NewErrUsage("cannot print only a range, and skip the top.")
		}

		var initialLine, finalLine int
		var err error

		switch lines := strings.Split(*optRange, "-"); len(lines) {
		case 1:
			if a := lines[0]; a != "" {
				initialLine, err = strconv.Atoi(a)
				if err != nil {
					return NewErrUsage("cannot parse initial value from range: %q.", a)
				}
				finalLine = initialLine // when given a single number for a range, only print that line number
			}
		case 2:
			if a := lines[0]; a != "" {
				initialLine, err = strconv.Atoi(a)
				if err != nil {
					return NewErrUsage("cannot parse initial value from range: %q.", a)
				}
			}

			if a := lines[1]; a != "" {
				finalLine, err = strconv.Atoi(a)
				if err != nil {
					return NewErrUsage("cannot parse final value from range: %q.", a)
				}
			}

			if finalLine > 0 && initialLine > finalLine {
				return NewErrUsage("cannot print lines %d thru %d because they are out of order.", initialLine, finalLine)
			}

		default:
			return NewErrUsage("cannot print invalid range of lines: %q.", *optRange)
		}

		return (filter(golf.Args(), func(r io.Reader, w io.Writer) error {
			return copyRange(r, w, initialLine, finalLine)
		}))
	}

	return (filter(golf.Args(), func(r io.Reader, w io.Writer) error {
		return skipRange(r, w, *optSkipTop, *optSkipBottom)
	}))
}

func filter(args []string, callback func(io.Reader, io.Writer) error) error {
	if len(args) == 0 {
		return callback(os.Stdin, os.Stdout)
	}

	for _, arg := range args {
		err := withOpenFile(arg, func(fh *os.File) error {
			return callback(fh, os.Stdout)
		})
		if err != nil {
			err = fmt.Errorf("cannot read %q: %s", arg, err)
			if !*optForce {
				return err
			}
			warning("%s\n", err)
		}
	}

	return nil
}

func withOpenFile(path string, callback func(*os.File) error) (err error) {
	var fh *os.File

	fh, err = os.Open(path)
	if err != nil {
		return err
	}

	defer func() {
		if err2 := fh.Close(); err == nil {
			err = err2
		}
	}()

	// Set err variable so deferred function can inspect it.
	err = callback(fh)
	return
}

// copyRange will copy lines from r to w, starting with the line number
// corresponding to start and ending with the line number corresponding to end.
func copyRange(r io.Reader, w io.Writer, start, end int) error {
	var lineNumber int

	br := gobls.NewScanner(r)

	for br.Scan() {
		lineNumber++

		if start > 0 && lineNumber < start {
			continue
		}

		if _, err := fmt.Fprintln(w, br.Text()); err != nil {
			return err
		}

		if end > 0 && lineNumber == end {
			return nil
		}
	}

	return br.Err()
}

// skipRange will copy lines from r to w, skipping the specified number of
// initial and final lines.
func skipRange(r io.Reader, w io.Writer, initial, final uint) error {
	// Use a cirular buffer, so we are processing the Nth previous line.
	cb, err := gotb.NewStrings(int(final))
	if err != nil {
		return err
	}

	var lineNumber uint // used to skip T lines from top

	br := gobls.NewScanner(r)

	for br.Scan() {
		if initial > 0 {
			// Only need to count lines while ignoring tops.
			if lineNumber++; lineNumber <= initial {
				continue
			}
			// No reason to count lines any longer.
			initial = 0
		}

		// Recall that the circular buffer always gives us the Nth previous
		// line. When fewer than N lines have been queued, the second return
		// value will be false.
		line, ok := cb.QueueDequeue(br.Text())
		if !ok {
			continue
		}

		if _, err = fmt.Fprintln(w, line); err != nil {
			return err
		}
	}

	return br.Err()
}

// top copies the initial num lines from io.Reader to io.Writer.
func top(num int, r io.Reader, w io.Writer) error {
	if num == 0 {
		return errors.New("cannot print the initial 0 lines.")
	}

	br := gobls.NewScanner(r)

	for br.Scan() {
		if _, err := fmt.Fprintln(w, br.Text()); err != nil {
			return err
		}
		if num--; num == 0 {
			return nil
		}
	}

	return br.Err()
}

// bottom copies the final num lines from io.Reader to io.Writer.
func bottom(num int, r io.Reader, w io.Writer) error {
	if num == 0 {
		return errors.New("cannot print the final 0 lines.")
	}

	cb, err := gotb.NewStrings(num)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		_, _ = cb.QueueDequeue(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	for _, line := range cb.Drain() {
		if _, err = fmt.Fprintln(w, line); err != nil {
			return err
		}
	}

	return nil
}
