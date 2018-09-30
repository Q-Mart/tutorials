package main

import (
	"testing"
	"testing/quick"
)

func TestMult(t *testing.T) {
	// Check a function in your code is operating correctly
	f := func(a uint32) bool {
		result := Mult(a)
		return (result == (a * 2))
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
