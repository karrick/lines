package golf

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	flags              []option
	remainingArguments []string // keep track of remaining arguments
	argsProcessed      int      // keep track of how many arguments have been set
	parsed             bool     // keep track of whether command line arguments have been parsed
)

func resetParser() {
	argsProcessed = 0
	flags = nil
	parsed = false
	remainingArguments = nil
}

// Parse parses the command line. On error, displays the usage of the command
// line and exits the program with status code 2.
func Parse() {
	// combine os.Args to a single string
	if err := parseArgs(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		Usage()
		os.Exit(2)
	}
	parsed = true
}

type slurpType uint

const (
	nothingToSlurp slurpType = iota
	slurpDuration
	slurpFloat
	slurpInt
	slurpInt64
	slurpString
	slurpUint
	slurpUint64
)

func (slurp slurpType) String() string {
	switch slurp {
	case nothingToSlurp:
		return "nothing to slurp"
	case slurpDuration:
		return "slurp duration"
	case slurpFloat:
		return "slurp float"
	case slurpInt:
		return "slurp int"
	case slurpInt64:
		return "slurp int64"
	case slurpString:
		return "slurp string"
	case slurpUint:
		return "slurp uint"
	case slurpUint64:
		return "slurp uint64"
	default:
		return fmt.Sprintf("unknown slurp type value: %d", int(slurp))
	}
}

// runeParserStateType represents the possible states of the parser, implemented as a state machine.
type runeParserStateType uint

const (
	beginArgument runeParserStateType = iota
	consumedHyphen
	wantLongName
	wantShortFlagsOnly
	wantText
	wantArgument
)

// String returns text representation of parser state.
func (state runeParserStateType) String() string {
	switch state {
	case beginArgument:
		return "new argument"
	case consumedHyphen:
		return "consumed single hyphen"
	case wantLongName:
		return "want long name"
	case wantShortFlagsOnly:
		return "want short flags only"
	case wantText:
		return "want text"
	case wantArgument:
		return "want argument"
	default:
		return fmt.Sprintf("unknown rune parser state value: %d", int(state))
	}
}

// Parsed reports whether the command-line flags have been parsed.
func Parsed() bool {
	return parsed
}

func slurpText(text string, nextSlurp slurpType, f option) error {
	var ui64 uint64
	var err error

	switch nextSlurp {
	case slurpDuration:
		*f.(*optionDuration).pv, err = time.ParseDuration(text)
	case slurpFloat:
		*f.(*optionFloat).pv, err = strconv.ParseFloat(text, 64)
	case slurpInt:
		*f.(*optionInt).pv, err = strconv.Atoi(text)
	case slurpInt64:
		*f.(*optionInt64).pv, err = strconv.ParseInt(text, 10, 64)
	case slurpUint:
		ui64, err = strconv.ParseUint(text, 10, 0)
		*f.(*optionUint).pv = uint(ui64)
	case slurpUint64:
		*f.(*optionUint64).pv, err = strconv.ParseUint(text, 10, 64)
	case slurpString:
		*f.(*optionString).pv = text
	default:
		err = fmt.Errorf("unexpected slurp state: %v", nextSlurp)
	}

	return err
}

func parseArgs(args []string) error {
	var flagType slurpType
	var flagName, flagText string
	var f option

	for ai, arg := range args { // ai (arg index)
		// fmt.Fprintf(os.Stderr, "arg %d: %q; start argParserState: %v\n", ai, arg, flagType)

		if flagType != nothingToSlurp {
			if err := slurpText(arg, flagType, f); err != nil {
				return err
			}
			flagType = nothingToSlurp
			argsProcessed++
			continue // with next argument
		}

		// nothing to slurp, so need to read runes one by one
		runeParserState := beginArgument
		flagText = ""

		for bi, r := range arg { // bi (byte index)
			// fmt.Fprintf(os.Stderr, "  runeParserState: %v; rune: %q\n", runeParserState, r)

			// NOTE: Cannot put this in switch statement because break we might
			// need to break out of enclosing for statement that loops over arg
			// runes.
			if runeParserState == wantText {
				flagText = arg[bi:]
				// fmt.Fprintf(os.Stderr, "  FLAG TEXT: %q\n", flagText)
				break // out of parsing this arg
			} else if runeParserState == wantLongName {
				flagName = arg[bi:]
				// fmt.Fprintf(os.Stderr, "  LONG FLAG NAME: %q\n", flagName)
				break // out of parsing this arg
			} else if runeParserState == beginArgument {
				if r != '-' {
					// fmt.Fprintf(os.Stderr, "index: %d; this rune ends processing: %q\n", ai, r)
					// remainingArguments = args[ai:]
					// argsProcessed = ai
					// return nil
					runeParserState = wantArgument
					break // out of parsing this arg
				}
				runeParserState = consumedHyphen
			} else if runeParserState == consumedHyphen {
				switch r {
				case '-':
					runeParserState = wantLongName
				default:
					if f = flagFromShortName(r); f == nil {
						argsProcessed = ai
						return fmt.Errorf("unknown flag: %q", r)
					}
					switch flagType = f.NextSlurp(); flagType {
					case nothingToSlurp:
						*f.(*optionBool).pv = true
						runeParserState = wantShortFlagsOnly
					default:
						runeParserState = wantText
					}
				}
			} else if runeParserState == wantShortFlagsOnly {
				if f = flagFromShortName(r); f == nil {
					argsProcessed = ai
					return fmt.Errorf("unknown flag: %q", r)
				}
				switch flagType = f.NextSlurp(); flagType {
				case nothingToSlurp:
					*f.(*optionBool).pv = true
				default:
					runeParserState = wantText
				}
			} else {
				panic(fmt.Errorf("TODO: runeParserState: %v; rune: %q", runeParserState, r))
			}
		}

		switch runeParserState {
		case consumedHyphen:
			argsProcessed = ai
			return errors.New("hyphen without flags")
		case wantArgument:
			remainingArguments = append(remainingArguments, arg)
		case wantText:
			if flagText == "" {
				// fmt.Fprintf(os.Stderr, "  finished arg and got no text: %q\n", flagText)
				argsProcessed++
				continue // with next arg, where we will slurp in the value
			}
			// fmt.Fprintf(os.Stderr, "  finished arg and got text: %q\n", flagText)
			if flagType == nothingToSlurp {
				panic(fmt.Errorf("got text %q but invalid nextSlurp: %v", flagText, flagType))
			}
			if err := slurpText(flagText, flagType, f); err != nil {
				return err
			}
			flagType = nothingToSlurp
			argsProcessed++
		case wantShortFlagsOnly:
			// fmt.Fprintf(os.Stderr, "  finished arg while looking for short flags\n")
			argsProcessed++
		case wantLongName:
			if flagName == "" {
				return nil
			}
			if f = flagFromLongName(flagName); f == nil {
				argsProcessed = ai
				return fmt.Errorf("unknown flag: %q", flagName)
			}
			if flagType = f.NextSlurp(); flagType == nothingToSlurp {
				*f.(*optionBool).pv = true
			}
			argsProcessed++
		default:
			panic(fmt.Errorf("TODO: handle runeParserState: %v", runeParserState))
		}

		// POST: end of an indexed argument
		// fmt.Fprintf(os.Stderr, "arg %d: %q end argParserState: %v\n", ai, arg, flagType)
	}

	// fmt.Fprintf(os.Stderr, "after all args: %v\n", flagType)

	// we might have hit the end of the command line while doing stuff
	if flagType != nothingToSlurp {
		if long := f.Long(); long != "" {
			return fmt.Errorf("flag requires argument: %q", long)
		}
		return fmt.Errorf("flag requires argument: %q", f.Short())
	}

	return nil
}
