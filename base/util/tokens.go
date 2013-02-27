package util

import (
	"math/rand"
)

func CreateToken(token_size int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	token := ""
	for i := 0; i < token_size; i++ {
		token = token + string(chars[rand.Intn(len(chars))])
	}
	return token
}
