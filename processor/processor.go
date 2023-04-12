// Package processor wraps a Syringe for easier integration.
package processor

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/KarelKubat/gtpl/syringe"
)

const (
	gtplNamePrefix = ".Gtpl" // Prefix for full function names, must match the membername in the `injected` struct
)

// Opts control how the processor works.
type Opts struct {
	AllowAliases     bool           // When true, allow short function names ("map") as aliases (for "".Gtpl.Map")
	LeftDelimiter    string         // When "", defaults to "{{"
	RightDelimter    string         // When "", defaults to "}}"
	RemoveEmptyLines bool           // When true, remove empty lines from the output
	Logger           syringe.Logger // When nil, defaults to https://pkg.go.dev/log
}

// Processor is the receiver.
type Processor struct {
	o          *Opts            // Input options
	needle     *syringe.Syringe // Actual template processor
	fmap       template.FuncMap // Non-empty when aliases are allowed
	leftDelim  string           // start-of-instruction
	rightDelim string           // end-of-instruction
}

func New(o *Opts) *Processor {
	p := &Processor{
		o: o,
		needle: syringe.New(&syringe.Opts{
			Logger: o.Logger,
		}),
		fmap:       template.FuncMap{},
		leftDelim:  o.LeftDelimiter,
		rightDelim: o.RightDelimter,
	}

	// Patch up non-standard options
	if o.AllowAliases {
		p.fmap = p.needle.AliasesMap()
	}
	return p
}

type injected struct {
	Gtpl *syringe.Syringe
}

// ProcessStreams reads the template to process from an io.Reader and runs it. The output goes to an io.Writer.
func (p *Processor) ProcessStreams(r io.Reader, w io.Writer) error {
	// Collect everything from the input reader.
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	if err != nil {
		return err
	}

	// Run the template.
	tpl, err := template.New("gtpl").Funcs(p.fmap).Delims(p.leftDelim, p.rightDelim).Parse(buf.String())
	if err != nil {
		return err
	}

	// If we don't need to postprocess the output for empty lines, then the template can be executed and the output goes
	// directly to the requrested writer.
	if !p.o.RemoveEmptyLines {
		return tpl.Execute(w, &injected{
			Gtpl: p.needle,
		})
	}

	// To remove empty lines, we need to collect the execution output and re-examine it.
	var wrbuf bytes.Buffer
	err = tpl.Execute(&wrbuf, &injected{
		Gtpl: p.needle,
	})
	var trimmed bytes.Buffer
	for _, line := range strings.Split(wrbuf.String(), "\n") {
		if strings.TrimSpace(line) != "" {
			trimmed.WriteString(line + "\n")
		}
	}
	trimmed.WriteTo(w)
	return err
}

// Builtins returns the "usage" information of the builtin functions, just as syringe.Overview does.
// This is just a pass-through.
func (p *Processor) Overview() string {
	out := ""
	for _, b := range p.needle.Builtins() {
		if p.o.AllowAliases {
			out += fmt.Sprintf("%v (longname: %v.%v)\n", b.Alias, gtplNamePrefix, b.Name)
		} else {
			out += fmt.Sprintf("%v.%v\n", gtplNamePrefix, b.Name)
		}
		for _, line := range strings.Split(b.Usage, "\n") {
			out += "  " + line + "\n"
		}
		out += "\n"
	}
	return out
}

// ProcessFiles reads templates from files. The output goes to an io.Writer.
func (p *Processor) ProcessFiles(files []string, w io.Writer) error {
	var total bytes.Buffer
	for _, f := range files {
		var b []byte
		var err error
		if f == "-" {
			var stdin bytes.Buffer
			_, err = stdin.ReadFrom(os.Stdin)
			if err != nil {
				return err
			}
			b = stdin.Bytes()
		} else {
			b, err = os.ReadFile(f)
			if err != nil {
				return err
			}
		}
		_, err = total.Write(b)
		if err != nil {
			return err
		}
	}
	return p.ProcessStreams(&total, w)
}
