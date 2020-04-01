package engine

import (
	"testing"
)

func TestGeneratePIN(t *testing.T) {
	in := make([]byte, 32)
	want := "7-9-3-2-7-4"
	got := generatePIN(in)
	if want != got {
		t.Error("expect", want, "got", got)
	}
}

func TestGenerate(t *testing.T) {
	in := make([]byte, 32)
	want := "a-a-a-a"
	got, err := Generate(in, PlainLowerCase)
	if err != nil {
		t.Error(err)
	}
	if want != got {
		t.Error("expect", want, "got", got)
	}
	want = "A-A-A-A-#1"
	got, err = Generate(in, WithSpecialCharacters)
	if err != nil {
		t.Error(err)
	}
	if want != got {
		t.Error("expect", want, "got", got)
	}
}
