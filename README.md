# lines

`lines` is a command line filter to print or omit lines.

## Description

`lines` is a line filter that either prints a range of lines, or skips
some specified number of lines from the header, and some specified
number of lines from the footer.

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

### Printing a single line

```Bash
$ lines sample.txt -r 3
3: test
```

### Printing a range of lines

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

### Omitting one or more header lines

```Bash
$ lines sample.txt -t 2
4: test
5: test
6: test
7: test
8: test
9: test
10: test
```

### Omitting one or more footer lines

```Bash
$ lines sample.txt -b 2
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
$ lines sample.txt -t 3 -b 2
4: test
5: test
6: test
7: test
8: test
```
