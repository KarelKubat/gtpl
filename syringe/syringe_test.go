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

func TestList(t *testing.T) {
	s := New(&Opts{Expander: expanderName, Version: versionID})
	l := s.List(0, 1, 2, 3, 4)
	for i, e := range l {
		if interface{}(reflect.ValueOf(e).Interface()) != i {
			t.Errorf("List(): value at index %v is %v, type %v", i, reflect.ValueOf(e), reflect.TypeOf(e))
		}
	}
}
