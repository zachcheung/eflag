package eflag

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type Flag struct {
	Name    string
	Flag    *flag.Flag
	Changed bool
}

func newFlag(p interface{}, name string, value interface{}, usage string) *Flag {
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
		fmt.Printf("invalid type: %T", value)
		os.Exit(1)
	}

	return &Flag{Name: name, Flag: flag.Lookup(name)}

}

type Flags []*Flag

func (ff Flags) Parse() {
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
}

var flags Flags

func Var(p interface{}, name string, value interface{}, usage string) {
	flags = append(flags, newFlag(p, name, value, usage))
}

func Parse() {
	flags.Parse()
}
