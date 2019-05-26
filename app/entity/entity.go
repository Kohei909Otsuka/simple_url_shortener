package entity

import (
	"math/rand"
	"os"
	"time"
)

const (
	tokenLength = 6
)

var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type OriginalUrl string
type ShortenUrl string

func (o OriginalUrl) Shorten() ShortenUrl {
	base := os.Getenv("BASE_URL")
	return ShortenUrl(base + "/" + genRandomToken())
}

func genRandomToken() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, tokenLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
