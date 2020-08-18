// +build !golf_debug

package golf

// debug is a no-op for release builds
func debug(_ string, _ ...interface{}) {}
