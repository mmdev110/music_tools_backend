package main

import (
	"os"
	"testing"
)

var expiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEwMCwidG9rZW5fdHlwZSI6ImFjY2VzcyIsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6NTAwMCIsImF1ZCI6WyJodHRwOi8vbG9jYWxob3N0OjMwMDAiXSwiZXhwIjoxNjkyNTA4ODc0LCJuYmYiOjE2OTI1MDg4MTQsImlhdCI6MTY5MjUwODgxNH0.Nw5g5FYh_uiZkvOg0bhxV0nIP_Z75lYZ72xjwOArbL0"

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}
