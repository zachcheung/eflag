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
		// don't read from env var
	} else {
		env = strings.ToUpper(env)
	}

	return &Flag{
		Name: name,
		Flag: flag.Lookup(name),
		Env:  env,
	}
}

// Flags is a slice of Flag pointers.
type Flags []*Flag

// Parse parses command-line flags and sets values from environment variables.
func (ff Flags) Parse(prefix string) {
	m := make(map[string]*Flag)
	for _, f := range ff {
		m[f.Name] = f
	}
	flag.Parse()

	flag.Visit(func(f *flag.Flag) {
		if v, ok := m[f.Name]; ok {
			v.Changed = true
		}
	})

	if prefix != "" && !strings.HasSuffix(prefix, "_") {
		prefix += "_"
	}
	for _, f := range ff {
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

var flags Flags

// Var registers a flag and associates it with a variable, environment variable, and usage description.
func Var(p interface{}, name string, value interface{}, usage, env string) {
	flags = append(flags, newFlag(p, name, value, usage, env))
}

// Parse parses all registered flags.
func Parse(prefix string) {
	flags.Parse(strings.ToUpper(prefix))
}
