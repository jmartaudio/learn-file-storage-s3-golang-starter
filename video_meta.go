package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math"
	"os/exec"
)

type VideoMetaStruct struct {
	Streams []struct {
		Width  int `json:"width,omitempty"`
		Height int `json:"height,omitempty"`
	} `json:"streams"`
}

func getVideoAspectRatio(filePath string) (string, error) {
	var out bytes.Buffer

	log.Println(filePath)

	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	var vidMeta VideoMetaStruct

	if err = json.Unmarshal(out.Bytes(), &vidMeta); err != nil {
		return "", err
	}
	log.Println(vidMeta)

	var vidType string

	ratio := Round(float64(vidMeta.Streams[0].Width) / float64(vidMeta.Streams[0].Height))

	log.Println(ratio)

	switch ratio {
	case 1.78:
		vidType = "landscape"
	case 0.56:
		vidType = "portrait"
	default:
		vidType = "other"
	}

	log.Println(vidType)

	return vidType, nil
}

func Round(x float64) float64 {
	return math.Round(x*100) / 100
}
