package playlist

import (
	"Test-Task-Go/internal/entity"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func InsertAdsIntoPlaylist(playlist string, creative entity.Creative) (string, error) {
	lines := strings.Split(playlist, "\n")
	var header []string
	var segments []string
	var modifiedPlaylist []string
	inCueOut := false

	for _, line := range lines {
		if strings.HasPrefix(line, "#EXTM3U") || strings.HasPrefix(line, "#EXT-X-VERSION") ||
			strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE") || strings.HasPrefix(line, "#EXT-X-TARGETDURATION") {
			header = append(header, line)
			continue
		}

		if len(header) > 0 && (strings.HasPrefix(line, "#EXTINF") || strings.HasPrefix(line, "#EXT-X-CUE-OUT") ||
			strings.HasPrefix(line, "#EXT-X-CUE-IN") || strings.HasPrefix(line, "seg")) {
			segments = append(segments, line)
		}
	}

	for _, line := range header {
		modifiedPlaylist = append(modifiedPlaylist, line)
	}

	re := regexp.MustCompile(`#EXTINF:(\d+\.\d+),\s*(\S+)`)

	for _, line := range segments {
		if strings.HasPrefix(line, "#EXT-X-CUE-OUT") {
			inCueOut = true
			modifiedPlaylist = append(modifiedPlaylist, line)
			continue
		}

		if strings.HasPrefix(line, "#EXT-X-CUE-IN") {
			matches := re.FindAllStringSubmatch(creative.PlaylistHLS, -1)

			for _, match := range matches {
				modifiedPlaylist = append(modifiedPlaylist, "#EXTINF:"+match[1]+",")
				modifiedPlaylist = append(modifiedPlaylist, match[2])
			}
			modifiedPlaylist = append(modifiedPlaylist, line)
			inCueOut = false
			continue
		}

		if !inCueOut {
			modifiedPlaylist = append(modifiedPlaylist, line)
		}
	}

	return strings.Join(modifiedPlaylist, "\n"), nil
}

func GetCueOutDuration(playlist string) (int, error) {
	re := regexp.MustCompile(`#EXT-X-CUE-OUT:(\d+)`)
	match := re.FindStringSubmatch(playlist)
	if len(match) < 2 {
		return 0, fmt.Errorf("cue-out duration not found")
	}

	duration, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, fmt.Errorf("error parsing cue-out duration: %w", err)
	}
	return duration, nil
}
