{{/*
  Demo of:
    strcat   - add elements into one string
    contains - `true` when a string contains a substring
    assert   - ensures a condition or exits
    addbyte  - adds an `int` value as the next byte to a string
  Also standard built ins:
    len      - returns the length of a string
    index    - returns byte at a given index
    printf   - `fmt.Printf()` like expansion
*/}}

{{/* strcat example */}}
{{ $ans := 42 }}
{{ $yrs := "7.5 million" }}
{{ $out := strcat "It took " $yrs " to come up with the number " $ans "." }}
"{{ $out }}" is {{ len $out }} is 50 bytes long.
The byte at position 12 is {{ index $out 12 }}.
{{ assert (contains $yrs "million") "assertion failure, quitting" }}

{{/* NOTE: instead of strcat, the builtin printf can be used: */}}
{{ $msg := printf "'%v' and '%v' are from the HHGttG." $ans $yrs }}
{{ $msg }}

{{/* string to bytes and reverse */}}
{{ $str := "Hello in Chinese is 你好" }}
{{ $values := list }}
{{ range $i := loop 0 (len $str) }}
  {{ $values = addelements $values (index $str $i) }}
{{ end }}
"{{ $str }}" as values is:
  {{ $values }}

{{ $new := "" }}               
{{ range $b := $values }}
  {{ $new = addbyte $new $b }}
{{ end }}
{{ $values }} as string is:
  "{{ $new }}"

{{ assert (eq $str $new) "To/from values conversion failed" }}
