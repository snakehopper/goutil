package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5Hash(plain string) string {
	return MD5HashRaw([]byte(plain))
}

func MD5HashRaw(b []byte) string {
	hasher := md5.New()
	hasher.Write(b)
	return hex.EncodeToString(hasher.Sum(nil))
}
