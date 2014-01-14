package goutil

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

type HashUtils struct {
	Key string
}

// Encrypt encrypts plaintext with the key.
// Returned ciphertext result
func (u *HashUtils) Encrypt(plain string) string {
	h := hmac.New(sha1.New, []byte(u.Key))
	h.Write([]byte(plain))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

// HasMatch check whether the encrpytion of provided plain-text with key equal with expected cipher
func (u *HashUtils) HasMatch(plain, cipher string) bool {
	msgHash := u.Encrypt(plain)
	return msgHash == cipher
}
