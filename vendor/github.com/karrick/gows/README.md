# gows

Golang library to get the terminal or console window size 

## Usage

```Go
package main

import (
	"fmt"
	"os"

	"github.com/karrick/gows"
)

func main() {
	col, row, err := gows.GetWinSize()
	if err != nil {
		exit(err)
	}
	fmt.Printf("%d %d\n", row, col) // output in same order as `stty size`
}

func exit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
```
