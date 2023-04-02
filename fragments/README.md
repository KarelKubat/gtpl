# gtpl: Generic (Golang) Template Expander

`gtpl` is a generic (Go-style) template parser that allows one to expand templates. `gtpl` prepares a number of handy functions that can be used in your templates, such as: maps management, lists, types.

I wrote `gtpl` because in some situations I have configuration files with a lot of boilerplate and just small variants for different cases. Typing out such configuration files is

- Toily, who wants to type almost the same thing 100 times.
- Error prone, what if you forget just one field in a sub-block.
- Hard to maintain, what if a sub-block suddenly needs another field. You need to go back and add that field to a hundred already existing sub-blocks.

Sounds familiar? `gtpl` to the rescue.
- It lets you define a template with a list of things to do, with maps (dicts) of settings
- It expands that template into the target configuration file.

To use `gtpl` you should know the standard built-in templating functions that the Go templating language provides, such as `index`, `len`. Read https://pkg.go.dev/text/template#hdr-Functions for more information. Additionally, the documentation below describes what `gtpl` adds.

Templates tend to generate a lot of noisy whitespace (unless they are very carefully crafted). To keep the output clean, `gtpl` understands a flag `-re` to remove empty lines. Also a prettyprinter that understands the output may come in handy.

## Usage

```shell
# Installation
go install https://github.com/KarelKubat/gtpl

# Quick overview of the built-ins
gtpl -b

# All flags
gtpl -h
```
