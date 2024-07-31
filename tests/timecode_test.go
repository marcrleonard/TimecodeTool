package main

import (
	"TimecodeTool/internal"
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestNDFIndexes(t *testing.T) {

	tc := internal.TimecodeFromFrames(0, 23.98, false)
	if tc.GetFrameIdx() != 0 {
		t.Fatalf(`%v != 0`, tc.GetFrameIdx())
	}

}

func TestAddNDFIndexes(t *testing.T) {

	tc := internal.TimecodeFromFrames(0, 23.98, false)
	tc.AddFrames(1)
	if tc.GetFrameIdx() != 1 {
		t.Fatalf(`%v != 1`, tc.GetFrameIdx())
	}

	tc.AddFrames(4)
	if tc.GetFrameIdx() != 5 {
		t.Fatalf(`%v != 5`, tc.GetFrameIdx())
	}

}

func TestDFIndexes(t *testing.T) {

	tcdf := internal.TimecodeFromFrames(0, 29.97, true)

	if tcdf.GetFrameIdx() != 0 {
		t.Fatalf(`%v != 0`, tcdf.GetFrameIdx())
	}
}

func TestDFNonValid(t *testing.T) {

	invalidTimecode := "00:07:00;00"

	tcdf := internal.TimecodeFromString(invalidTimecode, 29.97)

	if tcdf.Validate() != nil {
		t.Fatalf(`%s should not be valid for 29.97 DF timecode`, invalidTimecode)
	}

}

func TestNDFNonValid(t *testing.T) {

	invalidTimecode := "00:07:00;24"

	tcdf := internal.TimecodeFromString(invalidTimecode, 23.98)

	if tcdf.Validate() != nil {
		t.Fatalf(`%s should not be valid for 23.98 NDF timecode`, invalidTimecode)
	}

}
