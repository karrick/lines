# golinewrap

A Go library for line wrapping when writing to an io.Writer.

[![GoDoc](https://godoc.org/github.com/karrick/golinewrap?status.svg)](https://godoc.org/github.com/karrick/golinewrap)

## Description

Writes the to the underlying io.Writer, wrapping lines as necessary to prevent
line lengths from exceeding the pre-configured width.

This library does not split words to keep them within the specified line length,
but rather emits them on a line of their own. A future version of this library
may address this shortcoming.

## Examples

```Go
package main

import (
	"os"

	"github.com/karrick/golinewrap"
)

func main() {
	// Creating a new golinewarp.Writer wraps an existing io.Writer, but also
	// accepts the number of columns to wrap text to, along with a prefix
	// string that will be emitted at the start of every line. The prefix
	// string may be empty, but it may not be as large as or larger than the
	// number of columns specified.
	//
	// NOTE: When line wrapping based on terminal width, remember to save one
	// column for the newline character, or your output will be stuttered.
	lw, err := golinewrap.New(os.Stdout, 79, "> ")
	if err != nil {
		// NOTE: New function will only return an error when prefix wider than
		// the number of columns specified, which should not be possible when
		// calling with constants 79 and a prefix only several runes long.
		panic(err)
	}

	lw.Write([]byte(`Because golinewrap.Writer implements the io.Writer interface, it may be used anywhere an io.Writer is expected, and will wrap lines accordingly.
    When its Write method is called, as is done here, it splits its buffer input by newline characters, then emits each line as a paragraph. This may or may not be what you want, but it ensures maximum compatibility when you must provide an io.Writer to another function and cannot pass a golinewrap.Writer.`))

	lw.WriteParagraph(`While golinewrap.Writer _is_ an io.Writer, it also has
    some other methods which make it more natural to use when formatting
    paragraphs of information.`)

	lw.Printf("One such handy method is %q, formatting its arguments using %q, then writing the resulting string to the underlying io.Writer.", "Printf", "fmt.Sprintf")

	lw.WriteParagraph(`The one I tend to use most often is
    golinewrap.WriteParagraph, which accepts a paragraph of text, then emits
    that paragraph to the underlying io.Writer. Calling WriteParagraph multiple
    times results in multiple paragraphs of text to be represented in the
    output.`)

	lw.WriteWord("It")
	lw.WriteWord("also")
	lw.WriteWord("supports")
	lw.WriteWord("writing")
	lw.WriteWord("one")
	lw.WriteWord("word")
	lw.WriteWord("at")
	lw.WriteWord("a")
	lw.WriteWord("time.")

	lw.WriteRune('\n') // It also provides emitting a single rune. This is
	lw.WriteRune('\n') // handy when you need to add some extra vertical rows
	lw.WriteRune('\n') // to your output, similar to calling `fmt.Println()`.

	lw.WriteParagraph(`Invoking WriteParagraph, WriteWord, and WriteRune will
    always flush output up to and including the final newline, so this library
    can be more easily used in a line streaming context. However, writing a
    single word or a single rune will trigger a newline only once a line is
    completed in terms of its length.`)
}
```
