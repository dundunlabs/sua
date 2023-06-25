package test

import "testing"

func Equal(t *testing.T, got any, want any) {
	if got != want {
		t.Errorf("\ngot : %v\nwant: %v", got, want)
	}
}
