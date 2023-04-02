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
