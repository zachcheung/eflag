package eflag

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Flag struct {
	Name    string
	Flag    *flag.Flag
	Env     string
	Changed bool
	p       interface{}
}

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
		p:    p,
	}
}

type Flags []*Flag

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
			// explicitly set flag has the highest precedence
			continue
		}

		if f.Env == "-" {
			continue
		}

		if prefix != "" {
			if !strings.HasPrefix(f.Env, prefix) {
				f.Env = prefix + f.Env
			}
		}
		if v := os.Getenv(f.Env); v != "" {
			if err := setFromStr(f.p, v); err != nil {
				fmt.Printf("invalid value %#v for env %s: parse error\n", v, f.Env)
				os.Exit(2)
			}
		}
	}
}

func setFromStr(ptr interface{}, s string) error {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("expected a pointer, got %T", ptr)
	}

	v = v.Elem()
	switch v.Kind() {
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		v.SetBool(boolVal)
	case reflect.String:
		v.SetString(s)
	case reflect.Int, reflect.Int64:
		intVal, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(intVal)
	case reflect.Uint, reflect.Uint64:
		uintVal, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(uintVal)
	case reflect.Float64:
		floatVal, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		v.SetFloat(floatVal)
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(time.Duration(0)) {
			durationVal, err := time.ParseDuration(s)
			if err != nil {
				return err
			}
			v.SetInt(int64(durationVal))
		} else {
			return fmt.Errorf("unsupported struct type")
		}
	default:
		return fmt.Errorf("unsupported type %v", v.Kind())
	}

	return nil
}

var flags Flags

func Var(p interface{}, name string, value interface{}, usage, env string) {
	flags = append(flags, newFlag(p, name, value, usage, env))
}

func Parse(prefix string) {
	flags.Parse(prefix)
}
