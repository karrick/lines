# lines

Print a range of lines.

## Description

`lines` is a line filter that either prints a range of lines, or skips
some specified number of lines from the header, some specified number
of lines from the footer, or both.

While use of the `head` and `tail` POSIX command line programs is
already easy to use, how does one go about skipping the initial N
lines of a file, or skipping the final N lines of a file? Every
solution I have seen involves an Internet search, followed by a
handfull of solutions using `awk`, `perl`, `python`, `ruby`, or some
other scripted solution, each with its own language.

I just want a UNIX filter command to print a range of lines, or skip
lines from being printed. Is that too much to ask?

I do not claim any of this is an original idea. But I have not found a
similar program elsewhere, so I decided to write it myself. I hope it
serves others well.

## Example

All of the examples assume the following input file, `sample.txt`.

```
1: test
2: test
3: test
4: test
5: test
6: test
7: test
8: test
9: test
10: test
```

### Printing a range of lines using '--range START-END'

```Bash
$ lines sample.txt -r 4-7
4: test
5: test
6: test
7: test
```

Either or both of the ends of the range parameter may be omitted. When
the first number is omitted, printing starts at the first line. When
the final number is omitted, printing ends at the final line.

```Bash
$ lines sample.txt -r -3
1: test
2: test
3: test
```

```Bash
$ lines sample.txt -r 7-
7: test
8: test
9: test
10: test
```

### Printing a single line using '--range N'

```Bash
$ lines sample.txt -r 3
3: test
```

### Omitting one or more header lines using '--header N'

```Bash
$ lines sample.txt --header 2
4: test
5: test
6: test
7: test
8: test
9: test
10: test
```

### Omitting one or more footer lines using '--footer N'

```Bash
$ lines sample.txt --footer 2
1: test
2: test
3: test
4: test
5: test
6: test
7: test
8: test
```

### Omitting lines from both the header and footer

```Bash
$ lines sample.txt --header 3 --footer 2
4: test
5: test
6: test
7: test
8: test
```

### Printing only the initial N lines

Duplicates behavior of invoking `head -n 3`, but included here for
completeness.

```Bash
$ lines sample.txt --head 3
1: test
2: test
3: test
```

### Printing only the final N lines

Duplicates behavior of invoking `tail -n 3`, but included here for
completeness.

```Bash
$ lines sample.txt --tail 3
8: test
9: test
10: test
```
