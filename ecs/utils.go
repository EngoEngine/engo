package ecs

import (
	"fmt"
	"math/rand"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateUUID() string {
	return fmt.Sprintf("%X-%X-%X-%X-%X", randSeq(2), randSeq(2), randSeq(2), randSeq(2), randSeq(2))
}
