package main

import "testing"

func TestSum(t *testing.T) {
	a := 3
	b := 4
	got := a + b
	want := 7
	if got != want {
		t.Errorf("got %d, want %d .", got, want)
	}
}
