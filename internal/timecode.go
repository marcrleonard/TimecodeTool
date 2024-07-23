package internal

import (
	"errors"
	"fmt"
	"math"
)

type Timecode struct {
	FrameRate float64
	DropFrame bool
	_hours    int
	_mins     int
	_secs     int
	_frames   int
	_timecode string
}

func (t *Timecode) getTimecode() string {
	// calling this will spit out the normalized timecode. For instance, you can instantiate
	// a timecode with a string that contains something like 00:00:10:99 (you can't have 99 frames)
	// But we will run divmod to convert that to a real timecode.

	fq, fr := divmod(int64(t._frames), int64(getTimeBase(t.FrameRate)))

	mq, sr := divmod(int64(t._secs)+fq, 60)

	hq, mr := divmod(int64(t._mins)+mq, 60)

	o, hr := divmod(int64(t._hours)+hq, 24)

	_ = o

	return formatTimecode(int64(hr), int64(mr), int64(sr), int64(fr), t.DropFrame)
}

func (t *Timecode) Validate() error {

	asd := TimecodeFromFrames(int64(t.getFrameCount()), t.FrameRate, t.DropFrame)

	if asd.getTimecode() != t._timecode {
		return errors.New(fmt.Sprintf("Timecode is not valid! %s != %s", asd.getTimecode(), t.getTimecode()))
	} else {
		return nil
	}
}

func (t *Timecode) Print() {
	fmt.Println(t.getTimecode())
}

func (t *Timecode) getFrameCount() int {

	var timeBase int = getTimeBase(t.FrameRate)
	var frameCount int = 0
	if t.DropFrame == false {
		hrsToSecs := t._hours * 60 * 60
		minsToSecs := t._mins * 60
		totalSeconds := hrsToSecs + minsToSecs + t._secs

		frameCount = (totalSeconds * timeBase) + t._frames
	} else {
		// //CONVERT DROP FRAME TIMECODE TO A FRAME NUMBER
		//Code by David Heidelberger, adapted from Andrew Duncan
		//Given ints called hours, minutes, seconds, frames, and a double called framerate

		// adapted from https://www.davidheidelberger.com/2010/06/10/drop-frame-timecode/

		dropFrames := int(math.Round(t.FrameRate * 0.066666)) //Number of drop frames is 6% of framerate rounded to nearest integer
		// Should this ^ be round or int? todo: This could be wrong.

		hourFrames := timeBase * 60 * 60          //Number of frames per hour (non-drop)
		minuteFrames := timeBase * 60             //Number of frames per minute (non-drop)
		totalMinutes := (60 * t._hours) + t._mins //Total number of minuts
		frameCount = ((hourFrames * t._hours) + (minuteFrames * t._mins) + (timeBase * t._secs) + t._frames) - (dropFrames * (totalMinutes - (totalMinutes / 10)))

	}

	return frameCount
}

func (t *Timecode) AddFrames(frames int) {
	if t.DropFrame {
		// todo: investigate why we need this +1
		newFC := int64(t.getFrameCount()) + int64(frames) + 1

		tt := TimecodeFromFrames(newFC, t.FrameRate, t.DropFrame)
		t._hours = tt._hours
		t._mins = tt._mins
		t._secs = tt._secs
		t._frames = tt._frames
		t.DropFrame = tt.DropFrame

		// fmt.Println(t.getTimecode(), tt.getTimecode())

	} else {

		// Do we need the divmod below?

		newFrames := t._frames + frames
		print(newFrames, "\n")
		fq, fr := divmod(int64(t._frames+frames), int64(getTimeBase(t.FrameRate)))
		t._frames = int(fr)
		// fmt.Println(fq, fr)

		mq, sr := divmod(int64(t._secs)+fq, 60)
		t._secs = int(sr)

		// fmt.Println(mq, sr)

		hq, mr := divmod(int64(t._mins)+mq, 60)
		t._mins = int(mr)

		// fmt.Println(hq, mr)

		o, hr := divmod(int64(t._hours)+hq, 24)
		t._hours = int(hr)

		_ = o
	}

}

func (t *Timecode) GetFrameIdx() int {
	return t.getFrameCount() - 1
}
