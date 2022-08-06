package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(length int) string {
	var sb strings.Builder
	num := len(alphabet)
	for i := 0; i < length; i++ {
		c := alphabet[rand.Intn(num)]
		sb.WriteByte(c)
	}
	return sb.String()
}

//这里使用的为微信openid，长度为28
func RandomSessionid() string {
	return RandomString(28)
}
