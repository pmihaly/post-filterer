package main

import (
	"fmt"
	"testing"
)

func TestBasic(t *testing.T) {
	var tests = []struct {
		x    bool
		want bool
	}{
		{true, true},
		{false, false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%t", tt.x)
		t.Run(testname, func(t *testing.T) {
			if tt.x != tt.want {
				t.Errorf("got %t, want %t", tt.x, tt.want)
			}
		})
	}
}
