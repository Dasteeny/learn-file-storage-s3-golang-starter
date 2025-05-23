package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
)

func getVideoAspectRatio(filePath string) (string, error) {
	command := exec.Command("ffprobe",
		"-v", "error",
		"-print_format", "json",
		"-show_streams",
		filePath,
	)
	var stdout bytes.Buffer
	command.Stdout = &stdout

	if err := command.Run(); err != nil {
		return "", fmt.Errorf("ffprobe error: %v", err)
	}

	var output struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		return "", fmt.Errorf("could not parse ffprobe output: %v", err)
	}

	if len(output.Streams) == 0 {
		return "", errors.New("no video streams found")
	}

	width := output.Streams[0].Width
	height := output.Streams[0].Height

	if width == 16*height/9 {
		return "16:9", nil
	} else if height == 16*width/9 {
		return "9:16", nil
	}
	return "other", nil
}

func processVideoForFastStart(filePath string) (string, error) {
	processedFilePath := filePath + ".processing"

	command := exec.Command("ffmpeg",
		"-i", filePath,
		"-movflags", "faststart",
		"-codec", "copy",
		"-f", "mp4",
		processedFilePath,
	)
	var stderr bytes.Buffer
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		return "", fmt.Errorf("error processing video: %s, %v", stderr.String(), err)
	}

	return processedFilePath, nil
}
