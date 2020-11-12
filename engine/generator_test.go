package engine

import (
	"testing"
)

func TestGeneratePIN(t *testing.T) {
	in := make([]byte, 32)
	want := "0-0-0-0-0-0"
	got, _ := generatePIN(in)
	if want != got {
		t.Error("expect", want, "got", got)
	}
}

func TestGeneratePassword(t *testing.T) {
	in := make([]byte, 32)
	want := "aaaaaaaaaaaaa#A"
	got, _ := generatePassword(in)
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
