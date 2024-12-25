package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestNDFIndexes(t *testing.T) {

	tc, _ := NewTimecodeFromFrames(0, 23.98, false)
	if tc.GetFrameIdx() != 0 {
		t.Fatalf(`%v != 0`, tc.GetFrameIdx())
	}

}

func TestAddNDFIndexes(t *testing.T) {

	tc, _ := NewTimecodeFromFrames(0, 23.98, false)
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

	tcdf, _ := NewTimecodeFromFrames(0, 29.97, true)

	if tcdf.GetFrameIdx() != 0 {
		t.Fatalf(`%v != 0`, tcdf.GetFrameIdx())
	}
}

func TestRolloverForwardsDF(t *testing.T) {
	tc, err := NewTimecodeFromString("23:59:59:29", 29.97)
	require.Nil(t, err)
	tc.AddFrames(1)
	require.Equal(t, "00:00:00:00", tc.GetTimecode())
}

// Not implemented!
func TestRolloverBackwardsDF(t *testing.T) {
	tc, err := NewTimecodeFromString("00:00:00;00", 29.97)
	require.Nil(t, err)
	tc.AddFrames(-1)
	require.Equal(t, "23:59:59;29", tc.GetTimecode())
}

func TestRolloverForwardsNDF(t *testing.T) {
	tc, err := NewTimecodeFromString("23:59:59:23", 23.976)
	require.Nil(t, err)
	tc.AddFrames(1)
	require.Equal(t, "00:00:00:00", tc.GetTimecode())
}

func TestRolloverBackwardsNDF(t *testing.T) {
	tc, err := NewTimecodeFromString("00:00:00:00", 23.976)
	require.Nil(t, err)
	tc.AddFrames(-1)
	require.Equal(t, "23:59:59:23", tc.GetTimecode())
}

func TestDFNonValid(t *testing.T) {

	tests := []struct {
		name      string
		error     error
		timecode  string
		framerate float64
	}{
		{
			name:      "Errors when invalid dropframe timecode",
			error:     fmt.Errorf("00:07:00;00 is valid timecode for drop frame."),
			timecode:  "00:07:00;00",
			framerate: 29.97,
		},
		{
			name:      "Errors when invalid frames in timecode field",
			error:     fmt.Errorf("Frames cannot be higher than 23"),
			timecode:  "00:07:00;24",
			framerate: 24,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tcdf, err := NewTimecodeFromString(tt.timecode, tt.framerate)
			err = tcdf.Validate()

			if tt.error != nil {
				require.EqualError(t, err, tt.error.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}

}
