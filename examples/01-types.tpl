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
