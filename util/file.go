package util

import (
	"errors"
	"os"
)

func IsFileExist(filePath string) bool {
    if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
        return false
    } else {
        return true
    }
}
