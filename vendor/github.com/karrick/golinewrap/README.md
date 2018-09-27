# golinewrap

A Go library for line wrapping when writing to an io.Writer.

[![GoDoc](https://godoc.org/github.com/karrick/golinewrap?status.svg)](https://godoc.org/github.com/karrick/golinewrap)

## Description

Writes the to the underlying io.Writer, wrapping lines as necessary to
prevent line lengths from exceeding the pre-configured width. Each
pilcrow rune, ¶, in the byte sequence causes a new paragraph to be
emitted.

## Example

```Go
func lineWrapping(w io.Writer) io.Writer {
	columns, _, err := gows.GetWinSize()
	if err != nil {
		return w
	}

	lw, err := golinewrap.New(w, columns, "# ")
	if err != nil {
		return w
	}

	return lw
}

func exit(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(lineWrapping(os.Stderr), "ERROR: %s", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func usage(err error) {
	_, _ = fmt.Fprintf(lineWrapping(os.Stderr), "USAGE: %s", err)
	os.Exit(2)
}
```
