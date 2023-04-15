# gtpl: Generic (Golang) Template Expander

<!-- toc -->
- [Usage](#usage)
- [Short Examples](#short-examples)
  - [Example: examples/00-general.tpl](#example-examples00-generaltpl)
  - [Example: examples/01-types.tpl](#example-examples01-typestpl)
  - [Example: examples/02-arith.tpl](#example-examples02-arithtpl)
  - [Example: examples/03-list.tpl](#example-examples03-listtpl)
  - [Example: examples/04-mambo.tpl](#example-examples04-mambotpl)
  - [Example: examples/05-maps.tpl](#example-examples05-mapstpl)
  - [Example: examples/06-fibo.tpl](#example-examples06-fibotpl)
- [Example: Generated SSH Configuration and a Ping Check](#example-generated-ssh-configuration-and-a-ping-check)
  - [Source of Truth](#source-of-truth)
  - [Generation of an SSH Configuration](#generation-of-an-ssh-configuration)
  - [Pinging the Configured Hosts](#pinging-the-configured-hosts)
- [List of Built in Functions](#list-of-built-in-functions)
- [Expanding <code>gtpl</code> or embedding it in your own Go programs](#expanding-gtpl-or-embedding-it-in-your-own-go-programs)
  - [Package <code>processor</code>](#package-processor)
  - [Package <code>syringe</code>](#package-syringe)
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

## Short Examples

See also `examples/*tpl`.

### Example: examples/00-general.tpl

```
{{/*
  Demo of:
    expander - the name of the expander program
    version  - its version
    log      - logging information
    assert   - ensuring that a condition is met
    die      - how to abort a run
    env      - returns environment setting
  Also:
    comments - you can't have spaces between the template and comment delimiters
*/}}

{{ log "This generates " "one" " log statement" }}
This template is processed by {{ expander }} version {{ version }}

{{ $answer := 42 }}
{{/* abort if it isn't 42 */}}
{{ assert (eq $answer 42) "answer is " $answer " but i need 42" }}

{{/* Uncomment to cause an error.
  {{ die "stop execution with an error" }}
*/}}

My homedir is {{ env "HOME" }}
```

**Output** (empty lines removed):

```
2023/04/15 12:47:42 gtpl: This generates one log statement
This template is processed by gtpl version 0.0.3
My homedir is /Users/karelk
```
### Example: examples/01-types.tpl

```
{{/*
    Demo of:
        list, map - constructions, see next sections
        type      - type of its argument: "int", "float", "number", "list", "map"
        isint     - true for integers
        isfloat   - true for floating point numbers
        isnumber  - true for floats or ints
        islist    - true for lists
        ismap     - true for maps
    Also standard built ins:
        range    - how to loop over a list
*/}}

{{/* 
  A few type variants: int, float, list, map. The variants themselves are wrapped
  in a list. The map is created using sequential arguments, so 
    key1 val1 key2 val2
  etc.. There is no syntax to express 
    key1: val1, key2: val2 // or key1 -> val1, key2 -> val2 whatever
*/}}
{{ $variants := list 
    42 
    3.14
    (list "a" "b" "c")
    (map "firstname" "Karel" 
          "lastname" "Kubat")
}}

{{ range $v := $variants }} 
  {{ $v }} is a(n) {{ type $v }} 
{{ end }}

42 is {{ if not (isint 42)}} not {{ end }} an int
42 is {{ if not (isfloat 42)}} not {{ end }} a float
42 is {{ if not (isnumber 42)}} not {{ end }} a number
42 is {{ if not (islist 42)}} not {{ end }} a list
42 is {{ if not (ismap 42)}} not {{ end }} a map
```

**Output** (empty lines removed):

```
  42 is a(n) int 
  3.14 is a(n) float 
  [a b c] is a(n) list 
  map[firstname:Karel lastname:Kubat] is a(n) map 
42 is  an int
42 is  not  a float
42 is  a number
42 is  not  a list
42 is  not  a map
```
### Example: examples/02-arith.tpl

```
{{/*
  Demo of:
    add - adds two numbers
    sub - subtracts the second number from the first one
    mul - multiplies two numbers
    div - divides the first number by the second one
*/}}

12 + 3 = {{ add 12 3 }}
12 - 3 = {{ sub 12 3 }}
12 * 3 = {{ mul 12 3 }}
12 / 3 = {{ div 12 3 }}
```

**Output** (empty lines removed):

```
12 + 3 = 15
12 - 3 = 9
12 * 3 = 36
12 / 3 = 4
```
### Example: examples/03-list.tpl

```
{{/*
    Demo of:
      list        - creates a list
      addelements - creates a new list by adding elements
      haselement  - checks whether an element is in the list
      indexof     - returns the index of an element
    Also standard built ins:
      len         - returns the length of a list
      slice       - how to get a partial list from a list          
      index       - how to get one element from a list
*/}}

{{ $list := list "one" "two" "three" }}
The list so far: {{ $list }}
It has {{ len $list }} elements.
The first two elements are: {{ slice $list 0 2 }}
The second element is {{ index $list 1}}
Element "three" occurs at index {{ indexof $list "three" }}

Let's add "four" and "five".
{{ $list = addelements $list "four" "five" }}
I've $got {{ range $sense := $list }}{{ $sense }} {{ end }}senses working overtime.

{{ if (haselement $list "five") }}
  "five" is in the list
{{ else }}
  "five" is not in the list (this would be an error)
{{ end }}
```

**Output** (empty lines removed):

```
The list so far: [one two three]
It has 3 elements.
The first two elements are: [one two]
The second element is two
Element "three" occurs at index 2
Let's add "four" and "five".
I've $got one two three four five senses working overtime.
  "five" is in the list
```
### Example: examples/04-mambo.tpl

```
{{/*
    Demo of:
        map   - creates a map
    Also standard built ins:
        range - ranging key,value pairs over a map
*/}}

{{ $lyrics := map
          "Monica"    "in my life"
          "Erica"     "by my side"
          "Rita"      "is all I need"
          "Tina"      "is what I see"
          "Sandra"    "in the sun"
          "Mary"      "all night long"
          "Jessica"   "here I am"
          "you"       "makes me your man" }}
{{ range $name, $what := $lyrics }}
A little bit of {{ $name }} {{ $what }}
{{ end }}
```

**Output** (empty lines removed):

```
A little bit of Erica by my side
A little bit of Jessica here I am
A little bit of Mary all night long
A little bit of Monica in my life
A little bit of Rita is all I need
A little bit of Sandra in the sun
A little bit of Tina is what I see
A little bit of you makes me your man
```
### Example: examples/05-maps.tpl

```

{{/*
    Demo of:
      map       - creates a map
      haskey    - checks whether a key is in a map
      setkeyval - adds a key/value pair to a map
      assert    - ensures that a condition is true
    Also standard built ins:
      define    - defining a template
      template  - calling a template
      range     - iterating key/value pairs over a map
*/}}

{{ 
    $parties := map
        "Alice"    (map "role"       "sender"
                        "isAttacker" false)
        "Bob"      (map "role"       "recipient"
                        "isAttacker" false)
        "Mallory"  (map "role"       "man in the middle"
                         "isAttacker" true)
}}

{{ define "showParty" }}
  {{/* Ensure that the argument is a map with the keys role and isAttacker: */}}
  {{ assert (ismap .) "showParty: arg is not a map" }}
  {{ assert (haskey . "role") "showParty: arg map doesn't have role key" }}
  {{ assert (haskey . "isAttacker") "showParty: arg map doesn't have isAttacker key" }}

  {{/* 
    Since it's a valid and asserted map, we can `getval` from it and we'l get
    real values. Otherwise `getval` would evaluate to "" for non-existing keys.
  */}}
  Role: {{getval . "role"}}
  Attacker: {{getval . "isAttacker"}}
{{ end }}

{{ range $name, $data := $parties }}
Name: {{ $name }}
{{ template "showParty" $data }}
{{ end }}

Alice {{ if haskey $parties "Alice" }} occurs {{ else }} doesn't occur {{ end }} in the map.

{{ if not (haskey $parties "Eve") }}
Eve is not listed as a party yet. Let's add her.
{{ setkeyval $parties "Eve"
    (map "role" "another attacker"
                "isAttacker" true) }}
{{ assert (haskey $parties "Eve") "Eve must now be known as a party." }}
{{ end }}
Name: Eve
{{ template "showParty" (getval $parties "Eve") }}
```

**Output** (empty lines removed):

```
Name: Alice
  Role: sender
  Attacker: false
Name: Bob
  Role: recipient
  Attacker: false
Name: Mallory
  Role: man in the middle
  Attacker: true
Alice  occurs  in the map.
Eve is not listed as a party yet. Let's add her.
Name: Eve
  Role: another attacker
  Attacker: true
```
### Example: examples/06-fibo.tpl

```
{{/*
  Demo of:
    loop - shorthand for a list of consecutive ints
    add  - adds two numbers
  Also standard built ins:
    range
*/}}

Fibonacci series

{{ $a := 1 }}
{{ $b := 2 }}

{{/* `loop 0 10` is a shorthand for `list 0 1 2 3 4 5 6 7 8 9` */}}
{{/* That means "up to 10", not "and including". */}}
{{/* Or: `loop 0 100` means: 100 times. */}} 
{{ range $i := loop 1 11 }}
  Number {{ $i }}: {{ $a }}
  {{ $tmp := $a }}
  {{ $a = $b }}
  {{ $b = add $tmp $b }}
{{ end }}
```

**Output** (empty lines removed):

```
Fibonacci series
  Number 1: 1
  Number 2: 2
  Number 3: 3
  Number 4: 5
  Number 5: 8
  Number 6: 13
  Number 7: 21
  Number 8: 34
  Number 9: 55
  Number 10: 89
```

## Example: Generated SSH Configuration and a Ping Check

This is a longer example, where a "source of truth" file is used

- to generate an SSH configuration, and 
- to generate `ping` commands to see if the hosts are up and reachable.

### Source of Truth

Here is a file that defines a bunch of (imaginary) hosts. It is the substrate for the generation of the SSH configuration `~/.ssh/config` and for a ping test.

Each entry in the list is a map that defines for a host:

- A description (for comment purposes)
- A hostname or IP address (required)
- A short hame. When given, this name is taken as an alias that `ssh` understands.
- A user name. When given, `ssh` connections are set up as this user, when not, the logged-in user will be taken.
- A port. When given, the `ssh` connections run over the stated port.
- Whether the host supports X11. When yes, X11 connections are forwarded over the `ssh` connection.
- An identity file. When given, that specific identity file is used for authentication.

You can also find this as `examples/hosts/hosts`.

```
{{{/* Configuration of ssh-able hosts */}}

{{
    $hosts := list
        (map "description"   "My computer at work"
             "hostname"      "ws12345.example.com"
             "shortname"     "ws"
             "user"          "user12345"
             "idfile"        "/home/user/.ssh/specific_id_rsa_file"
             "hasX11"        true)

        (map "description"   "Bastion, emergency access from outside"
             "hostname"      "123.456.789.012"
             "shortname"     "bastion"
             "user"          "emergency"
             "port"          2222)

        (map "description"   "Office DHCP/Router"
             "hostname"      "192.168.1.1"
             "user"          "pi")

        (map "description"   "Raspberry Pi DNS/Blackhole"
             "hostname"      "192.168.1.10"
             "shortname"     "pi"
             "user"          "pi")
}}
```

### Generation of an SSH Configuration

In order to generate a configuration that `ssh` understands, the following template is used. It iterates over the above list, and for each map it picks up keys and values. Additionally, a few defaults are injected into each SSH host configuration: to use compression, to multiplex connections once logged in, and so on.

You can also find this as `examples/hosts/ssh-config`.

```
{{/* Expansion of hosts */}}

{{/* Use as follows:
    gtpl -re hosts ssh-config > ~/.ssh/config
*/}}

{{ range $h := $hosts }}
#
# {{ getval $h "description" }}
# ------------------------------------------------------------------------
Host {{ if haskey $h "shortname"}}{{ getval $h "shortname"}}{{ else }}{{ getval $h "hostname" }}{{ end }}
    Hostname {{ getval $h "hostname" }}
    {{ if haskey $h "user" }}User {{ getval $h "user" }}{{ end }}
    Compression yes
    ControlMaster auto
    ControlPath ~/.ssh/ctrl-%C
    ControlPersist yes
    ServerAliveInterval 300
{{ if haskey $h "idfile" }}
    IdentityFile {{ getval $h "idfile" }}
    IdentitiesOnly yes
{{ end }}
{{ if haskey $h "hasX11" }}
    ForwardX11Trusted yes
{{ end }}
{{ if haskey $h "port" }}
    Port {{ getval $h "port" }}
{{ end }}
{{ end }}
```

Sample output:

```
#
# My computer at work
# ------------------------------------------------------------------------
Host ws
    Hostname ws12345.example.com
    User user12345
    Compression yes
    ControlMaster auto
    ControlPath ~/.ssh/ctrl-%C
    ControlPersist yes
    ServerAliveInterval 300
    IdentityFile /home/user/.ssh/specific_id_rsa_file
    IdentitiesOnly yes
    ForwardX11Trusted yes
#
# Bastion, emergency access from outside
# ------------------------------------------------------------------------
Host bastion
    Hostname 123.456.789.012
    User emergency
    Compression yes
    ControlMaster auto
    ControlPath ~/.ssh/ctrl-%C
    ControlPersist yes
    ServerAliveInterval 300
    Port 2222
#
# Office DHCP/Router
# ------------------------------------------------------------------------
Host 192.168.1.1
    Hostname 192.168.1.1
    User pi
    Compression yes
    ControlMaster auto
    ControlPath ~/.ssh/ctrl-%C
    ControlPersist yes
    ServerAliveInterval 300
#
# Raspberry Pi DNS/Blackhole
# ------------------------------------------------------------------------
Host pi
    Hostname 192.168.1.10
    User pi
    Compression yes
    ControlMaster auto
    ControlPath ~/.ssh/ctrl-%C
    ControlPersist yes
    ServerAliveInterval 300
```

### Pinging the Configured Hosts

The following template picks up hosts information and generates `ping` comands to see if the hosts are accessible.

You can also find this as `examples/hosts/ping-test`.

```
{{/* Ping test for hosts. */}}

{{/* Use as follows:
    gtpl -re hosts ping-test | sh
*/}}

{{ range $h := $hosts }}
ping -c3 -t3 {{ getval $h "hostname" }} >/dev/null 2>&1 || echo {{ getval $h "hostname" }} is unreachable!
{{ end }}
```

Example output (which can be fed to the `/bin/sh`):

```
ping -c3 -t3 ws12345.example.com >/dev/null 2>&1 || echo ws12345.example.com is unreachable!
ping -c3 -t3 123.456.789.012 >/dev/null 2>&1 || echo 123.456.789.012 is unreachable!
ping -c3 -t3 192.168.1.1 >/dev/null 2>&1 || echo 192.168.1.1 is unreachable!
ping -c3 -t3 192.168.1.10 >/dev/null 2>&1 || echo 192.168.1.10 is unreachable!
```

## List of Built in Functions

The list can be generated using `gtpl -b`.
The lowercase aliases (e.g., `add` for `.Gtpl.Add`) are **not** available
when the flag `--allow-aliases=false` is given. 

```
add (longname: .Gtpl.Add)
  21 + 21 is {{ add (21 21) }}

addelements (longname: .Gtpl.AddElements)
  {{ $newlist := (addelements $list "d" "e") }} - creates a new list with added element

assert (longname: .Gtpl.Assert)
  asserts a condition and stops if not met: {{ assert (len $list) gt 0) "list is empty!" }}

die (longname: .Gtpl.Die)
  {{ die "some" "info" }} - prints args, logs them if logging was used, stops

div (longname: .Gtpl.Div)
  42 / 4 = {{ div 42 4 }}

env (longname: .Gtpl.Env)
  my homedir is {{ env "HOME" }} - returns environment setting

expander (longname: .Gtpl.Expander)
  {{ expander }} - the name of this template expander

getval (longname: .Gtpl.GetVal)
  a cat says {{ get $map "cat" }} - gets a value from a map, "" if absent

haselement (longname: .Gtpl.HasElement)
  {{ if (haselement $list "a") }} 'a' occurs in the list {{ end }}

haskey (longname: .Gtpl.HasKey)
  {{ if haskey $map "cat" }} yes {{ else }} no {{ end }} - tests whether a key is in a map

indexof (longname: .Gtpl.IndexOf)
  'a' occurs at index {{ indexof $list "a" }} in the list

isfloat (longname: .Gtpl.IsFloat)
  true when its argument is a float

isint (longname: .Gtpl.IsInt)
  true when its argument is an integer

islist (longname: .Gtpl.IsList)
  true when its argument is a list (or a slice)

ismap (longname: .Gtpl.IsMap)
  true when its argument is a map

isnumber (longname: .Gtpl.IsNumber)
  true when its argument is an int or a float

list (longname: .Gtpl.List)
  {{ $list := list "a" "b" "c" }} - creates a list

log (longname: .Gtpl.Log)
  {{ log "some" "info" }} - sends args to the log

loop (longname: .Gtpl.Loop)
  1 up to and including 10: {{ range $i := loop 1 11 }} {{ $i }} {{ end }}

map (longname: .Gtpl.Map)
  {{ $map := map "cat" "meow" "dog" "woof" }} - creates a map

mul (longname: .Gtpl.Mul)
  7 * 4 = {{ mul 7 4 }}

setkeyval (longname: .Gtpl.SetKeyVal)
  {{ set $map "frog" "ribbit" }} - adds a key/value pair to a map

sub (longname: .Gtpl.Sub)
  42 - 2 = {{ sub 42 2}}

type (longname: .Gtpl.Type)
  expands to "int", "float", "list" or "map"
  {{ $t := type $map }} {{ if $t ne "map" }} something is very wrong {{ end }}

version (longname: .Gtpl.Version)
  {{ version }} - the version of this template expander


```

## Expanding `gtpl` or embedding it in your own Go programs

### Package `processor`

If you want to embed the template processor in your own Go programs, then the easiest way is to import `github.com/KarelKubat/gtpl/processor` and to use that. An example is in the top level main program `gtpl.go`.

The processor is instantiated using options that define whether to remove empty lines from the output, whether to aliases (`map` as alias for `.Gtpl.Map` etc.). Then the processor can be started to expand the templates from a reader stream or from files. The output is sent to a writer stream for the caller to process. The minimal example is:

```go
// No special options:
// - No function aliases (builtins are `.Gtpl.Map`, no alias `map` etc.
// - Left delimiter is Go's default `{{`, right delimter is `}}`
// - Empty lines in the output are not removed
// - `.Gtpl.Log` invokes the standard Go logger
p := processor.New(&processor.Opts{
    // Nothing to see here
})

// Template(s) are expected on stdin, output goes to stdout
err := p.ProcessStreams(os.Stdin, os.Stdout)
```

The logger that `.Gtpl.Log` invokes (the alias `log` exists when aliases are enabled) must satisfy the interface `syringe.Logger`, which means that it must have a member function `Print()`. A customized logger can be plugged in as follows:

- You can pass a receiver to anything that implements `Print()`
- You can instantiate Go's default logger using `log.Default()`, customize it, then pass that
- A very simple version is in `github.com/KarelKubat/gtpl/logger`. This package uses the standard Go logger but sends output to stderr, stdout or to a file. The top-level main program `gtpl.go` uses that.

### Package `syringe`

A more low-level library is `github.com/KarelKubat/gtpl/syringe`. This package actually implements the functions such as `list` or `map` and injects them into the template processor. Supplying the template and expanding it (using the standard `text/template` package) is left to the caller.

**Do not change the fingerprint of builtins**, that breaks backwards compatibility. If needed, implement a new functions that does what you need. Adding checks to an existing function, fixing bugs or the like is of course okay.

To expand the list of builtins or to fix a bug, please proceed as follows:

- Implement new functions by adding them to the correct section (general, list-related etc.).
- If you add a function, then also state it in the list of builtins which is constructed in `New()`. This list maps function names such as `SomeLongName` to their aliases and provides very short descriptions.
- Update the version string at the top of the file.
- Send me a pull request :)

