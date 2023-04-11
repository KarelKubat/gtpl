#!/usr/bin/env perl
use File::Basename;
use strict;

my $top = dirname($0);
$top =~ s{tools$}{};
if ($top ne "") {
    chdir($top) or die("cannot cd to $top: $!");
}

print(
    "## List of Built in Functions\n\n",
    "The list can be generated using `gtpl -b`.\n",
    "The lowercase aliases (e.g., `add` for `.Gtpl.Add`) are **not** available\n",
    "when the flag `--allow-aliases=false` is given. \n\n",
    "```\n");

open(my $if, "go run gtpl.go -b |") or die;
my $out = do { local $/; <$if> };
close($if) or die;

print(
    $out,
    "```\n")