#!/usr/bin/env perl
use strict;

die('usage: gendoc FILE.TPL') unless ($#ARGV eq 0);

open(my $srcf, $ARGV[0]) or die;
my $file_content = do { local $/; <$srcf> };
close($srcf) or die;

my $cmd = "go run gtpl.go -re $ARGV[0]";
open(my $gtplf, "$cmd 2>&1 |") or die;
my $expanded_content = do { local $/; <$gtplf> };
close($gtplf) or die ("$cmd failed: $expanded_content");

print(
    "### Example: $ARGV[0]\n\n",
    "```\n",
    $file_content,
    "```\n\n",
    "**Output** (empty lines removed):\n\n",
    "```\n",
    $expanded_content,
    "```\n");
