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
	Name    string     // Name of the flag
	Flag    *flag.Flag // Reference to the underlying flag
	Env     string     // Environment variable associated with the flag
	Changed bool       // Indicates whether the flag has been changed
}

// newFlag creates a new Flag based on the provided parameters.
// It associates a flag with a variable, a default value, a usage description,
// and an optional environment variable.
func newFlag(p interface{}, name string, value interface{}, usage, env string) *Flag {
	switch value.(type) {
	case bool:
		flag.BoolVar(p.(*bool), name, value.(bool), usage)
	case time.Duration:
		flag.DurationVar(p.(*time.Duration), name, value.(time.Duration), usage)
	case float64:
		flag.Float64Var(p.(*float64), name, value.(float64), usage)
	case int:
		flag.IntVar(p.(*int), name, value.(int), usage)
	case int64:
		flag.Int64Var(p.(*int64), name, value.(int64), usage)
	case string:
		flag.StringVar(p.(*string), name, value.(string), usage)
	case uint:
		flag.UintVar(p.(*uint), name, value.(uint), usage)
	case uint64:
		flag.Uint64Var(p.(*uint64), name, value.(uint64), usage)
	default:
		fmt.Printf("invalid type: %T\n", value)
		os.Exit(1)
	}

	if env == "" {
		env = strings.ToUpper(name)
	} else if env == "-" {
		// Don't read from env var
	} else {
		env = strings.ToUpper(env)
	}

	return &Flag{
		Name: name,
		Flag: flag.Lookup(name),
		Env:  env,
	}
}

// FlagSet represents a set of flags and provides methods for parsing them.
type FlagSet struct {
	Flags  []*Flag // List of flags in the set
	prefix string  // Prefix for environment variables associated with flags
}

// SetPrefix set environment variable prefix.
func (ff *FlagSet) SetPrefix(prefix string) {
	if prefix != "" {
		prefix = strings.ToUpper(prefix)
		if !strings.HasSuffix(prefix, "_") {
			prefix += "_"
		}
	}
	ff.prefix = prefix
}

// Parse parses command-line flags and sets values from environment variables.
func (ff *FlagSet) Parse() {
	m := make(map[string]*Flag)
	for _, f := range ff.Flags {
		m[f.Name] = f
	}
	flag.Parse()

	flag.Visit(func(f *flag.Flag) {
		if v, ok := m[f.Name]; ok {
			v.Changed = true
		}
	})

	ff.parse()
}

// ReParse re-parses flags. This can be useful in scenarios where the
// environment variables have changed, and you want to update the flag values.
func (ff *FlagSet) ReParse() {
	ff.parse()
}

// parse sets flag values from environment variables and respects
// the precedence of explicitly set flags over environment variables.
func (ff *FlagSet) parse() {
	prefix := ff.prefix
	for _, f := range ff.Flags {
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
}

// CommandLine represents the default set of flags that are parsed
// when the package is imported. It provides a convenient global
// instance for managing flags.
var CommandLine = &FlagSet{}

// Var registers a flag and associates it with a variable, environment
// variable, and usage description. It is recommended to use this function
// during the initialization phase to register flags.
//
// The env parameter determines the association with an environment variable:
//   - When env is an empty string (""): The environment variable will be derived
//     from the flag name, optionally with a prefix, to check for the environment
//     variable.
//   - When env is "-": The flag will not be associated with any environment
//     variable, and environment variable checking will be ignored.
func Var(p interface{}, name string, value interface{}, usage, env string) {
	CommandLine.Flags = append(CommandLine.Flags, newFlag(p, name, value, usage, env))
}

// SetPrefix set environment variable prefix.
func SetPrefix(prefix string) {
	CommandLine.SetPrefix(prefix)
}

// Parse parses all registered flags.
func Parse() {
	CommandLine.Parse()
}

// ReParse re-parses all registered flags. This is useful when
// environment variables have changed, and you want to update the flag values.
func ReParse() {
	CommandLine.ReParse()
}
