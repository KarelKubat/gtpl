package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/KarelKubat/flagnames"
	"github.com/KarelKubat/gtpl/logger"
	"github.com/KarelKubat/gtpl/processor"
)

const (
	usageInfo = `
Welcome to gtpl, the Generic (Go-style) Template Expander.
Usage: gtpl [FLAGS] FILE [FILE...]

All files are scanned and executed as one template. File - (one hyphen) makes
gtpl read from stdin (you'll need a -- as flag terminator).

The below flags can be abbreviated
to a unique selector (-l can mean two things, -le is unique, -b is fine too).

`
)

var (
	logDest          = flag.String("log-output", "stderr", `log output: "stdout", "stderr" or a file to append`)
	allowAliases     = flag.Bool("allow-aliases", true, `when true, one can use "map" instead of ".Gtpl.Map" etc.`)
	leftDelimiter    = flag.String("left-delimiter", "", "opening delimiter in templates, {{ when unset")
	rightDelimiter   = flag.String("right-delimiter", "", "closing delimiter in templates, }} when unset")
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

	// Instantiate the logger.
	l, err := logger.New(*logDest)
	check(err)

	// Instantiate the processor.
	p := processor.New(&processor.Opts{
		AllowAliases:     *allowAliases,
		LeftDelimiter:    *leftDelimiter,
		RightDelimter:    *rightDelimiter,
		RemoveEmptyLines: *removeEmptyLines,
		Logger:           l,
	})

	// Show a short overview of builtins and stop, if requested.
	if *builtinsFlag {
		fmt.Println(p.Overview())
		os.Exit(0)
	}

	// At this point we want to process some files. We need at least 1 positional argument.
	if flag.NArg() < 1 {
		usage()
	}

	// All ready.
	check(p.ProcessFiles(flag.Args(), os.Stdout))

}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
