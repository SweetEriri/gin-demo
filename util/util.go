package util

import "math/rand"

func RandomString(n int) string {
	var letter = []byte("asdfghjklqwertyuiopASDFGHJKLQWERTYUIOP")
	result :=make([]byte, n)
	for i := range result{
		result[i] = letter[rand.Intn(len(letter))]
	}

	return string(result)
}