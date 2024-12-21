package timecode

import (
	"fmt"
	"math"
	"strconv"
)

func ParseStringToFrames(in string, fps float64, excludeLastTimecode bool) (int64, error) {
	frames, err := strconv.Atoi(in)
	if err == nil {
		return int64(frames), nil
	}
	tc, err := NewTimecodeFromString(in, fps)
	if err == nil {
		if excludeLastTimecode {
			return int64(tc.GetFrameIdx()), nil
		}
		return int64(tc.GetFrameCount()), nil
	}

	return 0, fmt.Errorf("Could not parse time of %s", in)
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
