{{/*
  Demo of:
    strcat - add elements into one string
*/}}

{{ $ans := 42 }}
{{ $yrs := "7.5 million" }}
{{ strcat "It took " $yrs " to come up with the number " $ans "." }}
