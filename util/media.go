package util

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
    "strings"
	"strconv"
)

func ExtractFrame(filePath string, seek int, height int, dest string) error {
    cmd := exec.Command(
        "ffmpeg",
        "-ss", strconv.Itoa(seek),
        "-i", filePath,
        "-y",
        "-vframes", "1",
        "-filter:v",
        "scale=w=trunc(oh*a/2)*2:h=" + strconv.Itoa(height),
        dest,
    )
    var stderr bytes.Buffer
    cmd.Stderr = &stderr

    err := cmd.Run()
    if err != nil {
        return errors.New(stderr.String())
    }
    return nil
}

func GetVideoDuration(filePath string) (int, error) {
    if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
        return -1, err
    }

    cmd := exec.Command(
        "ffprobe",
        "-i", filePath,
        "-show_entries", "format=duration",
        "-v", "error",
        "-of", "csv=p=0",
    )
    var stdout bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    err := cmd.Run()
    if err != nil {
        return -1, errors.New(stderr.String())
    }

    duration, err := strconv.ParseFloat(strings.TrimSuffix(stdout.String(), "\n"), 32);
    if err != nil {
        return -1, err
    }

    return int(duration), nil
}
