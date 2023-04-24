{{/*
    Demo of:
        map   - creates a map
    Also standard built ins:
        range - ranging key,value pairs over a map

Note that ranging over maps has an undefined order, which may be
alphabetical by key, but that is not guaranteed.
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
