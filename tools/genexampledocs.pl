#!/usr/bin/env perl
use File::Basename;
use strict;

my $top = dirname($0);
$top =~ s{tools$}{};
if ($top ne "") {
    chdir($top) or die("cannot cd to $top: $!");
}

print(
    "## Examples of `gtpl` builtins\n\n",
    "See also `examples/*tpl`.\n\n");

my @tpls = sort(glob("examples/*tpl"));
die("no example files") unless ($#tpls >= 0);
for my $f (@tpls) {
    system("tools/genexampledoc.pl $f") and die;
}
