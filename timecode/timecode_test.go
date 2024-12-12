package timecode

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestNDFIndexes(t *testing.T) {

	tc := TimecodeFromFrames(0, 23.98, false)
	if tc.GetFrameIdx() != 0 {
		t.Fatalf(`%v != 0`, tc.GetFrameIdx())
	}

}

func TestAddNDFIndexes(t *testing.T) {

	tc := TimecodeFromFrames(0, 23.98, false)
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

	tcdf := TimecodeFromFrames(0, 29.97, true)

	if tcdf.GetFrameIdx() != 0 {
		t.Fatalf(`%v != 0`, tcdf.GetFrameIdx())
	}
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
			error:     fmt.Errorf("00:07:00;00 is not a valid timecode."),
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
