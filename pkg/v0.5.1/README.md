# lines

Print a range of lines from standard input or one or more files.

## Description

`lines` is a command line filter that either prints a range of lines,
or skips some specified number of lines from the header, some
specified number of lines from the footer, or both.

While use of the `head` and `tail` POSIX command line programs is
already easy to use, how does one go about skipping the initial M
lines of a file, or skipping the final N lines of a file? How does one
go about skipping both M lines from the top _and_ N lines from the
bottom?

Every time I need to do this I spend time doing research on the
Internet, judging among a handfull of solutions using `awk`, `perl`,
`python`, `ruby`, or some other scripted solution. Each of the
proposals has a slightly different syntax, and some of them don't even
work. However all I really want is a filter that does what I need
without having to search and study those respective man pages. Is that
too much to ask?

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

`lines` accepts either `--range STRING` or `-r STRING` to specify a
range of lines to print. In these examples, I always use the short
flag.

In general, `lines -r M-N` is equivalent to `sed -n M,Np`, although
lines allows omitting M, in which case it will default to the first
line, and N, in which case it will default to the last line.

```Bash
$ lines sample.txt -r 4-7
4: test
5: test
6: test
7: test
```

Interestingly, when using the short flag name, one may omit the
space between the short flag letter and the argument. Therefore `-r
4-7` is the same as `-r4-7`.

Either or both of the ends of the range parameter may be omitted. When
the first number is omitted, printing starts at the first line. When
the final number is omitted, printing ends at the final line.

As previously described, the intervening space between the flag letter
and the argument may be omitted, causing `-r -3` to have the same
meaning as `-r-3`, printing the first 3 lines of the file. Printing
the first 3 lines of the file is equivalent to `lines --top 3`.

```Bash
$ lines sample.txt -r -3
1: test
2: test
3: test
```

Both `-r 7-` and `-r7-` both print lines 7 thru the end of the
file. Note this is different than printing the final 7 lines of the
file, as one might do with `lines --bottom 7`.

```Bash
$ lines sample.txt -r 7-
7: test
8: test
9: test
10: test
```

### Printing a single line using '--range N'

Equivalent to `sed -n Np`.

```Bash
$ lines sample.txt -r 3
3: test
```

### Omitting one or more header lines using '--skip-top N'

Equivalent to `(( M+=1 )) ; sed -n "$M,\$p"`, although that modifies M
along the way.

```Bash
$ lines sample.txt --skip-top 2
3: test
4: test
5: test
6: test
7: test
8: test
9: test
10: test
```

### Omitting one or more footer lines using '--skip-bottom N'

Equivalent to dying your hair gray, because I have not found a way to
do this with `sed`. Maybe I should give `awk` a swing...

```Bash
$ lines sample.txt --skip-bottom 2
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

See above, regarding omitting lines from the footer.

```Bash
$ lines sample.txt --skip-top 3 --skip-bottom 2
4: test
5: test
6: test
7: test
8: test
```

### Printing only the initial N lines

Equivalent to `head -n N`. Also note this will have the same effect as
calling `-r -N`.

```Bash
$ lines sample.txt --top 3
1: test
2: test
3: test
```

### Printing only the final N lines

Equivalent to `tail -n 3`.

```Bash
$ lines sample.txt --bottom 3
8: test
9: test
10: test
```

## Installation

If you don't have the Go programming language installed, then you'll
need to install a copy from [https://golang.org/dl](https://golang.org/dl).

Once you have Go installed:

    $ go get github.com/karrick/lines
