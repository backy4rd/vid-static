package util

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

func Compress360p(input string, dest string) (error) {
    if _, err := os.Stat(input); errors.Is(err, os.ErrNotExist) {
        return err
    }

    cmd := exec.Command(
        "ffmpeg",
        "-y",
        "-i", input,
        "-c:a", "copy",
        "-c:v", "libx264",
        "-vf", "scale=-2:360",
        "-preset", "ultrafast",
        "-strict", "-2",
        dest,
    )
    var stdout bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    err := cmd.Run()
    if err != nil {
        return errors.New(stderr.String())
    }

    return nil;
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

func GetVideoQuality(filePath string) (int, error) {
    if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
        return -1, err
    }

    cmd := exec.Command(
        "ffprobe",
        "-i", filePath,
        "-show_entries", "stream=height",
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

    height, err := strconv.ParseInt(strings.TrimSuffix(stdout.String(), "\n\n"), 10, 32);
    if err != nil {
        return -1, err
    }

    if height < 420 {
        return 360, nil
    } else if height < 600 {
        return 480, nil
    } else if height < 900 {
        return 720, nil
    } else {
        return 1080, nil
    }
}

/*
   GoLang: os.Rename() give error "invalid cross-device link" for Docker container with Volumes.
   MoveFile(source, destination) will work moving file between folders
*/

func MoveFile(sourcePath, destPath string) error {
    inputFile, err := os.Open(sourcePath)
    if err != nil {
        return err
    }
    outputFile, err := os.Create(destPath)
    if err != nil {
        inputFile.Close()
        return err
    }
    defer outputFile.Close()
    _, err = io.Copy(outputFile, inputFile)
    inputFile.Close()
    if err != nil {
        return err
    }
    // The copy was successful, so now delete the original file
    err = os.Remove(sourcePath)
    if err != nil {
        return err
    }
    return nil
}
