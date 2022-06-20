package util

import (
	"fmt"
	"strings"
)

func RemoveElem[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func AddAudioToWord(word string) string {
	return fmt.Sprintf("%s[sound:%s]", word, word)
}

func RemoveAudioFromWord(word string) string {
	b, _, _ := strings.Cut(word, "[")
	return b
}
