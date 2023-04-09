package syringe

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"text/template"
)

const (
	intString     = "int"
	floatString   = "float"
	numberString  = "number"
	listString    = "list"
	mapString     = "map"
	unknownString = "unknown"
)

// Syringe is the receiver of the template functions injector.
type Syringe struct {
	expander string
	version  string
	logUsed  bool
	builtins []builtinFunc
}

// Opts are the options for New.
type Opts struct {
	Expander string
	Version  string
}

type builtinFunc struct {
	function  interface{}
	longname  string
	shortname string
	usage     string
}

// New returns an initialized Syringe.
func New(o *Opts) *Syringe {
	s := &Syringe{
		expander: o.Expander,
		version:  o.Version,
	}
	s.builtins = []builtinFunc{
		// General
		{
			function:  s.Expander,
			longname:  ".Gtpl.Expander",
			shortname: "expander",
			usage:     "{{ expander }} - the name of this template expander",
		},
		{
			function:  s.Version,
			longname:  ".Gtpl.Version",
			shortname: "version",
			usage:     "{{ version }} - the version of this template expander",
		},
		{
			function:  s.Log,
			longname:  ".Gtpl.Log",
			shortname: "log",
			usage:     `{{ log "some" "info" }} - sends args to the log`,
		},
		{
			function:  s.Die,
			longname:  ".Gtpl.Die",
			shortname: "die",
			usage:     `{{ die "some" "info" }} - prints args, logs them if logging was used, stops`,
		},
		{
			function:  s.Assert,
			longname:  ".Gtpl.Assert",
			shortname: "assert",
			usage:     `asserts a condition and stops if not met: {{ assert (len $list) gt 0) "list is empty!" }}`,
		},

		// Lists
		{
			function:  s.List,
			longname:  ".Gtpl.List",
			shortname: "list",
			usage:     `{{ $list := list "a" "b" "c" }} - creates a list`,
		},
		{
			function:  s.HasElement,
			longname:  ".Gtpl.HasElement",
			shortname: "haselement",
			usage:     `{{ if (haselement $list "a") }} 'a' occurs in the list {{ end }}`,
		},
		{
			function:  s.IndexOf,
			longname:  ".Gtpl.IndexOf",
			shortname: "indexof",
			usage:     `'a' occurs at index {{ indexof $list "a" }} in the list`,
		},
		{
			function:  s.AddElements,
			longname:  ".Gtpl.AddElements",
			shortname: "addelements",
			usage:     `{{ $newlist := (addelements $list "d" "e") }} - creates a new list with added element`,
		},

		// Maps
		{
			function:  s.Map,
			longname:  ".Gtpl.Map",
			shortname: "map",
			usage:     `{{ $map := map "cat" "meow" "dog" "woof" }} - creates a map`,
		},
		{
			function:  s.HasKey,
			longname:  ".Gtpl.HasKey",
			shortname: "haskey",
			usage:     `{{ if haskey $map "cat" }} yes {{ else }} no {{ end }} - tests whether a key is in a map`,
		},
		{
			function:  s.GetVal,
			longname:  ".Gtpl.GetVal",
			shortname: "getval",
			usage:     `a cat says {{ get $map "cat" }} - gets a value from a map, "" if absent`,
		},
		{
			function:  s.SetKeyVal,
			longname:  ".Gtpl.SetKeyVal",
			shortname: "setkeyval",
			usage:     `{{ set $map "frog" "ribbit" }} - adds a key/value pair to a map`,
		},

		// Types
		{
			function:  s.Type,
			longname:  ".Gtpl.Type",
			shortname: "type",
			usage:     `expands to "int", "float", "list" or "map": {{ $t := type $map }} {{ if $t ne "map" }} something is very wrong {{ end }}`,
		},
		{
			function:  s.IsInt,
			longname:  ".Gtpl.IsInt",
			shortname: "isint",
			usage:     `true when its argument is an integer`,
		},
		{
			function:  s.IsFloat,
			longname:  ".Gtpl.IsFloat",
			shortname: "isfloat",
			usage:     `true when its argument is a float`,
		},
		{
			function:  s.IsNumber,
			longname:  ".Gtpl.IsNumber",
			shortname: "isnumber",
			usage:     `true when its argument is an int or a float`,
		},
		{
			function:  s.IsList,
			longname:  ".Gtpl.IsList",
			shortname: "islist",
			usage:     `true when its argument is a list (or a slice)`,
		},
		{
			function:  s.IsMap,
			longname:  ".Gtpl.IsMap",
			shortname: "ismap",
			usage:     `true when its argument is a map`,
		},

		// Arithmetic / misc
		{
			function:  s.Add,
			longname:  ".Gtpl.Add",
			shortname: "add",
			usage:     `21 + 21 is {{ add (21 21) }}`,
		},
		{
			function:  s.Sub,
			longname:  ".Gtpl.Sub",
			shortname: "sub",
			usage:     `42 - 2 = {{ sub 42 2}}`,
		},
		{
			function:  s.Mul,
			longname:  ".Gtpl.Mul",
			shortname: "mul",
			usage:     `7 * 4 = {{ mul 7 4 }}`,
		},
		{
			function:  s.Div,
			longname:  ".Gtpl.Div",
			shortname: "div",
			usage:     `42 / 4 = {{ div 42 4 }}`,
		},
		{
			function:  s.Loop,
			longname:  ".Gtpl.Loop",
			shortname: "loop",
			usage:     `1 up to and including 10: {{ range $i := loop 1 11 }} {{ $i }} {{ end }}`,
		},
	}
	return s
}

// FlatNamespace returns a `template.FuncMap` that can be passed to text/template so that shorthand
// builtins can be used.
func (s *Syringe) FlatNamespace() template.FuncMap {
	fmap := template.FuncMap{}
	for _, b := range s.builtins {
		fmap[b.shortname] = b.function
	}
	return fmap
}

// Overview returns a short usage text of the builtins.
func (s *Syringe) Overview() string {
	str := []string{}
	for _, b := range s.builtins {
		str = append(str,
			fmt.Sprintf("%v (long name: %v)", b.shortname, b.longname),
			"  "+b.usage,
			"")
	}
	return strings.Join(str, "\n")
}

// Builtin functions.
// Remember to update the above info when adding/modifying!

/* General */

// Expander is the builtin returning the name of the expander program.
func (s *Syringe) Expander() string {
	return s.expander
}

// Version is the builtin returning the version of the expander program.
func (s *Syringe) Version() string {
	return s.version
}

// Log is the builtin that logs information using the `log.Print` function.
func (s *Syringe) Log(args ...interface{}) string {
	s.logUsed = true
	msg := fmt.Sprint(args...)
	log.Printf("%s: %s", s.expander, msg)
	return ""
}

// Die is the builtin that stops execution. If previous `Log` invocations occurred, then the
// the reason for stopping is logged, else, the reason is shown on `os.Stderr`.
func (s *Syringe) Die(args ...interface{}) string {
	msg := fmt.Sprint(args...)
	if s.logUsed {
		s.Log(msg)
	} else {
		fmt.Fprintf(os.Stderr, "%s: %s\n", s.expander, msg)
	}
	os.Exit(1)
	return ""
}

// Assert is the builtin that ensures a condition.
func (s *Syringe) Assert(cond bool, args ...interface{}) string {
	if !cond {
		s.Die(args...)
	}
	return ""
}

/* List related */

// List is the builtin that returns a list.
func (s *Syringe) List(args ...interface{}) []interface{} {
	return args
}

// HasElement is the builtin that checks whether a list contains an element.
func (s *Syringe) HasElement(list []interface{}, el interface{}) bool {
	for _, e := range list {
		if e == el {
			return true
		}
	}
	return false
}

// IndexOf returns the index of an element in a list, or -1.
func (s *Syringe) IndexOf(list []interface{}, el interface{}) int {
	for i, e := range list {
		if e == el {
			return i
		}
	}
	return -1
}

// AddElements is the builtin that adds elements to a list.
func (s *Syringe) AddElements(list []interface{}, els ...interface{}) []interface{} {
	list = append(list, els...)
	return list
}

/* Map related */

// Map is the builtin that returns a map.
func (s *Syringe) Map(args ...interface{}) map[interface{}]interface{} {
	out := map[interface{}]interface{}{}
	for i := 0; i < len(args); i += 2 {
		out[args[i]] = args[i+1]
	}
	return out
}

// HasKey is the builtin that checks whether a map contains a key.
func (s *Syringe) HasKey(m map[interface{}]interface{}, key interface{}) bool {
	_, ok := m[key]
	return ok
}

// GetVal is the builtin that returns a value from a map given a key, or "".
func (s *Syringe) GetVal(m map[interface{}]interface{}, key interface{}) interface{} {
	val, ok := m[key]
	if ok {
		return val
	}
	return ""
}

// SetKeyVal is the builtin that sets a value for a key in a map, or adds it.
func (s *Syringe) SetKeyVal(m map[interface{}]interface{}, key, val interface{}) string {
	m[key] = val
	return ""
}

/* Type related */

// Type is the builtin to return the type of something as "int", "float", etc.
func (s *Syringe) Type(i interface{}) string {
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		return intString
	case reflect.Float32, reflect.Float64:
		return floatString
	case reflect.Slice, reflect.Array:
		return listString
	case reflect.Map:
		return mapString
	default:
		return unknownString
	}
}

// IsInt is the builtin that evaluates to `true` when its argument is an integer.
func (s *Syringe) IsInt(i interface{}) bool {
	return s.Type(i) == intString
}

// IsFloat is the builtin that evaluates to `true` when its argument is a floating point number.
func (s *Syringe) IsFloat(i interface{}) bool {
	return s.Type(i) == floatString
}

// IsNumber is the builtin that evaluates to `true` when its argument is a number (int or float).
func (s *Syringe) IsNumber(i interface{}) bool {
	return s.IsInt(i) || s.IsFloat(i)
}

// IsList is the builtin that evaluates to `true` when its argument is a list.
func (s *Syringe) IsList(i interface{}) bool {
	return s.Type(i) == listString
}

// IsMap is the builtin that evaluates to `true` when its argument is a map.
func (s *Syringe) IsMap(i interface{}) bool {
	return s.Type(i) == mapString
}

// Arithmetic.

// Add adds two numbers.
func (s *Syringe) Add(a, b interface{}) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Int() + bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Int() + int64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Int()) + bv.Float(), nil
		default:
			return nil, fmt.Errorf("add: unknown type for %q (%T)", bv, b)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int64(av.Uint()) + bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Uint() + bv.Uint(), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Uint()) + bv.Float(), nil
		default:
			return nil, fmt.Errorf("add: unknown type for %q (%T)", bv, b)
		}
	case reflect.Float32, reflect.Float64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Float() + float64(bv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Float() + float64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return av.Float() + bv.Float(), nil
		default:
			return nil, fmt.Errorf("add: unknown type for %q (%T)", bv, b)
		}
	default:
		return nil, fmt.Errorf("add: unknown type for %q (%T)", av, a)
	}
}

// Sub subtracts two numbers.
func (s *Syringe) Sub(a, b interface{}) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Int() - bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Int() - int64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Int()) - bv.Float(), nil
		default:
			return nil, fmt.Errorf("subtract: unknown type for %q (%T)", bv, b)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int64(av.Uint()) - bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Uint() - bv.Uint(), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Uint()) - bv.Float(), nil
		default:
			return nil, fmt.Errorf("subtract: unknown type for %q (%T)", bv, b)
		}
	case reflect.Float32, reflect.Float64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Float() - float64(bv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Float() - float64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return av.Float() - bv.Float(), nil
		default:
			return nil, fmt.Errorf("subtract: unknown type for %q (%T)", bv, b)
		}
	default:
		return nil, fmt.Errorf("subtract: unknown type for %q (%T)", av, a)
	}
}

// Mul multiplies two numbers.
func (s *Syringe) Mul(a, b interface{}) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Int() * bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Int() * int64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Int()) * bv.Float(), nil
		default:
			return nil, fmt.Errorf("multiply: unknown type for %q (%T)", bv, b)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int64(av.Uint()) * bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Uint() * bv.Uint(), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Uint()) * bv.Float(), nil
		default:
			return nil, fmt.Errorf("multiply: unknown type for %q (%T)", bv, b)
		}
	case reflect.Float32, reflect.Float64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Float() * float64(bv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Float() * float64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return av.Float() * bv.Float(), nil
		default:
			return nil, fmt.Errorf("multiply: unknown type for %q (%T)", bv, b)
		}
	default:
		return nil, fmt.Errorf("multiply: unknown type for %q (%T)", av, a)
	}
}

// Div divides two numbers.
func (s *Syringe) Div(a, b interface{}) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Int() / bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Int() / int64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Int()) / bv.Float(), nil
		default:
			return nil, fmt.Errorf("divide: unknown type for %q (%T)", bv, b)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int64(av.Uint()) / bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Uint() / bv.Uint(), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Uint()) / bv.Float(), nil
		default:
			return nil, fmt.Errorf("divide: unknown type for %q (%T)", bv, b)
		}
	case reflect.Float32, reflect.Float64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Float() / float64(bv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Float() / float64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return av.Float() / bv.Float(), nil
		default:
			return nil, fmt.Errorf("divide: unknown type for %q (%T)", bv, b)
		}
	default:
		return nil, fmt.Errorf("divide: unknown type for %q (%T)", av, a)
	}
}

// Loop is like a list of increasing ints.
func (s *Syringe) Loop(a, b int) <-chan int {
	ch := make(chan int)
	go func() {
		for i := a; i < b; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}
