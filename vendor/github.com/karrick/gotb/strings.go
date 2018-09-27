package gotb

import "fmt"

// Strings is a data structure for storing the previous N string items, where 0
// <= N <= arbitrary limit, and the arbitrary limit is proportional to how much
// memory ought to be allocated to the data structure.
type Strings struct {
	items  []string // items slice will be used as a circular buffer
	index  int      // index points to current item in buffer
	looped bool     // true once circular buffer has all slots filled
}

// NewStrings returns a newly initialized buffer for N string items, where 0 <=
// N.
func NewStrings(n int) (*Strings, error) {
	switch {
	case n < 0:
		return nil, fmt.Errorf("cannot create buffer with negative item count: %d", n)
	case n == 0:
		// Rather than returning an error, below methods provide shortcut when n
		// is 0.
		return new(Strings), nil
	default:
		return &Strings{items: make([]string, n)}, nil
	}
}

// QueueDequeue stores the newly provided item in the queue and returns the Nth
// previous item from the queue, along with a second return value of true. If
// exactly N or fewer than N items have thus far been stored in the buffer, an
// empty string will be returned along with a second return value of false.
func (tb *Strings) QueueDequeue(newItem string) (string, bool) {
	// Special case when the circular buffer was not allocated: just return the
	// provided item.
	if tb.items == nil {
		return newItem, true
	}

	// Swap item previously stored at index with new item.
	prevItem := tb.items[tb.index]
	tb.items[tb.index] = newItem
	valid := tb.looped

	// Increment index making note whether it wraps.
	if tb.index++; tb.index == cap(tb.items) {
		tb.index = 0
		tb.looped = true
	}

	return prevItem, valid
}

// Drain returns all items from the structure. This implimentation is not
// designed to handle invocation of any other methods after calling Drain.
func (tb *Strings) Drain() []string {
	if tb.looped {
		return append(tb.items[tb.index:], tb.items[:tb.index]...) // f g c d e
	}
	return tb.items[:tb.index] // a b c _ _
}
