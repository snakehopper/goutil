package strconv

import (
	"reflect"
	"testing"
)

func TestIncreaseInt(t *testing.T) {
	const in = "1"
	out, err := IncreaseInt(in)

	if err != nil {
		t.Errorf("strconv.IncreaseInt returned error: %v", err)
	}

	want := "2"
	if !reflect.DeepEqual(out, want) {
		t.Errorf("strconv.IncreaseInt returned %v, want %v", out, want)
	}
}

func TestIncreaseInt_parseError(t *testing.T) {
	const in = "1e"
	if _, err := IncreaseInt(in); err == nil {
		t.Errorf("Expected error.")
	}
}

func TestDecreaseInt(t *testing.T) {
	const in = "2"
	out, err := DecreaseInt(in)

	if err != nil {
		t.Errorf("strconv.IncreaseInt returned error: %v", err)
	}

	want := "1"
	if !reflect.DeepEqual(out, want) {
		t.Errorf("strconv.IncreaseInt returned %v, want %v", out, want)
	}
}

func TestDecreaseInt_parseError(t *testing.T) {
	const in = "1e"
	if _, err := DecreaseInt(in); err == nil {
		t.Errorf("Expected error.")
	}

}
