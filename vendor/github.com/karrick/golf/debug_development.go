// +build golf_debug

package golf

import (
	"fmt"
	"os"
)

// debug formats and prints arguments to stderr for development builds
func debug(f string, a ...interface{}) {
	os.Stderr.Write([]byte("golf: " + fmt.Sprintf(f, a...)))
}
