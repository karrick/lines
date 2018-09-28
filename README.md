# lines

Print a range of lines.

## Description

`lines` is a line filter that either prints a range of lines, or skips
some specified number of lines from the header, some specified number
of lines from the footer, or both.

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
