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

To use `gtpl` you should know the standard built-in templating functions that the Go templating language provides, such as `index`, `len`. Read https://pkg.go.dev/text/template#hdr-Functions for more information. Additionally, the documentation below describes what `gtpl` adds.

Templates tend to generate a lot of noisy whitespace (unless they are very carefully crafted). To keep the output clean, `gtpl` understands a flag `-re` to remove empty lines. Also a prettyprinter that understands the output may come in handy.

## Usage

```shell
# Installation: get the repo, then:
make install

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

In the case that template expansion fails, the error message will not clearly lead to the file and line number where the error occurs. In the above example the reported line number will point to somewhere in the total information of `file1`, `file2`, whatever was sent to `stdin`, and `file3`. To help with finding the offending error, you can run:

```shell
# --list-template, or abbreviated -li, will list the template with
# line numbers before processing.
gtpl -re -li -- file1 file2 - file3
```
