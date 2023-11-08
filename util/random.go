package util

import "math/rand"

func RandomString(n int) string {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func RandomUserName() string {
	return RandomString(6)
}

func RandomRole() string {
	return RandomString(6)
}

func RandomEmail() string {
	return RandomString(6) + "@" + RandomString(4) + ".com"
}
