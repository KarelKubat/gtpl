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
