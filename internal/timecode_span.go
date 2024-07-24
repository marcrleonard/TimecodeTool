package internal

import (
	"fmt"
	"strings"
)

type TimecodeSpan struct {
	StartTimecode Timecode
	EndTimecode   Timecode
	Framerate     float64
	Dropframe     bool
}

func (t *TimecodeSpan) GetTotalFrames() int {
	tf := t.EndTimecode.GetFrameIdx() - t.StartTimecode.GetFrameIdx() + 1
	//	return int(math.Max(float64(tf), 1))

	return tf
}

func (t *TimecodeSpan) GetSpanTimecode() string {
	_t := Timecode{}
	_t.DropFrame = t.Dropframe
	_t._hours = 0
	_t._mins = 0
	_t._secs = 0
	_t._frames = 0
	_t.AddFrames(t.GetTotalFrames())

	return _t.GetTimecode()

}

func (t *TimecodeSpan) GetSpanRealtime() string {
	_totalFrames := t.GetTotalFrames()

	// TODO: I HAVE NO CLUE IF THIS IS t.Framerate or getTimeBase()
	_fq, fr := divmod(int64(_totalFrames), int64(t.Framerate))
	// _fq, fr := divmod(int64(_totalFrames), int64(getTimeBase(t.Framerate)))

	_mq, sr := divmod(int64(_fq), 60)

	// fmt.Println(mq, sr)
	_hq, mr := divmod(int64(_mq), 60)

	_o, hr := divmod(int64(_hq), 24)

	_ = _o

	ms := (float64(fr) / float64(t.Framerate))

	str := fmt.Sprintf("%.3f", ms)
	// Remove the "0." prefix
	result := strings.TrimPrefix(str, "0.")

	return formatTimeSpan(int64(hr), int64(mr), int64(sr), result)

}
