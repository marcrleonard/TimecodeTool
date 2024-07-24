package internal

import (
	"math"
	"strconv"
	"strings"
)

func NewTimecodeSpan(startTimecode string, endTimecode string, frameRate float64) TimecodeSpan {

	if strings.Contains(startTimecode, ";") != strings.Contains(endTimecode, ";") {
		panic("startTimecode and endTimecode must both be dropframe or non-dropframe.")
	}

	return TimecodeSpan{
		StartTimecode: TimecodeFromString(startTimecode, frameRate),
		EndTimecode:   TimecodeFromString(endTimecode, frameRate),
		Framerate:     frameRate,
		Dropframe:     strings.Contains(startTimecode, ";") || strings.Contains(endTimecode, ";"),
	}
}

func TimecodeFromFrames(inputFrameCount int64, frameRate float64, isDropframe bool) Timecode {

	// if inputFrameCount < 1 {
	// 	panic("Framecount must be >= 1")
	// }

	// inputFrameIdx := inputFrameCount - 1

	tcObj := Timecode{}
	if isDropframe {
		//CONVERT A FRAME NUMBER TO DROP FRAME TIMECODE
		//Code by David Heidelberger, adapted from Andrew Duncan
		//Given an int called framenumber and a double called framerate
		//Framerate should be 29.97 or 59.94, otherwise the calculations will be off.

		// adapted from https://www.davidheidelberger.com/2010/06/10/drop-frame-timecode/

		// I think JUST this function needs the idx
		inputFrameIdx := inputFrameCount - 1

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
		tcObj = TimecodeFromString(tc_string, frameRate)

	} else {

		sr, frames := divmod(inputFrameCount, int64(getTimeBase(frameRate)))
		mr, seconds := divmod(sr, 60)
		hr, minutes := divmod(mr, 60)
		hours, _ := divmod(hr, 60)
		tc_string := formatTimecode(hours, minutes, seconds, frames, isDropframe)
		tcObj = TimecodeFromString(tc_string, frameRate)
	}

	return tcObj

}

func TimecodeFromString(inputTimecode string, frameRate float64) Timecode {

	_timecode := inputTimecode

	if len(inputTimecode) != 11 {
		panic("Timecode is malformed. Please format as hh:mm:ss:ff")
	}

	var tc Timecode

	dropFrame := strings.Contains(inputTimecode, ";")

	if dropFrame {
		inputTimecode = strings.Replace(inputTimecode, ";", ":", -1)
	}

	hmsf := strings.Split(inputTimecode, ":")

	if len(hmsf) != 4 {
		panic("Timecode is malformed. Please format as hh:mm:ss:ff")
	}

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

	tc._hours = _hours
	tc._mins = _mins
	tc._secs = _secs
	tc._frames = _frames
	tc._timecode = _timecode
	tc.FrameRate = frameRate
	tc.DropFrame = dropFrame

	return tc
}
