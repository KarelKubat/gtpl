{{/*
    Demo of:
        list  - creates an array 
        map   - creates a map
    Also standard built ins:
        range - ranging a list, or key,value pairs over a map

In this example verses (which are maps) are contained in a list
so that in-order traversal is guaranteed.
*/}}

{{ $lyrics := list
          (map "Monica"    "in my life")
          (map "Erica"     "by my side")
          (map "Rita"      "is all I need")
          (map "Tina"      "is what I see")
          (map "Sandra"    "in the sun")
          (map "Mary"      "all night long")
          (map "Jessica"   "here I am")
          (map "you"       "makes me your man") }}
{{ range $verse := $lyrics }}
  {{ range $name, $what := $verse }}
    A little bit of {{ $name }} {{ $what }}
  {{ end }}
{{ end }}
