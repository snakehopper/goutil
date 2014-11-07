package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5Hash(plain string) string {
	hasher := md5.New()
	hasher.Write([]byte(plain))
	return hex.EncodeToString(hasher.Sum(nil))
}
