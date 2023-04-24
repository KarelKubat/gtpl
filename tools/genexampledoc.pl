#!/usr/bin/env perl
use strict;

die('usage: genexampledoc.pl FILE.TPL') unless ($#ARGV eq 0);

open(my $srcf, $ARGV[0]) or die;
my $file_content = do { local $/; <$srcf> };
close($srcf) or die;
$file_content .= "\n" unless ($file_content =~ m{\n$}m);

my $cmd = "go run gtpl.go -re $ARGV[0]";
open(my $gtplf, "$cmd 2>&1 |");
my $expanded_content = do { local $/; <$gtplf> };
close($gtplf) or die ("$cmd failed: $expanded_content");
$expanded_content .= "\n" unless ($expanded_content =~ m{\n$}m);

print(
    "\n",
    "### Example: $ARGV[0]\n",
    "\n",
    "```C\n",
    $file_content,
    "```\n",
    "\n",
    "**Output** (empty lines removed):\n",
    "\n",
    "```plain\n",
    $expanded_content,
    "```\n");
