package eflag

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type Flag struct {
	Name         string
	DefaultValue interface{}
	Usage        string
	Value        interface{}
	Changed      bool
}

func (f *Flag) Init() {
	switch f.DefaultValue.(type) {
	case bool:
		flag.BoolVar(f.Value.(*bool), f.Name, f.DefaultValue.(bool), f.Usage)
	case time.Duration:
		flag.DurationVar(f.Value.(*time.Duration), f.Name, f.DefaultValue.(time.Duration), f.Usage)
	case float64:
		flag.Float64Var(f.Value.(*float64), f.Name, f.DefaultValue.(float64), f.Usage)
	case int:
		flag.IntVar(f.Value.(*int), f.Name, f.DefaultValue.(int), f.Usage)
	case int64:
		flag.Int64Var(f.Value.(*int64), f.Name, f.DefaultValue.(int64), f.Usage)
	case string:
		flag.StringVar(f.Value.(*string), f.Name, f.DefaultValue.(string), f.Usage)
	case uint:
		flag.UintVar(f.Value.(*uint), f.Name, f.DefaultValue.(uint), f.Usage)
	case uint64:
		flag.Uint64Var(f.Value.(*uint64), f.Name, f.DefaultValue.(uint64), f.Usage)
	default:
		fmt.Printf("invalid type: %T", f.DefaultValue)
		os.Exit(1)
	}
}

type Flags []*Flag

func (ff Flags) Parse() {
	m := make(map[string]*Flag)
	for _, f := range ff {
		f.Init()
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
	flags = append(flags, &Flag{
		Value:        p,
		Name:         name,
		DefaultValue: value,
		Usage:        usage,
	})
}

func Parse() {
	flags.Parse()
}
