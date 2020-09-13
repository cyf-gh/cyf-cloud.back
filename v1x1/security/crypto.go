package security

import (
	"crypto/sha512"
	"encoding/hex"
	"github.com/google/uuid"
)

func ToSHA512( str string ) string {
	h := sha512.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func CryptoPasswd( raw string ) string {
	return ToSHA512( raw )
}

func GenerateAccessToken() string {
	return uuid.New().String()
}