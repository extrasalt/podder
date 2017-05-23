package main

import (
	"strings"
	"testing"
)

func TestGetShortHash(t *testing.T) {

	some_string := "Some string"
	reader := strings.NewReader(some_string)
	hash := getShortHash(reader)

	if len(hash) != 6 {
		t.Fatal("length of the short hash is not the required number")
	}
}
