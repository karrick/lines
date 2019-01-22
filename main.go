package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/karrick/gobls"
	"github.com/karrick/golf"
	"github.com/karrick/gotb"
)

var (
	optHelp = golf.BoolP('h', "help", false, "Print command line help then exit.")

	optRange = golf.StringP('r', "range", "", "Only print lines START-END.")

	optSkipTop    = golf.Int("skip-top", 0, "Skip printing the top N header lines.")
	optSkipBottom = golf.Int("skip-bottom", 0, "Skip printing the bottom N footer lines.")

	optTop    = golf.IntP('t', "top", 0, "Only print the top N lines.")
	optBottom = golf.IntP('b', "bottom", 0, "Only print the bottom N lines.")
)

func helpThenExit(w io.Writer, err error) {
	if err != nil {
		_, _ = fmt.Fprintf(w, "ERROR: %s\n", err)
	}

	name := filepath.Base(os.Args[0])
	_, _ = fmt.Fprintf(w, "%s: Print a range of lines.\n", name)
	_, _ = fmt.Fprintf(w, "\nUSAGE: %s [ --top N | --bottom N ] [file1 [file2 ...]", name)
	_, _ = fmt.Fprintf(w, "\n  or   %s [ --range M-N | --range M- | --range -N | --range N ] [file1 [file2 ...]", name)
	_, _ = fmt.Fprintf(w, "\n  or   %s [ --skip-top N ] [ --skip-bottom N ] [file1 [file2 ...]", name)

	_, _ = fmt.Fprintf(w, "\n\n\tWithout command line arguments, reads from standard input and writes to\nstandard output. With command line arguments, reads from each file in sequence,\nand applies the below logic independently for each file.")

	_, _ = fmt.Fprintf(w, "\n\n\tWhen given the '--range N' command line argument, prints the line number\ncorresponding to N. When given the '--range START-END' command line argument,\nprints lines 'START' thru 'END', inclusively. START must not be greater than\nthe value of END. When START is omitted, the first line printed will be the\nfirst line of the input. When END is omitted, the final line printed will be\nthe final line of the input.")

	_, _ = fmt.Fprintf(w, "\n\n\tWhen given the '--skip-top N' command line argument, omits printing the\ninitial N lines, handy for removing a possibly multiline header from some text.")

	_, _ = fmt.Fprintf(w, "\n\n\tWhen given the '--skip-bottom N' command line argument, omits printing the\nfinal N lines, handy for removing a possibly multiline footer from some text.")

	_, _ = fmt.Fprintf(w, "\n\n\tWhen given the '--top N' command line argument, prints only the initial N\nlines, similar to the behavior of 'head -n N', but included in this tool for\ncompleteness.")

	_, _ = fmt.Fprintf(w, "\n\n\tWhen given the '--bottom N' command line argument, prints only the final N\nlines, similar to the behavior of 'tail -n N', but included in this tool for\ncompleteness.\n\n")

	golf.Usage()

	if err != nil {
		os.Exit(2)
	}
	os.Exit(0)
}

func main() {
	golf.Parse()

	if *optHelp {
		helpThenExit(os.Stderr, nil)
	}

	if *optTop != 0 {
		if *optTop < 0 {
			helpThenExit(os.Stderr, errors.New("cannot print a negative number of lines."))
		}
		if *optBottom != 0 {
			helpThenExit(os.Stderr, errors.New("cannot print only the top, and only the bottom."))
		}
		if *optRange != "" {
			helpThenExit(os.Stderr, errors.New("cannot print only the top, and only a range."))
		}
		if *optSkipBottom != 0 {
			helpThenExit(os.Stderr, errors.New("cannot print only the top, and skip the bottom."))
		}
		if *optSkipTop != 0 {
			helpThenExit(os.Stderr, errors.New("cannot print only the top, and skip the top."))
		}
		exit(filter(golf.Args(), func(r io.Reader, w io.Writer) error {
			return top(int(*optTop), r, w)
		}))
	}

	if *optBottom != 0 {
		if *optBottom < 0 {
			helpThenExit(os.Stderr, errors.New("cannot print a negative number of lines."))
		}
		if *optRange != "" {
			helpThenExit(os.Stderr, errors.New("cannot print only the bottom, and only a range."))
		}
		if *optSkipBottom != 0 {
			helpThenExit(os.Stderr, errors.New("cannot print only the bottom, and skip the bottom."))
		}
		if *optSkipTop != 0 {
			helpThenExit(os.Stderr, errors.New("cannot print only the bottom, and skip the top."))
		}
		exit(filter(golf.Args(), func(r io.Reader, w io.Writer) error {
			return bottom(int(*optBottom), r, w)
		}))
	}

	if *optRange != "" {
		if *optSkipBottom != 0 {
			helpThenExit(os.Stderr, errors.New("cannot print only a range, and skip the bottom."))
		}
		if *optSkipTop != 0 {
			helpThenExit(os.Stderr, errors.New("cannot print only a range, and skip the top."))
		}

		var initialLine, finalLine int
		var err error

		switch lines := strings.Split(*optRange, "-"); len(lines) {
		case 1:
			if a := lines[0]; a != "" {
				initialLine, err = strconv.Atoi(a)
				if err != nil {
					helpThenExit(os.Stderr, fmt.Errorf("cannot parse initial value from range: %q.", a))
				}
				finalLine = initialLine // when given a single number for a range, only print that line number
			}
		case 2:
			if a := lines[0]; a != "" {
				initialLine, err = strconv.Atoi(a)
				if err != nil {
					helpThenExit(os.Stderr, fmt.Errorf("cannot parse initial value from range: %q.", a))
				}
			}

			if a := lines[1]; a != "" {
				finalLine, err = strconv.Atoi(a)
				if err != nil {
					helpThenExit(os.Stderr, fmt.Errorf("cannot parse final value from range: %q.", a))
				}
			}

			if finalLine > 0 && initialLine > finalLine {
				helpThenExit(os.Stderr, fmt.Errorf("cannot print lines %d thru %d because they are out of order.", initialLine, finalLine))
			}

		default:
			helpThenExit(os.Stderr, fmt.Errorf("cannot print invalid range of lines: %q.", *optRange))
		}

		exit(filter(golf.Args(), func(r io.Reader, w io.Writer) error {
			return copyRange(r, w, initialLine, finalLine)
		}))
	}

	if *optSkipTop < 0 || *optSkipBottom < 0 {
		helpThenExit(os.Stderr, errors.New("cannot print a negative number of lines."))
	}

	exit(filter(golf.Args(), func(r io.Reader, w io.Writer) error {
		return skipRange(r, w, *optSkipTop, *optSkipBottom)
	}))
}

func exit(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
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
			_, _ = fmt.Fprintf(os.Stderr, "WARNING: %s\n", err)
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
	if err := br.Err(); err != nil {
		return err
	}

	return nil
}

// skipRange will copy lines from r to w, skipping the specified number of
// initial and final lines.
func skipRange(r io.Reader, w io.Writer, initial, final int) error {
	// Use a cirular buffer, so we are processing the Nth previous line.
	cb, err := gotb.NewStrings(final)
	if err != nil {
		return err
	}

	var lineNumber int // used to skip T lines from top

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
	if err := br.Err(); err != nil {
		return err
	}

	return nil
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
	if err := br.Err(); err != nil {
		return err
	}

	return nil
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
