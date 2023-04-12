package processor

import (
	"bytes"
	"strings"
	"testing"
)

func TestOverview(t *testing.T) {
	// Longnames are only stated in the overview when aliases are turned on.
	bare := New(&Opts{})
	if strings.Contains(bare.Overview(), "(longname") {
		t.Error("Overview() for a processor without aliases states longnames separately")
	}

	withAliases := New(&Opts{AllowAliases: true})
	if !strings.Contains(withAliases.Overview(), "(longname") {
		t.Errorf("Overview() for a processor with aliases fails to state longnames")
	}
}

func TestProcessStreams(t *testing.T) {
	// This is a bit of an integration test, going into package syringe as well.
	tpl := `
	{{ $list := list 1 2 3 4 5 }}
	{{ $elem := (index $list 1) }}
	Number: {{ $elem }}
	`
	rd := strings.NewReader(tpl)
	wr := &bytes.Buffer{}
	p := New(&Opts{AllowAliases: true})
	if err := p.ProcessStreams(rd, wr); err != nil {
		t.Fatalf("ProcessStreams(...) = %v, need nil error", err)
	}
	wantString := "Number: 2"
	if !strings.Contains(wr.String(), wantString) {
		t.Errorf("ProcessStreams(...): output is %q, doesn't contain %q", wr.String(), wantString)
	}
}
