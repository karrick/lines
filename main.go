package main

import (
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
	optBottom = golf.IntP('b', "bottom", 0, "omit final 'BOTTOM' lines")
	optHelp   = golf.BoolP('h', "help", false, "Print command line help and exit")
	optRange  = golf.StringP('r', "range", "", "Print lines START-END")
	optTop    = golf.IntP('t', "top", 0, "omit initial 'TOP' lines")
)

func help(w *golinewrap.Writer, err error) {
	if err != nil {
		_, _ = w.WriteParagraph(fmt.Sprintf("ERROR: %s", err))
	}

	name := filepath.Base(os.Args[0])
	_, _ = w.WriteParagraph(fmt.Sprintf("%s: Either print lines 'START' thru 'END', or skip 'TOP' top lines and 'BOTTOM' bottom lines.", name))
	_, _ = w.WriteParagraph(fmt.Sprintf("USAGE:\t%s [ [-r [NUMBER | START-END]] | [-b BOTTOM] [-t TOP] ] [file1 [file2 ...] ]", name))

	_, _ = w.WriteParagraph(`When given the '--range NUMBER' command line
	argument, prints the line number corresponding to NUMBER. When given the
	'--range START-END' command line argument, prints lines 'START' thru 'END',
	inclusively. START must not be greater than the value of END. When START is
	omitted, the first line printed will be the first line of the input. When
	END is omitted, the final line printed will be the final line of the
	input.`)

	_, _ = w.WriteParagraph(`When given the '--top TOP' command line argument,
	omits the initial 'TOP' lines. When given the '--bottom BOTTOM' command line
	argument, omits the final 'BOTTOM' lines.`)

	_, _ = w.WriteParagraph(`Without command line arguments, reads from standard
	input and writes to standard output. With command line arguments, reads from
	each file in sequence, and applies the above logic independently for each
	file.`)

	golf.Usage()
}

func main() {
	golf.Parse()

	lw := lineWrapping(os.Stderr, "")

	if *optHelp {
		help(lw, nil)
		os.Exit(0)
	}

	if *optRange != "" {
		var initialLine, finalLine int
		var err error

		switch lines := strings.Split(*optRange, "-"); len(lines) {
		case 1:
			if a := lines[0]; a != "" {
				initialLine, err = strconv.Atoi(a)
				if err != nil {
					help(lw, fmt.Errorf("cannot parse initial value from range: %q.", a))
					os.Exit(2)
				}
				finalLine = initialLine // when given a single number for a range, only print that line number
			}
		case 2:
			if a := lines[0]; a != "" {
				initialLine, err = strconv.Atoi(a)
				if err != nil {
					help(lw, fmt.Errorf("cannot parse initial value from range: %q.", a))
					os.Exit(2)
				}
			}

			if a := lines[1]; a != "" {
				finalLine, err = strconv.Atoi(a)
				if err != nil {
					help(lw, fmt.Errorf("cannot parse final value from range: %q.", a))
					os.Exit(2)
				}
			}

			if finalLine > 0 && initialLine > finalLine {
				help(lw, fmt.Errorf("cannot print lines %d thru %d because they are out of order.", initialLine, finalLine))
				os.Exit(2)
			}

		default:
			help(lw, fmt.Errorf("cannot print invalid range of lines: %q.", *optRange))
			os.Exit(2)
		}

		if golf.NArg() == 0 {
			if err := rangeReader(os.Stdin, initialLine, finalLine); err != nil {
				_, _ = lw.WriteParagraph(fmt.Sprintf("ERROR: %s", err))
			}
			os.Exit(0)
		}

		for _, arg := range golf.Args() {
			err := withOpenFile(arg, func(fh *os.File) error {
				return rangeReader(fh, initialLine, finalLine)
			})
			if err != nil {
				_, _ = lw.WriteParagraph(fmt.Sprintf("WARNING: %s", err))
			}
		}

		os.Exit(0)
	}

	if golf.NArg() == 0 {
		if err := skipReader(os.Stdin, *optTop, *optBottom); err != nil {
			_, _ = lw.WriteParagraph(fmt.Sprintf("ERROR: %s", err))
		}
		os.Exit(0)
	}

	for _, arg := range golf.Args() {
		err := withOpenFile(arg, func(fh *os.File) error {
			return skipReader(fh, *optTop, *optBottom)
		})
		if err != nil {
			_, _ = lw.WriteParagraph(fmt.Sprintf("WARNING: %s", err))
		}
	}
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

func rangeReader(ior io.Reader, top, bottom int) error {
	var lineNumber int

	br := gobls.NewScanner(ior)

	for br.Scan() {
		lineNumber++

		if top > 0 && lineNumber < top {
			continue
		}

		if bottom > 0 && lineNumber > bottom {
			return nil
		}

		if _, err := fmt.Println(br.Text()); err != nil {
			return err
		}
	}
	if err := br.Err(); err != nil {
		return err
	}

	return nil
}

func skipReader(ior io.Reader, top, bottom int) error {
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

		if _, err = fmt.Println(line); err != nil {
			return err
		}
	}
	if err := br.Err(); err != nil {
		return err
	}

	return nil
}
