package timecodetool

import (
	"errors"
	"fmt"

	"github.com/marcrleonard/TimecodeTool/internal"
)

func ValidateTimecode(startTc string, fps float64) *ValidateResponse {

	firstTc, err := internal.NewTimecodeFromString(startTc, fps)
	if err != nil {
		return newFailedValidateResponse(startTc, fps, err.Error())
	}
	if err := firstTc.Validate(); err != nil {
		return newFailedValidateResponse(startTc, fps, err.Error())
	}

	nextFrame, _ := internal.NewTimecodeFromString(startTc, fps)
	nextFrame.AddFrames(1)

	return newOkValidateResponse(startTc, fps, firstTc.DropFrame, nextFrame.GetTimecode())

}

func SpanTimecode(startTc string, endTc string, fps float64, excludeLastTimecode bool) *SpanResponse {

	var allErrors []error

	firstTc, err := internal.NewTimecodeFromString(startTc, fps)
	if err != nil {
		allErrors = append(allErrors, fmt.Errorf("First timecode error: %w", err))
	}
	lastTimecode, err := internal.NewTimecodeFromString(endTc, fps)
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
		return newFailedSpanResponse(startTc, endTc, fps, excludeLastTimecode, errors.Join(allErrors...).Error())
	}

	span, err := internal.NewTimecodeSpan(firstTc, lastTimecode)
	if err != nil {
		return newFailedSpanResponse(startTc, endTc, fps, excludeLastTimecode, err.Error())
	}

	nextTimecode, err := internal.NewTimecodeFromString(endTc, fps)
	if err != nil {
		panic("Error generating next timecode")
	}
	nextTimecode.AddFrames(1)

	return newOkSpanResponse(
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
}

func CalculateTimecodes(inTc string, operations []string, fps float64, excludeLastTimecode bool) *CalcResponse {
	firstTc, _ := internal.NewTimecodeFromString(inTc, fps)
	lastTimecode, _ := internal.NewTimecodeFromString(inTc, fps)
	curIdx := 0

	calcSteps := []CalculationStep{}

	for {
		if curIdx >= len(operations)-1 {
			break
		}

		opperator := operations[curIdx]

		nextTime := operations[curIdx+1]

		nexTc, err := internal.ParseStringToTimecode(nextTime, fps, excludeLastTimecode, firstTc.DropFrame)

		if err != nil {
			return newFailedCalcResponse(
				inTc,
				"",
				fps,
				excludeLastTimecode,
				err.Error(),
				[]CalculationStep{},
			)
		}

		switch opperator {
		case "-":
			lastTimecode.AddFrames(int(nexTc.GetFrameCount()) * -1)
		case "+":
			lastTimecode.AddFrames(int(nexTc.GetFrameCount()))
		}

		calcSteps = append(calcSteps, CalculationStep{
			Operation: opperator,
			Timecode:  nexTc.GetTimecode(), // fix this
			Frames:    int(nexTc.GetFrameCount()),
		})

		curIdx += 2
	}

	span := SpanTimecode(inTc, lastTimecode.GetTimecode(), fps, excludeLastTimecode)

	return newOkCalcResponse(
		firstTc.GetTimecode(),
		lastTimecode.GetTimecode(),
		fps,
		firstTc.DropFrame,
		excludeLastTimecode,
		firstTc.GetFrameIdx(),
		lastTimecode.GetFrameIdx(),
		span.LengthFrames,
		span.LengthTime,
		span.LengthTimecode,
		span.LengthSeconds,
		span.NextTimecode,
		calcSteps,
	)

}
