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
