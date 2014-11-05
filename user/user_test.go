package user

import (
	"testing"
)

func TestNewValidateDigit(t *testing.T) {
	d := "1234"
	vd := NewValidateDigit(d)
	if vd.Value != d {
		t.Errorf("NewValidateDigit.Value return %v, want %v", vd.Value, d)
	}
}
func TestValidateDigit_HasMatch(t *testing.T) {
	d := "234"
	storedDigit := "0234"
	vd := NewValidateDigit(storedDigit)
	if !vd.HasMatch(d) {
		t.Errorf("HasMatch return false, want true")
	}
}
func TestValidateDigit_HasMatch_failed(t *testing.T) {
	d := "1234"
	storedDigit := "0234"
	vd := NewValidateDigit(storedDigit)
	if vd.HasMatch(d) {
		t.Errorf("HasMatch return true, want false")
	}
}
func TestRandomValidateDigit(t *testing.T) {
	length := 6
	vd := RandomValidateDigit(length)
	t.Logf("%v", vd.Value)
	if len(vd.Value) != length {
		t.Errorf("Length of value return %d, want %d", len(vd.Value), length)
	}
}
