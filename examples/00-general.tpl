{{/*
  Demo of:
    expander - the name of the expander program
    version  - its version
    log      - logging information
    assert   - ensuring that a condition is met
    die      - how to abort a run
    env      - returns environment setting
  Also:
    comments - you can't have spaces between the template and comment delimiters
*/}}

{{ log "This generates " "one" " log statement" }}
This template is processed by {{ expander }} version {{ version }}

{{ $answer := 42 }}
{{/* abort if it isn't 42 */}}
{{ assert (eq $answer 42) "answer is " $answer " but i need 42" }}

{{/* Uncomment to cause an error.
  {{ die "stop execution with an error" }}
*/}}

My homedir is {{ env "HOME" }}
