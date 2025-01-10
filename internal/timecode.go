package internal

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
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

// NewTimecodeFromFrames will create a Timecode object for given frames.
// The only time it will return an error is if DF is specified for a non-DF framerate.
func NewTimecodeFromFrames(inputFrameIdx int64, frameRate float64, isDropframe bool) (*Timecode, error) {

	if isDropframe {
		//CONVERT A FRAME NUMBER TO DROP FRAME TIMECODE
		//Code by David Heidelberger, adapted from Andrew Duncan
		//Given an int called framenumber and a double called framerate
		//Framerate should be 29.97 or 59.94, otherwise the calculations will be off.

		// adapted from https://www.davidheidelberger.com/2010/06/10/drop-frame-timecode/

		// I think JUST this function needs the idx
		// inputFrameIdx := inputFrameIdx - 1

		framenumber := int(inputFrameIdx)

		var d int
		var m int

		dropFrames := int(math.Round(frameRate * 0.066666))               //Number of frames to drop on the minute marks is the nearest integer to 6% of the framerate
		framesPerHour := int(math.Round(frameRate * 60 * 60))             //Number of frames in an hour
		framesPer24Hours := framesPerHour * 24                            //Number of frames in a day - timecode rolls over after 24 hours
		framesPer10Minutes := int(math.Round(frameRate * 60 * 10))        //Number of frames per ten minutes
		framesPerMinute := (int(math.Round(frameRate)) * 60) - dropFrames //Number of frames per minute is the round of the framerate * 60 minus the number of dropped frames

		//In some languages, a % operation will work here
		//But since % for negative numbers varies by language, we'll do it manually
		//Negative time. Add 24 hours.
		for framenumber < 0 {
			framenumber = framesPer24Hours + framenumber
		}

		//If framenumber is greater than 24 hrs, next operation will rollover clock
		framenumber = framenumber % framesPer24Hours //% is the modulus operator, which returns a remainder. a % b = the remainder of a/b

		d = framenumber / framesPer10Minutes // \ means integer division, which is a/b without a remainder. Some languages you could use floor(a/b)
		m = framenumber % framesPer10Minutes

		//In the original post, the next line read m>1, which only worked for 29.97. Jean-Baptiste Mardelle correctly pointed out that m should be compared to dropFrames.
		if m > dropFrames {
			framenumber = framenumber + (dropFrames * 9 * d) + dropFrames*((m-dropFrames)/framesPerMinute)
		} else {
			framenumber = framenumber + dropFrames*9*d
		}

		frRound := int(math.Round(frameRate))
		frames := framenumber % frRound
		seconds := (framenumber / frRound) % 60
		minutes := ((framenumber / frRound) / 60) % 60
		hours := (((framenumber / frRound) / 60) / 60)
		tc_string := formatTimecode(int64(hours), int64(minutes), int64(seconds), int64(frames), true)

		// Fix this deref and deal with the error
		return NewTimecodeFromString(tc_string, frameRate)

	} else {

		sr, frames := divmod(inputFrameIdx, int64(getTimeBase(frameRate)))
		mr, seconds := divmod(sr, 60)
		hr, minutes := divmod(mr, 60)
		//_, _ := divmod(hr, 24)
		tc_string := formatTimecode(hr, minutes, seconds, frames, isDropframe)

		// Fix this deref and deal with the error
		return NewTimecodeFromString(tc_string, frameRate)
	}

}

func NewTimecodeFromString(inputTimecode string, frameRate float64) (*Timecode, error) {

	_timecode := inputTimecode

	re := regexp.MustCompile(`^([0-9]{2}):([0-5][0-9]):([0-5][0-9])[;:]([0-9]{2})$`)

	if !re.MatchString(inputTimecode) {
		return nil, fmt.Errorf("Timecode is malformed. Please format as hh:mm:ss:ff or hh:mm:ss;ff")
	}

	dropFrame := strings.Contains(inputTimecode, ";")

	inputTimecode = strings.Replace(inputTimecode, ";", ":", -1)

	hmsf := strings.Split(inputTimecode, ":")

	_hours, err := strconv.Atoi(hmsf[0])
	if err != nil {
		panic("Hours are malformed.")
	}
	_mins, err := strconv.Atoi(hmsf[1])
	if err != nil {
		panic("Minutes are malformed.")
	}
	_secs, err := strconv.Atoi(hmsf[2])
	if err != nil {
		panic("Seconds are malformed.")
	}
	_frames, err := strconv.Atoi(hmsf[3])
	// println(_frames)
	if err != nil {
		panic("Seconds are malformed.")
	}

	return &Timecode{
		FrameRate: frameRate,
		DropFrame: dropFrame,
		_hours:    _hours,
		_mins:     _mins,
		_secs:     _secs,
		_frames:   _frames,
		_timecode: _timecode,
	}, nil
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

		validDfTimecode := false
		for _, fr := range []float64{29.97, 59.94} {
			if t.FrameRate == fr {
				validDfTimecode = true
				break
			}
		}
		if !validDfTimecode {
			return fmt.Errorf("%s is not a valid framerate for drop frame timecode", t.GetFramerateString())
		}

		fc := t.GetFrameIdx()

		tccTest, err := NewTimecodeFromFrames(int64(fc), t.FrameRate, true)
		if err != nil {
			return err
		}

		if tccTest.GetTimecode() != t.GetTimecode() {
			return fmt.Errorf("%s is not valid drop frame timecode", t.GetTimecode())
		}
	}

	return nil
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

	timeBase := getTimeBase(t.FrameRate)
	frameCount := 0
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
