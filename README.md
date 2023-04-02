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
## Examples

See also `examples/*tpl`.

### Example: examples/00-general.tpl

```
This template is processed by {{ expander }} version {{ version }}

{{ log "This generates " "one" " log statement" }}

{{ $answer := 42 }}
{{/* abort if it isn't 42 */}}
{{ assert (eq $answer 42) "answer is " $answer " but i need 42" }}

{{/* Uncomment to cause an error.
  {{ die "stop execution with an error" }}
*/}}
```

**Output** (empty lines removed):

```
2023/04/02 14:02:59 gtpl: This generates one log statement
This template is processed by gtpl version 0.0.1
```
### Example: examples/01-types.tpl

```
{{/*
    Demo of:
        type     - type of its argument: "int", "float", "number", "list", "map"
        isint    - true for integers
        isfloat  - true for floating point numbers
        isnumber - true for floats or ints
        islist   - true for lists
        ismap    - true for maps
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
### Example: examples/03-list.tpl

```
{{/*
    Demo of:
      list        - creates a list
      addelements - creates a new list by adding elements
      haselement  - checks whether an element is in the list
    Also standard built ins:
      slice
      index
*/}}

{{- $list := list "one" "two" "three" -}}
The list so far: {{ $list }}
The first two elements are: {{ slice $list 0 2 }}
The second element is {{ index $list 1}}

Let's add "four" and "five".
{{- $list = addelements $list "four" "five" }}
I've got {{ range $sense := $list }}{{ $sense }} {{ end }}senses working overtime.

{{ if (haselement $list "five") }}
  "five" is in the list
{{ else }}
  "five" is not in the list (this would be an error)
{{ end }}
```

**Output** (empty lines removed):

```
The list so far: [one two three]
The first two elements are: [one two]
The second element is two
Let's add "four" and "five".
I've got one two three four five senses working overtime.
  "five" is in the list
```
### Example: examples/04-mambo.tpl

```
{{/*
    Demo of:
        map - creates a map
    Also standard built ins:
        range
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
      define
      template
      range
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

  {{/* Since it's a map, we can `getval` from it. */}}
  Role: {{getval . "role"}}
  Attacker: {{getval . "isAttacker"}}
{{ end }}

{{ range $name, $data := $parties }}
Name: {{ $name }}
{{ template "showParty" $data }}
{{ end }}

Alice {{ if haskey $parties "Alice" }} occurs {{ else }} doesn't occur {{ end }} in the map as a party.

{{ if not (haskey $parties "Eve") }}
Eve is not listed as a party yet. Let's add her.
{{ setkeyval $parties "Eve"
    (map "role" "another attacker"
                "isAttacker" true) }}
{{ assert (haskey $parties "Eve") "Eve must now be known as a party." }}
{{ end }}
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
Alice  occurs  in the map as a party.
Eve is not listed as a party yet. Let's add her.
```
## List of Built in Functions

The list can be generated using `gtpl -b`.

```
expander (long name: .Gtpl.Expander)
  {{expander}} - the name of this template expander

version (long name: .Gtpl.Version)
  {{version}} - the version of this template expander

log (long name: .Gtpl.Log)
  {{log "some" "info"}} - sends args to the log

die (long name: .Gtpl.Die)
  {{die "some" "info"}} - prints args, logs them if logging was used, stops

assert (long name: .Gtpl.Assert)
  asserts a condition and stops if not met: {{$len := len $list}}{{assert ($len gt 0) "list is empty!"}}

list (long name: .Gtpl.List)
  {{$list := list "a" "b" "c"}} - creates a list

haselement (long name: .Gtpl.HasElement)
  {{if (haselement $list "a")}} 'a' occurs in the list {{end}}

addelements (long name: .Gtpl.AddElements)
  {{$newlist := (addelements $list "d" "e")}} - creates a new list with added element

map (long name: .Gtpl.Map)
  {{$map := map "cat" "meow" "dog" "woof"}} - creates a map

haskey (long name: .Gtpl.HasKey)
  {{if haskey $map "cat"}} yes {{else}} no {{end}} - tests whether a key is in a map

getval (long name: .Gtpl.GetVal)
  a cat says {{get $map "cat"}} - gets a value from a map, "" if absent

setkeyval (long name: .Gtpl.SetKeyVal)
  {{set $map "frog" "ribbit"}} - adds a key/value pair to a map

type (long name: .Gtpl.Type)
  expands to "int", "float", "list" or "map": {{$t := type $map}} {{if $t ne "map"}} something is very wrong {{end}}

isint (long name: .Gtpl.IsInt)
  true when its argument is an integer

isfloat (long name: .Gtpl.IsFloat)
  true when its argument is a float

isnumber (long name: .Gtpl.IsNumber)
  true when its argument is an int or a float

islist (long name: .Gtpl.IsList)
  true when its argument is a list (or a slice)

ismap (long name: .Gtpl.IsMap)
  true when its argument is a map

```
