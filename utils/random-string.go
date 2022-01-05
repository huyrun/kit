package utils

import "math/rand"

func RandomString(n int, alphabet string) string {
	var letters = []rune(alphabet)

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
