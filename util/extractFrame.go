package util

import (
	"bytes"
	"errors"
	"os/exec"
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
