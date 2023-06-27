package test

import (
	"reflect"
	"testing"
)

func Equal(t *testing.T, got any, want any) {
	if got != want {
		t.Errorf("\ngot : %v\nwant: %v", got, want)
	}
}

func DeepEqual(t *testing.T, got any, want any) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot : %v\nwant: %v", got, want)
	}
}
