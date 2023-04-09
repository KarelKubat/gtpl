package syringe

import (
	"reflect"
	"strings"
	"testing"
)

const (
	longnamePrefix = ".Gtpl."
	expanderName   = "plugh" // You are in a maze of twisty little passages,
	versionID      = "xyzzy" // all different.
)

func TestDocStrings(t *testing.T) {
	s := New(&Opts{Expander: expanderName, Version: versionID})
	for _, b := range s.builtins {
		if !strings.HasPrefix(b.longname, longnamePrefix) {
			t.Errorf("docstrings: longname %q doesn't start with %q", b.longname, longnamePrefix)
		}
		if strings.ToLower(b.shortname) != b.shortname {
			t.Errorf("docstrings: shortname %q isn't all lowercase", b.shortname)
		}
		if b.usage == "" {
			t.Errorf("docstrings: function %q lacks usage", b.longname)
		}
	}
}

func TestExpanderAndVersion(t *testing.T) {
	s := New(&Opts{Expander: expanderName, Version: versionID})
	if s.Expander() != expanderName {
		t.Errorf("Expander() = %q, want %q", s.Expander(), expanderName)
	}
	if s.Version() != versionID {
		t.Errorf("Version() = %q, want %q", s.Version(), versionID)
	}
}

func TestListIsList(t *testing.T) {
	s := New(&Opts{Expander: expanderName, Version: versionID})
	l := s.List(0, 1, 2, 3, 4)
	for i, e := range l {
		if interface{}(reflect.ValueOf(e).Interface()) != i {
			t.Errorf("List(): value at index %v is %v, type %v", i, reflect.ValueOf(e), reflect.TypeOf(e))
		}
	}
	v := reflect.ValueOf(l)
	if v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
		t.Errorf("List(...) is a %v, want reflect.Array or reflect.Slice", v.Kind())
	}
}

func TestHasElement(t *testing.T) {
	s := New(&Opts{Expander: expanderName, Version: versionID})
	l := s.List(0, 1, 2)
	for i := 0; i <= 5; i++ {
		want := i <= 2
		if found := s.HasElement(l, interface{}(i)); found != want {
			t.Errorf("HasElement(%v, %v) = %v, want %v", l, i, found, want)
		}

	}
}

func TestIndexOf(t *testing.T) {
	s := New(&Opts{Expander: expanderName, Version: versionID})
	l := s.List(0, 1, 2, 3, 4, 5)
	for i := 0; i <= 5; i++ {
		if got := s.IndexOf(l, interface{}(i)); got != i {
			t.Errorf("IndexOf(%v, %v) = %v, want %v", l, i, got, i)
		}
	}
}

func TestAddElements(t *testing.T) {
	s := New(&Opts{Expander: expanderName, Version: versionID})
	l0 := s.List(0, 1, 2, 3)
	l1 := s.AddElements(l0, 4, 5)
	if len(l1) != 6 {
		t.Errorf("AddElements(%v, 4, 5) yields len %v, want 6", l0, len(l1))
	}
	for i := 0; i <= 5; i++ {
		if e := interface{}(reflect.ValueOf(l1[i]).Interface()); e != i {
			t.Errorf("after AddElments(%v, 4, 5): at index %v = %v, want %v", l0, i, e, i)
		}
	}
}

func TestMapIsMap(t *testing.T) {
	s := New(&Opts{Expander: expanderName, Version: versionID})
	m := s.Map(0, "zero", 1, "one", 2, "two", 3, "three")
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		t.Errorf("Map(...) is a %v, want reflect.Map", v.Kind())
	}
}

func TestHasKey(t *testing.T) {
	s := New(&Opts{Expander: expanderName, Version: versionID})
	m := s.Map(0, "zero", 1, "one", 2, "two", 3, "three")
	for _, test := range []struct {
		key        int
		wantHasKey bool
	}{
		{
			key:        0,
			wantHasKey: true,
		},
		{
			key:        4,
			wantHasKey: false,
		},
	} {
		if got := s.HasKey(m, test.key); got != test.wantHasKey {
			t.Errorf("HasKey(%v, %v) = %v, want %v", m, test.key, got, test.wantHasKey)
		}
	}
}

func TestGetAndSetVal(t *testing.T) {
	s := New(&Opts{Expander: expanderName, Version: versionID})
	m := s.Map(0, "zero", 1, "one", 2, "two", 3, "three")
	for _, test := range []struct {
		key     int
		wantVal string
	}{
		{
			key:     0,
			wantVal: "zero",
		},
		{
			key:     1,
			wantVal: "one",
		},
		{
			key:     4,
			wantVal: "",
		},
		{
			key:     42,
			wantVal: "",
		},
	} {
		if v := reflect.ValueOf(s.GetVal(m, test.key)); v.String() != test.wantVal {
			t.Errorf("Getval(%v, %v) = %v, want %v", m, test.key, v.String(), test.wantVal)
		}
	}

	fortyTwoString := "forty two"
	s.SetKeyVal(m, 42, fortyTwoString)
	if v := reflect.ValueOf(s.GetVal(m, 42)); v.String() != fortyTwoString {
		t.Errorf("Getval(%v, 42) = %v, want %v", m, v.String(), fortyTwoString)
	}
}

func TestIsKind(t *testing.T) {
	s := New(&Opts{Expander: expanderName, Version: versionID})
	for _, test := range []struct {
		v    interface{}
		want string
	}{
		{12, intString},
		{3.14, floatString},
		{[]string{"a", "b"}, listString},
		{[]int{1, 2, 3}, listString},
		{map[string]string{"a": "b"}, mapString},
	} {
		if got := s.Type(test.v); got != test.want {
			t.Errorf("Type(%v) = %q, want %q", test.v, got, test.want)
		}
	}

}
