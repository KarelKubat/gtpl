{{/*
  Demo of:
    strcat   - add elements into one string
    contains - `true` when a string contains a substring
    assert   - ensures a condition or exits
  Also standard built ins:
    len      - returns the length of a string
    index    - returns the ordinal number of the rune at a given index
*/}}

{{ $ans := 42 }}
{{ $yrs := "7.5 million" }}
{{ $out := strcat "It took " $yrs " to come up with the number " $ans "." }}

"{{ $out }}" is {{ len $out }} runes long.
The ordinal of the rune at position 12 is {{ index $out 12 }}.

{{ assert (contains $yrs "million") "assertion failure, quitting" }}
