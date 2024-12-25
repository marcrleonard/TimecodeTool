package internal

import (
	"fmt"
	"math"
	"strings"
)

type TimecodeSpan struct {
	StartTimecode *Timecode
	LastTimecode  *Timecode
	Framerate     float64
	Dropframe     bool
}

func NewTimecodeSpan(firstTimecode, lastTimecode *Timecode) (*TimecodeSpan, error) {

	return &TimecodeSpan{
		StartTimecode: firstTimecode,
		LastTimecode:  lastTimecode,
		Framerate:     firstTimecode.FrameRate,
		Dropframe:     firstTimecode.DropFrame,
	}, nil
}

// TOdo: This calc appears to be wrong
func (t *TimecodeSpan) GetTotalSeconds() float64 {
	tf := t.GetTotalFrames()

	return float64(tf) / t.Framerate
}

// todo: note these two heavily rely on the fact that there is at
//
//	least one frame in the span (see +1)

func (t *TimecodeSpan) GetTotalFrames() int {
	tf := t.LastTimecode.GetFrameIdx() - t.StartTimecode.GetFrameIdx() + 1
	return tf
}

func (t *TimecodeSpan) GetSpanTimecode() string {
	_t, _ := NewTimecodeFromFrames(int64(t.GetTotalFrames()), t.Framerate, t.Dropframe)
	return _t.GetTimecode()
}

func (t *TimecodeSpan) GetSpanRealtime() string {
	_totalSeconds := t.GetTotalSeconds()

	_sq, _sr := divmod(int64(_totalSeconds), int64(60))

	_hq, mr := divmod(int64(_sq), 60)

	_o, hr := divmod(int64(_hq), 24)

	_ = _o

	ms := _totalSeconds - math.Floor(_totalSeconds)

	str := fmt.Sprintf("%.3f", ms)
	// Remove the "0." prefix
	result := strings.TrimPrefix(str, "0.")

	return formatTimeSpan(int64(hr), int64(mr), int64(_sr), result)

}
