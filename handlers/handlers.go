package handlers

import (
	"errors"
	"fmt"

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

func NewTimecodeSpan(startTc string, endTc string, fps float64, excludeLastTimecode bool) *timecode.SpanResponse {

	var allErrors []error

	firstTc, err := timecode.NewTimecodeFromString(startTc, fps)
	if err != nil {
		allErrors = append(allErrors, fmt.Errorf("First timecode error: %w", err))
	}
	lastTimecode, err := timecode.NewTimecodeFromString(endTc, fps)
	if err != nil {
		allErrors = append(allErrors, fmt.Errorf("Last timecode error: %w", err))
	}

	if err := firstTc.Validate(); err != nil {
		allErrors = append(allErrors, err)
	}

	if err := lastTimecode.Validate(); err != nil {
		allErrors = append(allErrors, err)
	}

	if len(allErrors) > 0 {
		return timecode.NewFailedSpanResponse(startTc, endTc, fps, excludeLastTimecode, errors.Join(allErrors...).Error())
	}

	span, err := timecode.NewTimecodeSpan(firstTc, lastTimecode)
	if err != nil {
		return timecode.NewFailedSpanResponse(startTc, endTc, fps, excludeLastTimecode, err.Error())
	}

	nextTimecode, err := timecode.NewTimecodeFromString(endTc, fps)
	if err != nil {
		panic("Error generating next timecode")
	}
	nextTimecode.AddFrames(1)

	return timecode.NewOkSpanResponse(
		startTc,
		endTc,
		fps,
		span.Dropframe,
		excludeLastTimecode,
		span.StartTimecode.GetFrameIdx(),
		span.LastTimecode.GetFrameIdx(),
		span.GetTotalFrames(),
		span.GetSpanRealtime(),
		span.GetSpanTimecode(),
		span.GetTotalSeconds(),
		nextTimecode.GetTimecode(),
	)
	//return &timecode.SpanResponse{}
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
