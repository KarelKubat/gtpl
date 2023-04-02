This template is processed by {{ expander }} version {{ version }}

{{ log "This generates " "one" " log statement" }}

{{ $answer := 42 }}
{{/* abort if it isn't 42 */}}
{{ assert (eq $answer 42) "answer is " $answer " but i need 42" }}

{{/* Uncomment to cause an error.
  {{ die "stop execution with an error" }}
*/}}
