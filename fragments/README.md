# gtpl: Generic (Golang) Template Expander

<!-- toc -->
<!-- /toc -->

`gtpl` is a generic (Go-style) template parser that allows one to expand templates. `gtpl` prepares a number of handy functions that can be used in your templates, such as: maps management, lists, types.

I wrote `gtpl` because in some situations I have configuration files with a lot of boilerplate and just small variants for different cases. Typing out such configuration files is:

- Toily, who wants to type almost the same thing 100 times.
- Error prone, what if you forget just one field in a sub-block.
- Hard to maintain, what if a sub-block suddenly needs another field. You need to go back and add that field to a hundred already existing sub-blocks.

Sounds familiar? `gtpl` to the rescue.

- It lets you define a template with a things to do, such as settings expressed in lists (arrrays), or maps (dicts).
- It expands such templates into a target configuration file.
- You can collect some "global" settings and apply them to different cases. `gtpl` will happily run:

    ```shell
    # assume common.tpl holds settings for both cases
    gtpl common.tpl onecase.tpl      > one.conf
    gtpl common.tpl anothercase.tpl  > another.conf
    ```

To use `gtpl` you should know the standard built-in templating functions that the Go templating language provides, such as `index`, `len`. 

- If you are completely new to Go's templates then there's a primer below.
- Go's templating system already provides a lot of [built-in functions](https://pkg.go.dev/text/template#hdr-Functions).
- `gtpl` adds functions to the standard list of built-ins, which is described below.

Templates tend to generate a lot of noisy whitespace (unless they are very carefully crafted). To keep the output clean, `gtpl` understands a flag `-re` to remove empty lines. Also a prettyprinter that understands the output may come in handy.

## Usage

```shell
# Installation: get the repo, then:
make install  # or just `go install gtpl.go`

# Quick overview of the built-ins
gtpl -b

# All flags
gtpl -h

# Run it
gtpl FILE1 FILE2 [FILE3...]
```

- Would you like to see all supported flags and the usage? Try `gtpl -h`.
- Would you like to see what builtins `gtpl` offers? Try `gtpl -b`.
- Do you dislike the action delimiters in template files, which default to `{{` and `}}`? Try `gtpl -left` and `gtpl -right`.
- See `gtpl -h` for a full overview.

`gtpl` also supports the filename `-` to indicate stdin; but to use it, you'll need the end-of-flags indicator `--`. Example:

```shell
gtpl -re -- file1 file2 - file3
```

This will:

- Suppress empty lines in the output (`-re` is a shorthand for `-remove-empty-lines`; the even shorter version `-r` can't be used as it is not distinguishable from `-right-delimiter`)
- Read `file1`
- Read `file2`
- Read whatever arrives on stdin
- Read `file3`
- Interpret everything that's read as a template.

In the case that template expansion fails, the error message will not clearly lead to the file and line number where the error occurs. In the above example the reported line number will point to somewhere in the bulk of of `file1`, `file2`, whatever was sent to `stdin`, and `file3`. To help with finding the offending error, you can re-run the command and supply `-li`:

```shell
# --list-template, or abbreviated -li, will list the template with
# line numbers before processing.
gtpl -re -li -- file1 file2 - file3
```

## Very Short Template Primer

You can skip this section if you know about Go's templating language. This section is meant for those who are completely new to it.

- **Delimiters:** A template holds text which is copied from the template to wherever it's supposed to go (*stdout* in the case of `gtpl`). Special actions are denoted by placing them inside delimiters, which are by default `{{` and `}}`, but you can choose your own if you want. You can use spaces or newlines following `{{` and preceding `}}`, that's just for readability and not required.

  ```C
  This is copied verbatim.
  {{ ... something special occurs here ... }}
  ```

- **Avoiding whitespace:** Since anything outside of `{{...}}` blocks is copied verbatim, templates can generate a lot of newlines and whitespace. You can use the opening delimiter `{{-` to suppress whitespace *before* the block, and the closing delimiter `-}}` to suppress whitespace *following* the block.

  ```C
  {{ ... something here ... -}}
  Part1
  {{- ... something here too ... -}}
  Part2
  ```

  This will ignore whitespace after the first block and around the second block. If the statements themselves don't output anything, this will lead to `Part1Part2`.

- **Comments:** Everything between `{{/*` and `*/}}` is skipped. There may be no space between `{{` and `/*` or between `*/` and `}}`.

- **Variables:** Variable names start with a `$` sign. New variables are created using the operator `:=`, when re-using variables, use `=`. The type of a variable can be anything (integer, floating point number, string, list, map) - the name is `$whatever` regardless of the type.

  ```C
  {{ $myvar := ... something ... }}
  {{ $myvar = ... another thing ... }}
  ```

- **No statement separators:** Each `{{ ... }}` block can contain one statement, there are no separators.

  ```C
  {{/* Good: */}}
  {{ $onevar := 12}}
  {{ $anothervar := 13 }}

  {{/* Bad: */}}
  {{ 
    $onevar := 12
    $anothervar := 13
  }}
  ```

- **Statement outcome:** A statement either sends something to the output, or can assign that "output" to a variable for later processing:

  ```C
  {{/* evaluated and if there's a result it is sent to the output */}}
  {{ ... something ... }}

  {{/* evaluated and assigned to $myvar, no output produced */}}
  {{ $myvar := ... something ... }}
  ```

- **Blocks:** Control structures have multiple `{{ ... }}` blocks. These are *if/end, if/else/end, range/end, define/end*.

  ```C
  {{ if ... some condition ... }}
    True-branch
  {{ else }}
    False-branch
  {{ end }}
  ```

- **Boolean expressions:** The template language supports the usual comparisons *ne/lt/le/gt/gt* and the booleans for *and/or/not*. However these are not operators in-between operands, but functions (so front of their operands):

  ```C
  {{ if lt $someval 10 }}
    That value is less than 10
  {{ end }}
  ```

- **Grouping using ( and )**: Using longer expressions will often involve grouping. A fragment like `and or lt $someval 10 gt $someval 20 ne $someval 0` is hard to read and also could be parsed in multiple ways. The templating language allows (insists on) `(` and `)` for grouping. In essence, wrap a fragment between `(` and `)` to force the evaluation into one result. Remember that you can add as many whitespace or newlines between `{{` and `}}` to improve readability.

  ```C
  {{/* $someval must be either less than 10 (but not zero) or greater than 20. */}}
  {{ if and (or (lt $someval 10)
                (gt $someval 20))
            (ne $someval 0) }} ... {{ else }} ... {{ end }}
  ```

- **Lists:** `len` is the length of a list, `index` returns an element, `range` iterates over all, `slice` creates a sub-list.

  ```C
  {{ $list := ... some list ... }}
  The length of the list is {{ len $list }}
  The first element is {{ index $list 0 }}
  The list contains:
  {{ range $el := $list }}
    {{ $el }}
  {{ end }}
  {{ $small := slice $list 0 3 }}
  A sub-list containing the first 3 elements is {{ $small }}
  ```

- **Maps:** Maps (or hashes, dicts) contain key/value pairs.

  ```C
  {{ $mymap := ... something ... }}
  {{ range $key,$val := $mymap }}
    The map holds value {{ $val }} at the key {{ $key }}
  {{ end }}
  ```

- **Templates within templates:** One can define a template within a template using `define`. Inside that template, the argument that the "caller" passes in is denoted by a dot. All usual functions can be used for that dot-value, such as `range` in the example below. The caller executes the template using the keyword `template`.

  ```C
  {{ define "mytemplate" }}
    {{ range $key,$val := . }}
      The map holds value {{ $val }} at the key {{ $key }}
    {{ end }}
  {{ end }}
  ...
  {{ $mymap := ... something ... }}
  {{ template "mytemplate" $mymap }}
  ```
  
- **Pipelines:** The last argument to a function can be either stated in the invocation, or it can be passed to the function using the pipeline symbol `|`. The following who examples are equivalent (`print` is a built-in function):

  ```C
  {{ print "Hello " "world"}}
  {{ "world" | print "Hello " }}
  ```
