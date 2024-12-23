package timecode

import (
	"fmt"
	"math"
	"strconv"
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

func (t *Timecode) GetFramerateString() string {
	s := strconv.FormatFloat(t.FrameRate, 'f', -1, 64)
	return s
}

func (t *Timecode) GetTimecode() string {
	// calling this will spit out the normalized timecode. For instance, you can instantiate
	// a timecode with a string that contains something like 00:00:10:99 (you can't have 99 frames)
	// But we will run divmod to convert that to a real timecode.
	fq, fr := divmod(int64(t._frames), int64(getTimeBase(t.FrameRate)))
	// println(fq, fr)
	mq, sr := divmod(int64(t._secs)+fq, 60)
	// println(mq, sr)
	hq, mr := divmod(int64(t._mins)+mq, 60)
	// println(hq, mr)
	o, hr := divmod(int64(t._hours)+hq, 24)
	// println(o, hr)
	_ = o

	return formatTimecode(int64(hr), int64(mr), int64(sr), int64(fr), t.DropFrame)
}

func (t *Timecode) GetValid() string {
	isValid := "Valid"
	err := t.Validate()
	if err != nil {
		// isValid = err.Error()
		isValid = "Not Valid"
	}

	return isValid
}

func (t *Timecode) Validate() error {
	// println("+++")
	if t._hours > 23 {
		return fmt.Errorf("Hours cannot be higher than 23")
	}

	if t._mins > 59 {
		return fmt.Errorf("Minutes cannot be higher than 59")
	}

	if t._secs > 59 {
		return fmt.Errorf("Minutes cannot be higher than 59")
	}

	lastAllowedFrame := int(math.Ceil(t.FrameRate)) - 1
	if t._frames > int(lastAllowedFrame) {
		return fmt.Errorf("Frames cannot be higher than %d", lastAllowedFrame)
	}

	if t.DropFrame {
		if 0 <= int(t._frames) || int(t._frames) < 2 {
			// This is potentially wrong for dropframe.
			newTcS := formatTimecode(
				int64(t._hours),
				int64(t._mins-1),
				59,
				int64(lastAllowedFrame),
				t.DropFrame,
			)

			tc, _ := NewTimecodeFromString(newTcS, t.FrameRate)

			valid := false
			for i := 0; i < 3; i++ {
				tc.AddFrames(1)
				if tc.GetTimecode() == t._timecode {
					valid = true
					break
				}
			}

			if !valid {
				return fmt.Errorf("%s is valid timecode for drop frame.", t._timecode)
			}

		}
	}

	return nil
}

func (t *Timecode) PrintPieces() {
	fmt.Println(t._hours, t._mins, t._secs, t._frames)
}

func (t *Timecode) Print() {
	fmt.Println(t.GetTimecode())
}

func (t *Timecode) AddFrames(frames int) {
	if t.DropFrame {
		// todo: investigate why we need this +1
		//  update: I have commented it out, as it seems wrong...
		newFrames := int64(t.GetFrameIdx()) + int64(frames) //+ 1

		if newFrames < 0 {
			// this means we have rolled over backwards.

			// last possible timecode
			// note: This is calculated slightly different than
			// for NDF. In our case, its much easier to construct the
			// last timecode str and then derive the frame versus
			// calculating the last frame.
			lastTimecodeStr := fmt.Sprintf("23:59:59;%d", getTimeBase(t.FrameRate)-1)
			lastTimecode, _ := NewTimecodeFromString(
				lastTimecodeStr,
				t.FrameRate,
			)
			lastFrame := int64(lastTimecode.GetFrameIdx() + 1)
			newFrames = lastFrame + newFrames
		}

		tt, _ := NewTimecodeFromFrames(newFrames, t.FrameRate, t.DropFrame)
		t._hours = tt._hours
		t._mins = tt._mins
		t._secs = tt._secs
		t._frames = tt._frames
		t.DropFrame = tt.DropFrame

	} else {

		newFrames := t._frames + frames
		if newFrames < 0 {
			// this means we have rolled over backwards.

			// last possible frame
			lastFrame := (24 * 60 * 60 * getTimeBase(t.FrameRate))
			newFrames = lastFrame + newFrames
		}

		fq, fr := divmod(int64(newFrames), int64(getTimeBase(t.FrameRate)))
		t._frames = int(fr)

		mq, sr := divmod(int64(t._secs)+fq, 60)
		t._secs = int(sr)

		hq, mr := divmod(int64(t._mins)+mq, 60)
		t._mins = int(mr)

		o, hr := divmod(int64(t._hours)+hq, 24)
		t._hours = int(hr)

		_ = o
	}

}

func (t *Timecode) GetFrameCount() int {
	return t.GetFrameIdx() + 1
}

func (t *Timecode) GetFrameIdx() int {

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
