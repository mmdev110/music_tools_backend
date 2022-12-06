package main

import "testing"

func TestSum(t *testing.T) {
	a := 3
	b := 4
	got := 7
	want := a + b
	if got != want {
		t.Errorf("got %d, want %d .", got, want)
	}
}
