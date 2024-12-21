package handlers

import (
	"fmt"
	"strconv"

	"TimecodeTool/timecode"
)

func TimecodeValidate(startTc string, fps float64) *timecode.ValidateResponse {

	firstTc, err := timecode.NewTimecodeFromString(startTc, fps)
	if err != nil {
		return timecode.NewFailedValidateResponse(startTc, fps, err.Error())
	}
	if err := firstTc.Validate(); err != nil {
		return timecode.NewFailedValidateResponse(startTc, fps, err.Error())
	}

	nextFrame, _ := timecode.NewTimecodeFromString(startTc, fps)
	nextFrame.AddFrames(1)

	return timecode.NewOkValidateResponse(startTc, fps, firstTc.DropFrame, nextFrame.GetTimecode())

}

func TimecodeSpan(startTc string, endTc string, fps float64, excludeLastTimecode bool) {

	firstTc, _ := timecode.NewTimecodeFromString(startTc, fps)
	lastTimecode, _ := timecode.NewTimecodeFromString(endTc, fps)
	if excludeLastTimecode {
		lastTimecode.AddFrames(-1)
	}

	fmt.Printf("First Timecode: %s (%s)\n", firstTc.GetTimecode(), firstTc.GetValid())
	fmt.Printf("Last Timecode: %s (%s)\n", lastTimecode.GetTimecode(), lastTimecode.GetValid())

	dfness := "NDF"
	if lastTimecode.DropFrame {
		dfness = "DF"
	}
	println("Last Timecode Frame Index (0 based):", lastTimecode.GetFrameIdx())
	fmt.Printf(
		"Framerate: %s%s\n",
		strconv.FormatFloat(fps, 'f', -1, 64),
		dfness,
	)
}

func TimecodeCalculate(inTc string, operations []string, fps float64, excludeLastTimecode bool) {
	firstTc, _ := timecode.NewTimecodeFromString(inTc, fps)
	curIdx := 0
	for {
		if curIdx >= len(operations)-1 {
			break
		}

		opperator := operations[curIdx]

		nextTime := operations[curIdx+1]

		frames, err := timecode.ParseStringToFrames(nextTime, fps, excludeLastTimecode)
		if err != nil {
			return
		}

		switch opperator {
		case "-":
			firstTc.AddFrames(int(frames) * -1)
		case "+":
			firstTc.AddFrames(int(frames))
		}

		curIdx += 2
	}

	fmt.Println(firstTc.GetTimecode())
}
