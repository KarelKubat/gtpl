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

{{/* `loop 1 11` is a shorthand for `list 1 2 3 4 5 6 7 8 9 10` */}}
{{/* That means "up to 11", not "and including". */}}
{{ range $i := loop 1 11 }}
  Number {{ $i }}: {{ $a }}
  {{ $tmp := $a }}
  {{ $a = $b }}
  {{ $b = add $tmp $b }}
{{ end }}
