{{/*
  Demo of:
    strcat   - add elements into one string
    contains - `true` when a string contains a substring
*/}}

{{ $ans := 42 }}
{{ $yrs := "7.5 million" }}
{{ $out := strcat "It took " $yrs " to come up with the number " $ans "." }}

{{ $out }}

{{ assert (contains $yrs "million") "assertion failure, quitting" }}