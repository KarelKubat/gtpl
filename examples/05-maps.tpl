
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
