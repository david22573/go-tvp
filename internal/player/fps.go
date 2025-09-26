package player

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// FFProbe response structures
type FFProbeResponse struct {
	Streams []Stream `json:"streams"`
}

type Stream struct {
	CodecType    string `json:"codec_type"`
	AvgFrameRate string `json:"avg_frame_rate"`
	RFrameRate   string `json:"r_frame_rate"`
}

// DetectFPS extracts the frame rate from a video file using ffprobe
func DetectFPS(videoPath string) (float64, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_streams",
		"-select_streams", "v:0",
		videoPath)

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to run ffprobe: %w", err)
	}

	var response FFProbeResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return 0, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	if len(response.Streams) == 0 {
		return 0, fmt.Errorf("no video streams found")
	}

	stream := response.Streams[0]
	if stream.CodecType != "video" {
		return 0, fmt.Errorf("first stream is not video")
	}

	// Try avg_frame_rate first, then r_frame_rate
	fps, err := parseFrameRate(stream.AvgFrameRate)
	if err != nil || fps <= 0 {
		fps, err = parseFrameRate(stream.RFrameRate)
		if err != nil || fps <= 0 {
			return 0, fmt.Errorf("could not parse frame rate")
		}
	}

	return fps, nil
}

// parseFrameRate converts ffprobe frame rate strings like "30000/1001" to float64
func parseFrameRate(frameRateStr string) (float64, error) {
	if frameRateStr == "" || frameRateStr == "0/0" {
		return 0, fmt.Errorf("invalid frame rate string")
	}

	parts := strings.Split(frameRateStr, "/")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid frame rate format")
	}

	numerator, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, err
	}

	denominator, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, err
	}

	if denominator == 0 {
		return 0, fmt.Errorf("division by zero in frame rate")
	}

	return numerator / denominator, nil
}
