## Expanding `gtpl` or embedding it in your own Go programs

### Package `processor`

If you want to embed the template processor in your own Go programs, then the easiest way is to import `github.com/KarelKubat/gtpl/processor` and to use that. An example is in the top level main program `gtpl.go`.

The processor is instantiated using options that define whether to remove empty lines from the output, whether to aliases (`map` as alias for `.Gtpl.Map` etc.). Then the processor can be started to expand the templates from a reader stream or from files. The output is sent to a writer stream for the caller to process. The minimal example is:

```go
// No special options:
// - No function aliases (builtins are `.Gtpl.Map`, no alias `map` etc.
// - Left delimiter is Go's default {{, right delimter is }}
// - Empty lines in the output are not removed
// - .Gtpl.Log functions invoke the standard Go logger
p := processor.New(&processor.Opts{
    // Nothing to see here
})

// Template(s) are expected on stdin, output goes to stdout
err := p.ProcessStreams(os.Stdin, os.Stdout)
```

The logger that `.Gtpl.Log` invokes (the alias `log` exists when aliases are enabled) must satisfy the interface `syringe.Logger`, which means that it must have a member function `Print()`. A customized logger can be plugged in as follows:

- You can pass any receiver to something that implements `Print()`
- You can instantiate the default logger, using `log.Default()` and customize it, then pass that
- A very simple version is in `github.com/KarelKubat/gtpl/logger`. This package uses the standard Go logger but sends output to stderr, stdout or to a file. The top-level main program `gtpl.go` uses that.

### Package `syringe`

A more low-level library is `github.com/KarelKubat/gtpl/syringe`. This package actually implements the functions such as `list` or `map` and injects them into the template processor. Supplying the template and expanding it (using the standard `text/template` package) is left to the caller.

**Do not change the fingerprint of builtins**, that breaks backwards compatibility. If needed, implement a new functions that does what you need. Adding checks to an existing function, fixing bugs or the like is of course okay.

To expand the list of builtins, please proceed as follows:

- Implement the function by adding it to the right section (general, list-related etc.).
- Expand the list of builtins which is constructed in `New()`. This list maps functions `.Gtpl.Whatever` to their short name and provides a very short description.
- Update the version string at the top of the file.
- Send me a pull request :)