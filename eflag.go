// Package eflag provides an extended flag package with enhanced features
// including the ability to set flag values from environment variables and
// a more convenient way to manage and parse multiple flags.
package eflag

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// Flag represents a command-line flag and its associated information.
type Flag struct {
	p interface{}
	*flag.Flag
	Name    string // Name of the flag
	Env     string // Environment variable associated with the flag
	Changed bool   // Indicates whether the flag has been changed
}

// newFlag creates a new Flag based on the provided parameters.
// It associates a flag with a variable, a default value, a usage description,
// and an optional environment variable.
func newFlag(fs *flag.FlagSet, p interface{}, name string, value interface{}, usage, env string) *Flag {
	switch p.(type) {
	case *bool:
		fs.BoolVar(p.(*bool), name, value.(bool), usage)
	case *time.Duration:
		fs.DurationVar(p.(*time.Duration), name, value.(time.Duration), usage)
	case *float64:
		fs.Float64Var(p.(*float64), name, value.(float64), usage)
	case *int:
		fs.IntVar(p.(*int), name, value.(int), usage)
	case *int64:
		fs.Int64Var(p.(*int64), name, value.(int64), usage)
	case *string:
		fs.StringVar(p.(*string), name, value.(string), usage)
	case *uint:
		fs.UintVar(p.(*uint), name, value.(uint), usage)
	case *uint64:
		fs.Uint64Var(p.(*uint64), name, value.(uint64), usage)
	case *StringList:
		fs.StringVar(&p.(*StringList).p, name, value.(string), usage)
	default:
		fmt.Printf("invalid type: %T\n", p)
		os.Exit(1)
	}

	if env == "" {
		env = MixedCapsToScreamingSnake(name)
	} else if env == "-" {
		// Don't read from env var
	} else {
		env = strings.ToUpper(env)
	}

	return &Flag{
		p:    p,
		Flag: fs.Lookup(name),
		Name: name,
		Env:  env,
	}
}

// ErrorHandling defines how FlagSet.Parse behaves if the parse fails.
type ErrorHandling flag.ErrorHandling

const (
	ContinueOnError ErrorHandling = iota // Return a descriptive error.
	ExitOnError                          // Call os.Exit(2) or for -h/-help Exit(0).
	PanicOnError                         // Call panic with a descriptive error.
)

// FlagSet represents a set of flags and provides methods for parsing them.
type FlagSet struct {
	*flag.FlagSet
	formal map[string]*Flag
	prefix string // Prefix for environment variables associated with flags
}

// NewFlagSet returns a new, empty flag set with the specified name and
// error handling property. If the name is not empty, it will be printed
// in the default usage message and in error messages.
func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet {
	flagSet := flag.NewFlagSet(name, flag.ErrorHandling(errorHandling))
	f := &FlagSet{
		FlagSet: flagSet,
	}
	return f
}

// SetPrefix set environment variable prefix.
func (fs *FlagSet) SetPrefix(prefix string) {
	if prefix != "" {
		prefix = strings.ToUpper(prefix)
		if !strings.HasSuffix(prefix, "_") {
			prefix += "_"
		}
	}
	fs.prefix = prefix
}

// Var registers a flag and associates it with a variable, environment
// variable, and usage description. It is recommended to use this function
// during the initialization phase to register flags.
//
// The env parameter determines the association with an environment variable:
//   - When env is an empty string (""): The environment variable name will be derived
//     from the flag name by converting it to uppercase and replacing any camel case with underscores.
//     For example, if the flag name is "mixedCaps", the derived environment variable name will be "MIXED_CAPS".
//     An optional prefix can be added to the environment variable name by using SetPrefix() function.
//   - When env is "-": The flag will not be associated with any environment
//     variable, and environment variable checking will be ignored.
func (fs *FlagSet) Var(p interface{}, name string, value interface{}, usage, env string) {
	if fs.formal == nil {
		fs.formal = make(map[string]*Flag)
	}
	fs.formal[name] = newFlag(fs.FlagSet, p, name, value, usage, env)
}

// Parse parses command-line flags and sets values from environment variables.
func (fs *FlagSet) Parse(arguments []string) {
	fs.FlagSet.Parse(arguments)

	fs.FlagSet.Visit(func(f *flag.Flag) {
		// Visit() visits only those flags that have been set.
		fs.formal[f.Name].Changed = true
	})

	fs.parse()
}

// ReParse re-parses flags. This can be useful in scenarios where the
// environment variables have changed, and you want to update the flag values.
func (fs *FlagSet) ReParse() {
	fs.parse()
}

// parse sets flag values from environment variables and respects
// the precedence of explicitly set flags over environment variables.
func (fs *FlagSet) parse() {
	prefix := fs.prefix
	for _, f := range fs.formal {
		if f.Changed {
			// Explicitly set flag has the highest precedence
			continue
		}

		if f.Env == "-" {
			continue
		}

		if prefix != "" && !strings.HasPrefix(f.Env, prefix) {
			f.Env = prefix + f.Env
		}
		if v := os.Getenv(f.Env); v != "" {
			if err := f.Flag.Value.Set(v); err != nil {
				fmt.Printf("invalid value %#v for env %s: parse error\n", v, f.Env)
				os.Exit(2)
			}
		}
	}

	for _, f := range fs.formal {
		switch f.p.(type) {
		case *StringList:
			f.p.(*StringList).setValue()
		}
	}
}

// CommandLine represents the default set of flags that are parsed
// when the package is imported. It provides a convenient global
// instance for managing flags.
var CommandLine = NewFlagSet(os.Args[0], ExitOnError)

// Var registers a command-line flag and associates it with a variable, environment.
func Var(p interface{}, name string, value interface{}, usage, env string) {
	CommandLine.Var(p, name, value, usage, env)
}

// SetPrefix set environment variable prefix.
func SetPrefix(prefix string) {
	CommandLine.SetPrefix(prefix)
}

// Parse parses all registered flags.
func Parse() {
	CommandLine.Parse(os.Args[1:])
}

// ReParse re-parses all registered flags. This is useful when
// environment variables have changed, and you want to update the flag values.
func ReParse() {
	CommandLine.ReParse()
}

// Func defines a flag with the specified name and usage string.
// Each time the flag is seen, fn is called with the value of the flag.
// If fn returns a non-nil error, it will be treated as a flag value parsing error.
func Func(name, usage string, fn func(string) error) {
	CommandLine.Func(name, usage, fn)
}

// Set sets the value of the named command-line flag.
func Set(name, value string) error {
	return CommandLine.Set(name, value)
}

// PrintDefaults prints, to standard error unless configured otherwise,
// a usage message showing the default settings of all defined
// command-line flags.
func PrintDefaults() {
	CommandLine.PrintDefaults()
}

// NFlag returns the number of command-line flags that have been set.
func NFlag() int { return CommandLine.NFlag() }

// Arg returns the i'th command-line argument. Arg(0) is the first remaining argument
// after flags have been processed. Arg returns an empty string if the
// requested element does not exist.
func Arg(i int) string {
	return CommandLine.Arg(i)
}

// NArg is the number of arguments remaining after flags have been processed.
func NArg() int { return CommandLine.NArg() }

// Args returns the non-flag command-line arguments.
func Args() []string { return CommandLine.Args() }

// Parsed reports whether the command-line flags have been parsed.
func Parsed() bool {
	return CommandLine.Parsed()
}
