package golf

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

// flagFromLongName performs linear search for the flag with a matching long
// flag name in the list of flags. It returns the flag found, or nil if the
// requested long flag name was not found.
func flagFromLongName(long string) option {
	for _, option := range flags {
		if option.Long() == long {
			return option
		}
	}
	return nil
}

// flagFromShortName performs linear search for the flag with a matching short
// flag name in the list of flags. It returns the flag found, or nil if the
// requested short flag name was not found.
func flagFromShortName(short rune) option {
	for _, option := range flags {
		if option.Short() == short {
			return option
		}
	}
	return nil
}

// redefinition ensures the specified short and long flags do not redefine any
// existing flag definitions.
func redefinition(short rune, long string) error {
	for _, opt := range flags {
		if long != "" && long == opt.Long() {
			return fmt.Errorf("cannot add option that duplicates long flag: %q", long)
		}
		if short != utf8.RuneError && short == opt.Short() {
			return fmt.Errorf("cannot add option that duplicates short flag: %q", short)
		}
	}
	return nil
}

// parseShortAndLongFlag is called when there is both a short and a long flag to
// validate and ensure there are no duplicates.
func parseShortAndLongFlag(short rune, long string) error {

	switch short {
	case utf8.RuneError:
		return fmt.Errorf("cannot use flag with invalid rune: %q", short)
	case '-':
		return fmt.Errorf("cannot use hyphen as a flag: %q", short)
	}

	switch {
	case long == "":
		return errors.New("cannot use empty flag string")
	case strings.HasPrefix(long, "-"):
		return fmt.Errorf("cannot use flag that starts with a hyphen: %q", long)
	}

	return redefinition(short, long)
}

// parseSingleFlag is called when there is a single flag, and it is not known
// whether the flag is short or long. If validates the flag and ensures it is
// not a duplicate.
func parseSingleFlag(flag string) (rune, string, error) {
	if flag == "" {
		return utf8.RuneError, "", errors.New("cannot use empty flag string")
	}

	var firstRune rune
	var bufIndex, runeCount int

	// Ensure all bytes are valid runes
	buf := []byte(flag)

	for bufIndex < len(buf) {
		thisRune, runeSize := utf8.DecodeRune(buf[bufIndex:])
		if thisRune == utf8.RuneError {
			return thisRune, "", fmt.Errorf("cannot use flag with invalid rune: %q", flag)
		}
		if runeCount == 0 {
			if thisRune == '-' {
				return thisRune, "", fmt.Errorf("cannot use flag that starts with a hyphen: %q", flag)
			}
			firstRune = thisRune
		}
		runeCount++
		bufIndex += runeSize
	}

	if runeCount == 1 {
		return firstRune, "", redefinition(firstRune, "")
	}
	return utf8.RuneError, flag, redefinition(utf8.RuneError, flag)
}

// option is list of methods any concrete option needs to have for use by
// parser.
type option interface {
	Default() interface{} // default value of command line option
	Description() string  // describes the command line option
	Long() string         // long flag
	NextSlurp() slurpType // next state for state machine
	Short() rune          // short flag
}

type optionBool struct {
	pv          *bool
	long        string
	description string
	short       rune
	def         bool
}

func (o optionBool) Default() interface{} { return o.def }
func (o optionBool) Description() string  { return o.description }
func (o optionBool) Long() string         { return o.long }
func (o optionBool) NextSlurp() slurpType { return nothingToSlurp }
func (o optionBool) Short() rune          { return o.short }

// Bool returns a pointer to a bool command line option, allowing for either a
// short or a long flag. If both are desired, use the BoolP function.
func Bool(flag string, value bool, description string) *bool {
	var v bool
	BoolVar(&v, flag, value, description)
	return &v
}

// BoolVar binds an existing boolean variable to a flag, allowing for either a
// short or a long flag. If both are desired, use the BoolVarP function.
func BoolVar(pv *bool, flag string, value bool, description string) {
	short, long, err := parseSingleFlag(flag)
	if err != nil {
		panic(err)
	}
	*pv = value
	flags = append(flags, &optionBool{
		description: description,
		long:        long,
		short:       short,
		pv:          pv,
		def:         value,
	})
}

// BoolP returns a pointer to a bool command line option, allowing for both a
// short and a long flag.
func BoolP(short rune, long string, value bool, description string) *bool {
	var v bool
	BoolVarP(&v, short, long, value, description)
	return &v
}

// BoolVarP binds an existing boolean variable to a flag, allowing for both a
// short and a long flag.
func BoolVarP(pv *bool, short rune, long string, value bool, description string) {
	if err := parseShortAndLongFlag(short, long); err != nil {
		panic(err)
	}
	*pv = value
	flags = append(flags, &optionBool{
		description: description,
		long:        long,
		short:       short,
		pv:          pv,
		def:         value,
	})
}

type optionDuration struct {
	pv          *time.Duration
	description string
	long        string
	short       rune
	def         time.Duration
}

func (o optionDuration) Default() interface{} { return o.def }
func (o optionDuration) Description() string  { return o.description }
func (o optionDuration) Long() string         { return o.long }
func (o optionDuration) NextSlurp() slurpType { return slurpDuration }
func (o optionDuration) Short() rune          { return o.short }

// Duration returns a pointer to a time.Duration command line option, allowing
// for either a short or a long flag. If both are desired, use the DurationP
// function.
func Duration(flag string, value time.Duration, description string) *time.Duration {
	var v time.Duration
	DurationVar(&v, flag, value, description)
	return &v
}

// DurationVar binds an existing time.Duration variable to a flag, allowing for
// either a short or a long flag. If both are desired, use the DurationVarP
// function.
func DurationVar(pv *time.Duration, flag string, value time.Duration, description string) {
	short, long, err := parseSingleFlag(flag)
	if err != nil {
		panic(err)
	}
	*pv = value
	flags = append(flags, &optionDuration{
		description: description,
		long:        long,
		short:       short,
		pv:          pv,
		def:         value,
	})
}

// DurationP returns a pointer to a time.Duration command line option, allowing
// for both a short and a long flag.
func DurationP(short rune, long string, value time.Duration, description string) *time.Duration {
	var v time.Duration
	DurationVarP(&v, short, long, value, description)
	return &v
}

// DurationVarP binds an existing time.Duration variable to a flag, allowing for
// both a short and a long flag.
func DurationVarP(pv *time.Duration, short rune, long string, value time.Duration, description string) {
	if err := parseShortAndLongFlag(short, long); err != nil {
		panic(err)
	}
	*pv = value
	flags = append(flags, &optionDuration{
		description: description,
		long:        long,
		short:       short,
		pv:          pv,
		def:         value,
	})
}

type optionFloat struct {
	pv          *float64
	description string
	short       rune
	long        string
	def         float64
}

func (o optionFloat) Default() interface{} { return o.def }
func (o optionFloat) Description() string  { return o.description }
func (o optionFloat) Long() string         { return o.long }
func (o optionFloat) NextSlurp() slurpType { return slurpFloat }
func (o optionFloat) Short() rune          { return o.short }

// Float returns a pointer to a float64 command line option, allowing for either
// a short or a long flag. If both are desired, use the FloatP function.
func Float(flag string, value float64, description string) *float64 {
	var v float64
	FloatVar(&v, flag, value, description)
	return &v
}

// FloatVar binds an existing float64 variable to a flag, allowing for either a
// short or a long flag. If both are desired, use the FloatVarP function.
func FloatVar(pv *float64, flag string, value float64, description string) {
	short, long, err := parseSingleFlag(flag)
	if err != nil {
		panic(err)
	}
	*pv = value
	flags = append(flags, &optionFloat{
		description: description,
		long:        long,
		short:       short,
		pv:          pv,
		def:         value,
	})
}

// FloatP returns a pointer to a float64 command line option, allowing for both
// a short and a long flag.
func FloatP(short rune, long string, value float64, description string) *float64 {
	var v float64
	FloatVarP(&v, short, long, value, description)
	return &v
}

// FloatVarP binds an existing float64 variable to a flag, allowing for both a
// short and a long flag.
func FloatVarP(pv *float64, short rune, long string, value float64, description string) {
	if err := parseShortAndLongFlag(short, long); err != nil {
		panic(err)
	}
	*pv = value
	flags = append(flags, &optionFloat{
		description: description,
		long:        long,
		short:       short,
		pv:          pv,
		def:         value,
	})
}

type optionInt struct {
	pv          *int
	description string
	short       rune
	long        string
	def         int
}

func (o optionInt) Default() interface{} { return o.def }
func (o optionInt) Description() string  { return o.description }
func (o optionInt) Long() string         { return o.long }
func (o optionInt) NextSlurp() slurpType { return slurpInt }
func (o optionInt) Short() rune          { return o.short }

// Int returns a pointer to a int command line option, allowing for either a
// short or a long flag. If both are desired, use the IntP function.
func Int(flag string, value int, description string) *int {
	var v int
	IntVar(&v, flag, value, description)
	return &v
}

// IntVar binds an existing int variable to a flag, allowing for either a short
// or a long flag. If both are desired, use the IntVarP function.
func IntVar(pv *int, flag string, value int, description string) {
	short, long, err := parseSingleFlag(flag)
	if err != nil {
		panic(err)
	}
	*pv = value
	flags = append(flags, &optionInt{
		description: description,
		long:        long,
		short:       short,
		pv:          pv,
		def:         value,
	})
}

// IntP returns a pointer to a int command line option, allowing for both a
// short and a long flag.
func IntP(short rune, long string, value int, description string) *int {
	var v int
	IntVarP(&v, short, long, value, description)
	return &v
}

// IntVarP binds an existing int variable to a flag, allowing for both a short
// and a long flag.
func IntVarP(pv *int, short rune, long string, value int, description string) {
	if err := parseShortAndLongFlag(short, long); err != nil {
		panic(err)
	}
	*pv = value
	flags = append(flags, &optionInt{
		description: description,
		long:        long,
		short:       short,
		pv:          pv,
		def:         value,
	})
}

type optionInt64 struct {
	pv          *int64
	description string
	short       rune
	long        string
	def         int64
}

func (o optionInt64) Default() interface{} { return o.def }
func (o optionInt64) Description() string  { return o.description }
func (o optionInt64) Long() string         { return o.long }
func (o optionInt64) NextSlurp() slurpType { return slurpInt64 }
func (o optionInt64) Short() rune          { return o.short }

// Int64 returns a pointer to a int64 command line option, allowing for either a
// short or a long flag. If both are desired, use the Int64P function.
func Int64(flag string, value int64, description string) *int64 {
	var v int64
	Int64Var(&v, flag, value, description)
	return &v
}

// Int64Var binds an existing int64 variable to a flag, allowing for either a
// short or a long flag. If both are desired, use the Int64VarP function.
func Int64Var(pv *int64, flag string, value int64, description string) {
	short, long, err := parseSingleFlag(flag)
	if err != nil {
		panic(err)
	}
	*pv = value
	flags = append(flags, &optionInt64{
		description: description,
		long:        long,
		short:       short,
		pv:          pv,
		def:         value,
	})
}

// Int64P returns a pointer to a int64 command line option, allowing for both a
// short and a long flag.
func Int64P(short rune, long string, value int64, description string) *int64 {
	var v int64
	Int64VarP(&v, short, long, value, description)
	return &v
}

// Int64VarP binds an existing int64 variable to a flag, allowing for both a
// short and a long flag.
func Int64VarP(pv *int64, short rune, long string, value int64, description string) {
	if err := parseShortAndLongFlag(short, long); err != nil {
		panic(err)
	}
	*pv = value
	flags = append(flags, &optionInt64{
		description: description,
		long:        long,
		short:       short,
		pv:          pv,
		def:         value,
	})
}

type optionUint struct {
	pv          *uint
	description string
	short       rune
	long        string
	def         uint
}

func (o optionUint) Default() interface{} { return o.def }
func (o optionUint) Description() string  { return o.description }
func (o optionUint) Long() string         { return o.long }
func (o optionUint) NextSlurp() slurpType { return slurpUint }
func (o optionUint) Short() rune          { return o.short }

// Uint returns a pouinter to a uint command line option, allowing for either a
// short or a long flag. If both are desired, use the UintP function.
func Uint(flag string, value uint, description string) *uint {
	var v uint
	UintVar(&v, flag, value, description)
	return &v
}

// UintVar binds an existing uint variable to a flag, allowing for either a short
// or a long flag. If both are desired, use the UintVarP function.
func UintVar(pv *uint, flag string, value uint, description string) {
	short, long, err := parseSingleFlag(flag)
	if err != nil {
		panic(err)
	}
	*pv = value
	flags = append(flags, &optionUint{
		description: description,
		long:        long,
		short:       short,
		pv:          pv,
		def:         value,
	})
}

// UintP returns a pouinter to a uint command line option, allowing for both a
// short and a long flag.
func UintP(short rune, long string, value uint, description string) *uint {
	var v uint
	UintVarP(&v, short, long, value, description)
	return &v
}

// UintVarP binds an existing uint variable to a flag, allowing for both a short
// and a long flag.
func UintVarP(pv *uint, short rune, long string, value uint, description string) {
	if err := parseShortAndLongFlag(short, long); err != nil {
		panic(err)
	}
	*pv = value
	flags = append(flags, &optionUint{
		description: description,
		long:        long,
		short:       short,
		pv:          pv,
		def:         value,
	})
}

type optionUint64 struct {
	pv          *uint64
	description string
	short       rune
	long        string
	def         uint64
}

func (o optionUint64) Default() interface{} { return o.def }
func (o optionUint64) Description() string  { return o.description }
func (o optionUint64) Long() string         { return o.long }
func (o optionUint64) NextSlurp() slurpType { return slurpUint64 }
func (o optionUint64) Short() rune          { return o.short }

// Uint64 returns a pointer to a uint64 command line option, allowing for either a
// short or a long flag. If both are desired, use the Uint64P function.
func Uint64(flag string, value uint64, description string) *uint64 {
	var v uint64
	Uint64Var(&v, flag, value, description)
	return &v
}

// Uint64Var binds an existing uint64 variable to a flag, allowing for either a
// short or a long flag. If both are desired, use the Uint64VarP function.
func Uint64Var(pv *uint64, flag string, value uint64, description string) {
	short, long, err := parseSingleFlag(flag)
	if err != nil {
		panic(err)
	}
	*pv = value
	flags = append(flags, &optionUint64{
		description: description,
		long:        long,
		short:       short,
		pv:          pv,
		def:         value,
	})
}

// Uint64P returns a pointer to a uint64 command line option, allowing for both a
// short and a long flag.
func Uint64P(short rune, long string, value uint64, description string) *uint64 {
	var v uint64
	Uint64VarP(&v, short, long, value, description)
	return &v
}

// Uint64VarP binds an existing uint64 variable to a flag, allowing for both a
// short and a long flag.
func Uint64VarP(pv *uint64, short rune, long string, value uint64, description string) {
	if err := parseShortAndLongFlag(short, long); err != nil {
		panic(err)
	}
	*pv = value
	flags = append(flags, &optionUint64{
		description: description,
		long:        long,
		short:       short,
		pv:          pv,
		def:         value,
	})
}

type optionString struct {
	pv          *string
	description string
	short       rune
	long        string
	def         string
}

func (o optionString) Default() interface{} { return o.def }
func (o optionString) Description() string  { return o.description }
func (o optionString) Long() string         { return o.long }
func (o optionString) NextSlurp() slurpType { return slurpString }
func (o optionString) Short() rune          { return o.short }

// String returns a postringer to a string command line option, allowing for either a
// short or a long flag. If both are desired, use the StringP function.
func String(flag string, value string, description string) *string {
	var v string
	StringVar(&v, flag, value, description)
	return &v
}

// StringVar binds an existing string variable to a flag, allowing for either a short
// or a long flag. If both are desired, use the StringVarP function.
func StringVar(pv *string, flag string, value string, description string) {
	short, long, err := parseSingleFlag(flag)
	if err != nil {
		panic(err)
	}
	*pv = value
	flags = append(flags, &optionString{
		description: description,
		long:        long,
		short:       short,
		pv:          pv,
		def:         value,
	})
}

// StringP returns a postringer to a string command line option, allowing for both a
// short and a long flag.
func StringP(short rune, long string, value string, description string) *string {
	var v string
	StringVarP(&v, short, long, value, description)
	return &v
}

// StringVarP binds an existing string variable to a flag, allowing for both a short
// and a long flag.
func StringVarP(pv *string, short rune, long string, value string, description string) {
	if err := parseShortAndLongFlag(short, long); err != nil {
		panic(err)
	}
	*pv = value
	flags = append(flags, &optionString{
		description: description,
		long:        long,
		short:       short,
		pv:          pv,
		def:         value,
	})
}
