package timecode

import (
	"fmt"
	"math"
	"strconv"
)

// ParseStringToTimcode Will take a string that is either a timecode string or a frame count string
// and return a bonefied timecode object. This is only used in the context of calculate
// where we know if it's df or ndf. This dropframeness is ignored if it's a timecode string
func ParseStringToTimcode(in string, fps float64, excludeLastTimecode bool, dropFrame bool) (*Timecode, error) {
	frames, err := strconv.Atoi(in)
	if err == nil {
		if excludeLastTimecode {
			frames = frames - 1
		}
		return NewTimecodeFromFrames(int64(frames), fps, dropFrame)
	}
	time, err := NewTimecodeFromString(in, fps)
	if excludeLastTimecode && err == nil {
		time.AddFrames(-1)
	}
	return time, err
}

func divmod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}

func formatTimecode(hours int64, minutes int64, seconds int64, frames int64, isDropframe bool) string {
	delim := ":"
	if isDropframe {
		delim = ";"
	}
	return fmt.Sprintf("%02d:%02d:%02d%s%02d", hours, minutes, seconds, delim, frames)
}

func formatTimeSpan(hours int64, minutes int64, seconds int64, ms string) string {
	return fmt.Sprintf("%02d:%02d:%02d.%s", hours, minutes, seconds, ms)
}

func getTimeBase(framerate float64) int {
	return int(math.Ceil(framerate))
}
