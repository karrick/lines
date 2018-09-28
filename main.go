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
	"github.com/karrick/golinewrap"
	"github.com/karrick/gotb"
	"github.com/karrick/gows"
)

var (
	optHelp = golf.BoolP('h', "help", false, "Print command line help then exit")

	optHead = golf.Uint("head", 0, "Only print the initial N lines")
	optTail = golf.Uint("tail", 0, "Only print the final N lines")

	optRange = golf.StringP('r', "range", "", "Only print lines START-END")

	optHeader = golf.Uint("header", 0, "Skip printing the initial N header lines")
	optFooter = golf.Uint("footer", 0, "Skip printing the final N footer lines")
)

func helpThenExit(w *golinewrap.Writer, err error) {
	if err != nil {
		_, _ = w.WriteParagraph(fmt.Sprintf("ERROR: %s", err))
	}

	name := filepath.Base(os.Args[0])
	_, _ = w.WriteParagraph(fmt.Sprintf("%s: Print a range of lines.", name))
	_, _ = w.WriteParagraph(fmt.Sprintf("USAGE:\t%s [ --head N | --tail N ] [file1 [file2 ...]", name))
	_, _ = w.WriteParagraph(fmt.Sprintf("USAGE:\t%s [ --range M-N | --range M- | --range -N | --range N ] [file1 [file2 ...]", name))
	_, _ = w.WriteParagraph(fmt.Sprintf("USAGE:\t%s [ --header N ] [ --footer N ] [file1 [file2 ...]", name))

	_, _ = w.WriteParagraph(`Without command line arguments, reads from standard
	input and writes to standard output. With command line arguments, reads from
	each file in sequence, and applies the below logic independently for each
	file.`)

	_, _ = w.WriteParagraph(`When given the '--head N' command line argument,
	prints only the initial N lines, similar to the behavior of 'head -n N', but
	included in this tool for completeness.`)

	_, _ = w.WriteParagraph(`When given the '--tail N' command line argument,
	prints only the final N lines, similar to the behavior of 'tail -n N', but
	included in this tool for completeness.`)

	_, _ = w.WriteParagraph(`When given the '--range N' command line argument,
	prints the line number corresponding to N. When given the '--range
	START-END' command line argument, prints lines 'START' thru 'END',
	inclusively. START must not be greater than the value of END. When START is
	omitted, the first line printed will be the first line of the input. When
	END is omitted, the final line printed will be the final line of the
	input.`)

	_, _ = w.WriteParagraph(`When given the '--header N' command line argument,
	omits printing the initial N lines, handy for removing a possibly multiline
	header from some text.`)

	_, _ = w.WriteParagraph(`When given the '--footer N' command line argument,
	omits printing the final N lines, handy for removing a possibly multiline
	footer from some text.`)

	golf.Usage()

	if err != nil {
		os.Exit(2)
	}
	os.Exit(0)
}

func main() {
	golf.Parse()

	lw := lineWrapping(os.Stderr, "")

	if *optHelp {
		helpThenExit(lw, nil)
	}

	if *optHead != 0 {
		if *optTail != 0 {
			helpThenExit(lw, errors.New("cannot print only the head, and only the tail."))
		}
		if *optRange != "" {
			helpThenExit(lw, errors.New("cannot print only the head, and only a range."))
		}
		if *optFooter != 0 {
			helpThenExit(lw, errors.New("cannot print only the head, and skip the footer."))
		}
		if *optHeader != 0 {
			helpThenExit(lw, errors.New("cannot print only the head, and skip the header."))
		}
		exit(filter(golf.Args(), func(r io.Reader, w io.Writer) error {
			return head(int(*optHead), r, w)
		}))
	}

	if *optTail != 0 {
		if *optRange != "" {
			helpThenExit(lw, errors.New("cannot print only the tail, and only a range."))
		}
		if *optFooter != 0 {
			helpThenExit(lw, errors.New("cannot print only the tail, and skip the footer."))
		}
		if *optHeader != 0 {
			helpThenExit(lw, errors.New("cannot print only the tail, and skip the header."))
		}
		exit(filter(golf.Args(), func(r io.Reader, w io.Writer) error {
			return tail(int(*optTail), r, w)
		}))
	}

	if *optRange != "" {
		if *optFooter != 0 {
			helpThenExit(lw, errors.New("cannot print only a range, and skip the footer."))
		}
		if *optHeader != 0 {
			helpThenExit(lw, errors.New("cannot print only a range, and skip the header."))
		}

		var initialLine, finalLine int
		var err error

		switch lines := strings.Split(*optRange, "-"); len(lines) {
		case 1:
			if a := lines[0]; a != "" {
				initialLine, err = strconv.Atoi(a)
				if err != nil {
					helpThenExit(lw, fmt.Errorf("cannot parse initial value from range: %q.", a))
				}
				finalLine = initialLine // when given a single number for a range, only print that line number
			}
		case 2:
			if a := lines[0]; a != "" {
				initialLine, err = strconv.Atoi(a)
				if err != nil {
					helpThenExit(lw, fmt.Errorf("cannot parse initial value from range: %q.", a))
				}
			}

			if a := lines[1]; a != "" {
				finalLine, err = strconv.Atoi(a)
				if err != nil {
					helpThenExit(lw, fmt.Errorf("cannot parse final value from range: %q.", a))
				}
			}

			if finalLine > 0 && initialLine > finalLine {
				helpThenExit(lw, fmt.Errorf("cannot print lines %d thru %d because they are out of order.", initialLine, finalLine))
			}

		default:
			helpThenExit(lw, fmt.Errorf("cannot print invalid range of lines: %q.", *optRange))
		}

		exit(filter(golf.Args(), func(r io.Reader, w io.Writer) error {
			return rangeReader(r, w, initialLine, finalLine)
		}))
	}

	exit(filter(golf.Args(), func(r io.Reader, w io.Writer) error {
		return skipReader(r, w, int(*optHeader), int(*optFooter))
	}))
}

func exit(err error) {
	if err != nil {
		_, _ = lineWrapping(os.Stderr, "").WriteParagraph(fmt.Sprintf("ERROR: %s", err))
		os.Exit(1)
	}
	os.Exit(0)
}

func lineWrapping(w io.Writer, prefix string) *golinewrap.Writer {
	columns, _, err := gows.GetWinSize()

	if columns == 0 || columns >= 80 {
		columns = 79
	} else {
		columns--
	}

	lw, err := golinewrap.New(w, columns, prefix)
	if err != nil {
		panic(err)
	}

	return lw
}

func filter(args []string, callback func(io.Reader, io.Writer) error) error {
	if len(args) == 0 {
		return callback(os.Stdin, os.Stdout)
	}

	lw := lineWrapping(os.Stderr, "")

	for _, arg := range args {
		err := withOpenFile(arg, func(fh *os.File) error {
			return callback(fh, os.Stdout)
		})
		if err != nil {
			_, _ = lw.WriteParagraph(fmt.Sprintf("WARNING: %s", err))
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

func rangeReader(ior io.Reader, w io.Writer, top, bottom int) error {
	var lineNumber int

	br := gobls.NewScanner(ior)

	for br.Scan() {
		lineNumber++

		if top > 0 && lineNumber < top {
			continue
		}

		if _, err := fmt.Fprintln(w, br.Text()); err != nil {
			return err
		}

		if bottom > 0 && lineNumber == bottom {
			return nil
		}
	}
	if err := br.Err(); err != nil {
		return err
	}

	return nil
}

func skipReader(ior io.Reader, w io.Writer, top, bottom int) error {
	// Use a cirular buffer, so we are processing the Nth previous line.
	cb, err := gotb.NewStrings(bottom)
	if err != nil {
		return err
	}

	var lineNumber int // used to skip T lines from top

	br := gobls.NewScanner(ior)

	for br.Scan() {
		if top > 0 {
			// Only need to count lines while ignoring tops.
			if lineNumber++; lineNumber <= top {
				continue
			}
			// No reason to count lines any longer.
			top = 0
		}

		// Recall circular buffer always gives us the Nth previous line, or a
		// false for the second return value.
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

// head copies the initial num lines from io.Reader to io.Writer.
func head(num int, r io.Reader, w io.Writer) error {
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

// tail copies the final num lines from io.Reader to io.Writer.
func tail(num int, r io.Reader, w io.Writer) error {
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
