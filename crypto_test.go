package goutil

import (
	"reflect"
	"testing"
)

func TestHashUtils_Encrypt(t *testing.T) {
	const k = "OGU2M2ExZmI3"
	plainText := "123"
	var h = HashUtils{Key: k}
	cipherText := h.Encrypt(plainText)

	want := "wXoWsOe4wOQK8QqlFpBreM7h_zQ="
	if !reflect.DeepEqual(cipherText, want) {
		t.Errorf("HashUtils.Encrypt returned %v, want %v", cipherText, want)
	}
}

func TestHashUtils_HasMatch(t *testing.T) {
	const k = "OGU2M2ExZmI3"
	const p, c = "123", "wXoWsOe4wOQK8QqlFpBreM7h_zQ="
	var h = HashUtils{Key: k}
	if !h.HasMatch(p, c) {
		t.Errorf("HasMatch returned `false`, want `true`")
	}
}
