package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/KarelKubat/flagnames"
	"github.com/KarelKubat/gtpl/syringe"
)

const (
	expander  = "gtpl"
	version   = "0.0.1"
	usageInfo = `
Welcome to gtpl, the Generic (Go-style) Template Expander.
Usage: gtpl [FLAGS] FILE [FILE...]

All files are scanned and executed as one template. The below flags can be abbreviated
to a unique selector (-l can mean two things, -le is unique, -b is fine too).

`
)

type injected struct {
	Gtpl *syringe.Syringe
}

var (
	logDest          = flag.String("log-output", "stderr", `log output: "stdout", "stderr" or a file to append`)
	flatNameSpace    = flag.Bool("flat-namespace", true, `when true, one can use "map" instead of ".Gtpl.Map" etc.`)
	leftDelimiter    = flag.String("left-delimiter", "{{", "opening delimiter in templates")
	rightDelimiter   = flag.String("right-delimiter", "}}", "closing delimiter in templates")
	builtinsFlag     = flag.Bool("builtins", false, "when true, list built in functions and stop")
	removeEmptyLines = flag.Bool("remove-empty-lines", false, "when true, remove empty lines from the output")
)

func main() {
	// Parse the commandline.
	flagnames.Patch()
	usage := func() {
		fmt.Fprint(flag.CommandLine.Output(), usageInfo)
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output())
		os.Exit(1)
	}
	flag.Usage = usage
	flag.Parse()

	// Builtins to be injected into template processing.
	needle := syringe.New(&syringe.Opts{
		Expander: expander,
		Version:  version,
	})
	if *builtinsFlag {
		fmt.Println(needle.Overview())
		os.Exit(0)
	}

	// We need at least 1 file to process.
	if flag.NArg() < 1 {
		usage()
	}

	// Set the log output.
	switch *logDest {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	default:
		logStream, err := os.OpenFile(*logDest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		checkErr(err)
		log.SetOutput(logStream)
		defer logStream.Close()
	}

	// Wrap syringe data in a .Gtpl accessor for the templates.
	payload := &injected{
		Gtpl: needle,
	}

	// Prepare lowercase (flat space) names for injected functions, if requested.
	fmap := template.FuncMap{}
	if *flatNameSpace {
		fmap = needle.FlatNamespace()
	}
	// Read input files.
	readBuf := &bytes.Buffer{}
	for i := 0; i < flag.NArg(); i++ {
		dat, err := os.ReadFile(flag.Arg(i))
		checkErr(err)
		readBuf.Write(dat)
	}

	// Collect expansion in a buffer to avoid leading output in the case of errors.
	// Remove empty lines if so requested.
	tpl, err := template.New("gtpl").Funcs(fmap).Delims(*leftDelimiter, *rightDelimiter).Parse(string(readBuf.String()))
	checkErr(err)
	writeBuf := &bytes.Buffer{}
	checkErr(tpl.Execute(writeBuf, payload))
	if !*removeEmptyLines {
		_, err := os.Stdout.Write(writeBuf.Bytes())
		checkErr(err)
	} else {
		for _, line := range strings.Split(writeBuf.String(), "\n") {
			if strings.TrimSpace(line) != "" {
				fmt.Println(line)
			}
		}
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
