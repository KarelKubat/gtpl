package syringe

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"
	"text/template"
)

const (
	// Type names
	intString     = "int"
	floatString   = "float"
	numberString  = "number"
	listString    = "list"
	mapString     = "map"
	unknownString = "unknown"

	// Name/version of this beast
	expanderName    = "gtpl"
	expanderVersion = "v1.0.5" // NOTE: Must match `gittag.txt`, TODO: make that automatic
)

// Logger is an interface that Syringe uses for "log" statements.
type Logger interface {
	Print(v ...interface{})
}

// Syringe is the receiver of the template functions injector.
type Syringe struct {
	logger   Logger
	logUsed  bool
	builtins []Builtin
}

// Opts are the options for New.
type Opts struct {
	Logger Logger // Used for "log" statements, defaults to https://pkg.go.dev/log
}

type Builtin struct {
	function interface{}
	Name     string
	Alias    string
	Usage    string
}

// New returns an initialized Syringe.
func New(o *Opts) *Syringe {
	s := &Syringe{
		logger: o.Logger,
	}
	if s.logger == nil {
		s.logger = log.Default()
	}
	s.builtins = []Builtin{
		// General
		{
			function: s.Expander,
			Name:     "Expander",
			Alias:    "expander",
			Usage:    "{{ expander }} - the name of this template expander",
		},
		{
			function: s.Version,
			Name:     "Version",
			Alias:    "version",
			Usage:    "{{ version }} - the version of this template expander",
		},
		{
			function: s.Log,
			Name:     "Log",
			Alias:    "log",
			Usage:    `{{ log "some" "info" }} - sends args to the log`,
		},
		{
			function: s.Die,
			Name:     "Die",
			Alias:    "die",
			Usage:    `{{ die "some" "info" }} - prints args, logs them if logging was used, stops`,
		},
		{
			function: s.Env,
			Name:     "Env",
			Alias:    "env",
			Usage:    `my homedir is {{ env "HOME" }} - returns environment setting`,
		},
		{
			function: s.Assert,
			Name:     "Assert",
			Alias:    "assert",
			Usage:    `asserts a condition and stops if not met: {{ assert (len $list) gt 0) "list is empty!" }}`,
		},

		// Strings
		{
			function: s.Strcat,
			Name:     "Strcat",
			Alias:    "strcat",
			Usage:    `{{ $all := strcat 12 " plus " 13 " is " 25 }}`,
		},
		{
			function: s.AddByte,
			Name:     "AddByte",
			Alias:    "addbyte",
			Usage:    `Add a '!': {{ $s := "Hello World" }} {{ $s = addbyte $s 33 }}`,
		},

		// Lists
		{
			function: s.List,
			Name:     "List",
			Alias:    "list",
			Usage:    `{{ $list := list "a" "b" "c" }} - creates a list`,
		},
		{
			function: s.HasElement,
			Name:     "HasElement",
			Alias:    "haselement",
		},
		{
			function: s.IndexOf,
			Name:     "IndexOf",
			Alias:    "indexof",
			Usage:    `'a' occurs at index {{ indexof $list "a" }} in the list`,
		},
		{
			function: s.AddElements,
			Name:     "AddElements",
			Alias:    "addelements",
			Usage:    `{{ $newlist := (addelements $list "d" "e") }} - creates a new list with added element(s)`,
		},

		// Maps
		{
			function: s.Map,
			Name:     "Map",
			Alias:    "map",
			Usage:    `{{ $map := map "cat" "meow" "dog" "woof" }} - creates a map`,
		},
		{
			function: s.HasKey,
			Name:     "HasKey",
			Alias:    "haskey",
		},
		{
			function: s.GetVal,
			Name:     "GetVal",
			Alias:    "getval",
			Usage:    `a cat says {{ get $map "cat" }} - gets a value from a map, "" if absent`,
		},
		{
			function: s.SetKeyVal,
			Name:     "SetKeyVal",
			Alias:    "setkeyval",
			Usage:    `{{ setkeyval $map "frog" "ribbit" }} - sets a key/value pair to a map`,
		},

		// Types
		{
			function: s.Type,
			Name:     "Type",
			Alias:    "type",
			Usage: `expands to "int", "float", "list" or "map"` + "\n" +
				`{{ $t := type $map }} {{ if $t ne "map" }} something is very wrong {{ end }}`,
		},
		{
			function: s.IsInt,
			Name:     "IsInt",
			Alias:    "isint",
			Usage:    `true when its argument is an integer`,
		},
		{
			function: s.IsFloat,
			Name:     "IsFloat",
			Alias:    "isfloat",
			Usage:    `true when its argument is a float`,
		},
		{
			function: s.IsNumber,
			Name:     "IsNumber",
			Alias:    "isnumber",
			Usage:    `true when its argument is an int or a float`,
		},
		{
			function: s.IsList,
			Name:     "IsList",
			Alias:    "islist",
			Usage:    `true when its argument is a list (or a slice)`,
		},
		{
			function: s.IsMap,
			Name:     "IsMap",
			Alias:    "ismap",
			Usage:    `true when its argument is a map`,
		},
		{
			function: s.Contains,
			Name:     "Contains",
			Alias:    "contains",
			Usage: `true when a map contains a key, a slice contains an element, or a string a substring` + "\n" +
				`{{ if contains $map "frog" }} .... {{ end }}`,
		},

		// Arithmetic / misc
		{
			function: s.Add,
			Name:     "Add",
			Alias:    "add",
			Usage:    `21 + 21 is {{ add 21 21 }}`,
		},
		{
			function: s.Sub,
			Name:     "Sub",
			Alias:    "sub",
			Usage:    `42 - 2 = {{ sub 42 2}}`,
		},
		{
			function: s.Mul,
			Name:     "Mul",
			Alias:    "mul",
			Usage:    `7 * 4 = {{ mul 7 4 }}`,
		},
		{
			function: s.Div,
			Name:     "Div",
			Alias:    "div",
			Usage:    `42 / 4 = {{ div 42 4 }}`,
		},
		{
			function: s.Loop,
			Name:     "Loop",
			Alias:    "loop",
			Usage:    `1 up to and including 10: {{ range $i := loop 1 11 }} {{ $i }} {{ end }}`,
		},
	}
	sort.Slice(s.builtins, func(i, j int) bool {
		return s.builtins[i].Name < s.builtins[j].Name
	})

	return s
}

// AliasesMap returns a `template.FuncMap` that can be passed to text/template so that shorthand
// builtins can be used.
func (s *Syringe) AliasesMap() template.FuncMap {
	fmap := template.FuncMap{}
	for _, b := range s.builtins {
		fmap[b.Alias] = b.function
	}
	return fmap
}

// Builtins returns the list of builtin functions.
func (s *Syringe) Builtins() []Builtin {
	return s.builtins
}

// Builtin functions.
// Remember to update the above info when adding/modifying!

/* General */

// Expander is the builtin returning the name of the expander program.
func (s *Syringe) Expander() string {
	return expanderName
}

// Version is the builtin returning the version of the expander program.
func (s *Syringe) Version() string {
	return expanderVersion
}

// Log is the builtin that logs information using the `log.Print` function.
func (s *Syringe) Log(args ...interface{}) string {
	s.logUsed = true
	parts := make([]string, len(args))
	for i, a := range args {
		parts[i] = fmt.Sprintf("%v", a)
	}
	s.logger.Print(fmt.Sprintf("%s: %s", expanderName, strings.Join(parts, " ")))
	return ""
}

// Die is the builtin that stops execution. If previous `Log` invocations occurred, then the
// the reason for stopping is logged, else, the reason is shown on `os.Stderr`.
func (s *Syringe) Die(args ...interface{}) (string, error) {
	msg := fmt.Sprint(args...)
	if s.logUsed {
		s.Log(msg)
	}
	return "", errors.New(msg)
}

// Env is the builtin that fetches the value of an environment variable.
func (s *Syringe) Env(str string) string {
	return os.Getenv(str)
}

// Assert is the builtin that ensures a condition.
func (s *Syringe) Assert(cond bool, args ...interface{}) (string, error) {
	if !cond {
		return "", fmt.Errorf("assert: %v", fmt.Sprint(args...))
	}
	return "", nil
}

/* String related */

// Strcat returns a string where all arguments are concatenated.
func (s *Syringe) Strcat(args ...interface{}) string {
	out := ""
	for _, a := range args {
		out += fmt.Sprintf("%v", a)
	}
	return out
}

// AddByte adds a byte to a string and returns the expanded string.
func (s *Syringe) AddByte(str string, val interface{}) (string, error) {
	var bt byte
	vk := reflect.ValueOf(val)
	switch vk.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bt = byte(vk.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bt = byte(vk.Uint())
	default:
		return "", fmt.Errorf("addbyte: unsupported type for %q (%T), only integers allowed", vk, val)
	}

	bytes := []byte(str)
	bytes = append(bytes, byte(bt))
	return string(bytes), nil
}

/* List related */

// List is the builtin that returns a list.
func (s *Syringe) List(args ...interface{}) []interface{} {
	return args
}

// HasElement is the builtin that checks whether a list contains an element.
// Deprecated: use Contains.
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
// Deprecated: use Contains.
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
func (s *Syringe) Type(i interface{}) (string, error) {
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		return intString, nil
	case reflect.Float32, reflect.Float64:
		return floatString, nil
	case reflect.Slice, reflect.Array:
		return listString, nil
	case reflect.Map:
		return mapString, nil
	default:
		return unknownString, errors.New("type: none of int, float, slice, array, map")
	}
}

// IsInt is the builtin that evaluates to `true` when its argument is an integer.
func (s *Syringe) IsInt(i interface{}) bool {
	t, _ := s.Type(i)
	return t == intString
}

// IsFloat is the builtin that evaluates to `true` when its argument is a floating point number.
func (s *Syringe) IsFloat(i interface{}) bool {
	t, _ := s.Type(i)
	return t == floatString
}

// IsNumber is the builtin that evaluates to `true` when its argument is a number (int or float).
func (s *Syringe) IsNumber(i interface{}) bool {
	return s.IsInt(i) || s.IsFloat(i)
}

// IsList is the builtin that evaluates to `true` when its argument is a list.
func (s *Syringe) IsList(i interface{}) bool {
	t, _ := s.Type(i)
	return t == listString
}

// IsMap is the builtin that evaluates to `true` when its argument is a map.
func (s *Syringe) IsMap(i interface{}) bool {
	t, _ := s.Type(i)
	return t == mapString
}

// Contains replaces HasElement or HasKey and offers substring matching.
func (s *Syringe) Contains(haystack interface{}, needle interface{}) (bool, error) {
	av := reflect.ValueOf(haystack)

	switch av.Kind() {
	case reflect.Map:
		m, ok := haystack.(map[interface{}]interface{})
		if !ok {
			return false, fmt.Errorf("contains: failed to convert %v to a map", haystack)
		}
		_, ok = m[needle]
		return ok, nil
	case reflect.Slice, reflect.Array:
		s, ok := haystack.([]interface{})
		if !ok {
			return false, fmt.Errorf("contains: failed to covert %v to a slice", haystack)
		}
		for _, e := range s {
			if e == needle {
				return true, nil
			}
		}
		return false, nil
	case reflect.String:
		s, ok := haystack.(string)
		if !ok {
			return false, fmt.Errorf("contains: failed to convert %v to a string", haystack)
		}
		bv := fmt.Sprintf("%v", needle)
		return strings.Contains(s, bv), nil
	default:
		return false, fmt.Errorf("contains: %v is neither a map, nor a slice, nor a string", haystack)
	}
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
