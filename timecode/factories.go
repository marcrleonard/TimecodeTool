package timecode

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func NewTimecodeSpan(firstTimecode, lastTimecode *Timecode) (*TimecodeSpan, error) {

	return &TimecodeSpan{
		StartTimecode: firstTimecode,
		LastTimecode:  lastTimecode,
		Framerate:     firstTimecode.FrameRate,
		Dropframe:     firstTimecode.DropFrame,
	}, nil
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
		hours, _ := divmod(hr, 60)
		tc_string := formatTimecode(hours, minutes, seconds, frames, isDropframe)

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
