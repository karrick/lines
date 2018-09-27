# gotb

Implements a few tail buffer abstract data types.

## Description

Tail buffers are useful when a program needs to track the N final
elements added to a list, but not necessarily track previous elements.

```Go
// tail copies the final num lines from io.Reader to io.Writer.
func tail(num int, r io.Reader, w io.Writer) error {
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
```
