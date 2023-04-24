package main

import "testing"

// There isn't a lot to test in package main. All relevant tests occur in sub-packages.

func TestCheck(t *testing.T) {
	check(nil) // If this does an os.Exit(1) then the test suite will fail.
}
