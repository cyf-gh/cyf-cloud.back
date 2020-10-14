package security

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"github.com/google/uuid"
	"math/big"
)

func ToSHA512( str string ) string {
	h := sha512.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func CryptoPasswd( raw string ) string {
	return ToSHA512( raw )
}

// $atk${{uuid}}
func GenerateAtk() string {
	return "$atk$"+uuid.New().String()
}

func GenerateAtkSession() string {
	return "$atk$"+uuid.New().String() + "$session$"
}

func GetRandom() string {
	result, _ := rand.Int(rand.Reader, big.NewInt(9223372036854775807))
	return ToSHA512(result.String())
}