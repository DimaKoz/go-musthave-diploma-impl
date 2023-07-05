package main

import "testing"

func TestDoNothing(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{name: "test code cov", want: 0},
	}

	for _, tt := range tests {
		ttt := tt
		t.Run(ttt.name, func(t *testing.T) {
			if got := DoNothing(); got != ttt.want {
				t.Errorf("DoNothing() = %v, want %v", got, ttt.want)
			}
		})
	}
}
