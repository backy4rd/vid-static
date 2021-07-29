package util

import (
	"errors"
	"math/rand"
	"strings"
)

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateRandomString(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}

func GetFileExtension(filename string) (string, error) {
    dotIdx := strings.LastIndex(filename, ".")
    if dotIdx == -1 {
        return "", errors.New("can't detect file extension")
    }
    if dotIdx == len(filename) - 1 {
        return "", errors.New("can't detect file extension")
    }

    return filename[dotIdx+1:], nil
}
