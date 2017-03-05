package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

func getShortHash(f io.Reader) string {

	hash := sha256.New()
	io.Copy(hash, f)
	key := hex.EncodeToString(hash.Sum(nil))

	return key[:6]

}
