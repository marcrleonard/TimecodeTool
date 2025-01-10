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

func Test_TimecodeIndexes(t *testing.T) {
	tests := []struct {
		name        string
		timecode    string
		framerate   float64
		expectedIdx int
	}{
		{
			name:        "29.97",
			timecode:    "00:00:00:00",
			framerate:   29.97,
			expectedIdx: 0,
		},
		{
			name:        "29.97",
			timecode:    "00:00:01:00",
			framerate:   29.97,
			expectedIdx: 30,
		},
		{
			name:        "29.97",
			timecode:    "01:00:00:00",
			framerate:   29.97,
			expectedIdx: 108000,
		},
		{
			name:        "29.97",
			timecode:    "10:00:00:00",
			framerate:   29.97,
			expectedIdx: 1080000,
		},
		{
			name:        "24",
			timecode:    "00:00:01:00",
			framerate:   24,
			expectedIdx: 24,
		},
		{
			name:        "24",
			timecode:    "01:00:00:00",
			framerate:   24,
			expectedIdx: 86400,
		},
		{
			name:        "24",
			timecode:    "10:00:00:00",
			framerate:   24,
			expectedIdx: 864000,
		},
		{
			name:        "24",
			timecode:    "20:00:00:00",
			framerate:   24,
			expectedIdx: 1_728_000,
		},
		{
			name:        "24",
			timecode:    "23:00:00:00",
			framerate:   24,
			expectedIdx: 1_728_000,
		},
	}
	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc, err := NewTimecodeFromString(tt.timecode, tt.framerate)
			require.NoError(t, err)
			require.Equal(t, tt.expectedIdx, tc.GetFrameIdx())
		})
	}
}

func TestDfIFrameRates(t *testing.T) {
	tests := []struct {
		name      string
		error     error
		timecode  string
		framerate float64
	}{
		{
			name:      "valid 29.97",
			error:     nil,
			timecode:  "00:00:00;10",
			framerate: 29.97,
		},
		{
			name:      "valid 59.94",
			error:     nil,
			timecode:  "00:00:00;10",
			framerate: 59.94,
		},
		{
			name:      "valid 23.976",
			error:     fmt.Errorf("23.976 is not a valid framerate for drop frame timecode"),
			timecode:  "00:00:00;10",
			framerate: 23.976,
		},
	}

	for _, tt := range tests {
		tc, err := NewTimecodeFromString(tt.timecode, tt.framerate)
		require.Nil(t, err)

		err = tc.Validate()

		if tt.error != nil {
			require.EqualError(t, err, tt.error.Error())
		} else {
			require.Nil(t, err)
		}

	}

}

func TestSeriesTc(t *testing.T) {

	tests := []struct {
		name      string
		timecodes []struct {
			timecode string
			valid    bool
		}
		framerate float64
	}{
		{
			name: "29.97df",
			timecodes: []struct {
				timecode string
				valid    bool
			}{
				{
					timecode: "00:00:00;28",
					valid:    true,
				},
				{
					timecode: "00:00:00;29",
					valid:    true,
				},
				{
					timecode: "00:00:01;00",
					valid:    true,
				},
				{
					timecode: "00:00:01;01",
					valid:    true,
				},
				{
					timecode: "00:01:00;28",
					valid:    true,
				},
				{
					timecode: "00:01:00;29",
					valid:    true,
				},
				{
					timecode: "00:02:00;00",
					valid:    false,
				},
				{
					timecode: "00:02:00;01",
					valid:    false,
				},
				{
					timecode: "00:02:00;02",
					valid:    true,
				},
				{
					timecode: "00:02:00;03",
					valid:    true,
				},
				{
					timecode: "00:09:59;28",
					valid:    true,
				},
				{
					timecode: "00:09:59;29",
					valid:    true,
				},
				{
					timecode: "00:10:00;00",
					valid:    true,
				},
				{
					timecode: "00:10:00;01",
					valid:    true,
				},
				{
					timecode: "00:10:00;02",
					valid:    true,
				},
				{
					timecode: "00:10:00;03",
					valid:    true,
				},
			},
			framerate: 29.97,
		},
		{
			name: "59.94df",
			timecodes: []struct {
				timecode string
				valid    bool
			}{
				{
					timecode: "00:00:00;58",
					valid:    true,
				},
				{
					timecode: "00:00:00;59",
					valid:    true,
				},
				{
					timecode: "00:00:01;00",
					valid:    true,
				},
				{
					timecode: "00:00:01;01",
					valid:    true,
				},
				{
					timecode: "00:01:00;58",
					valid:    true,
				},
				{
					timecode: "00:01:00;59",
					valid:    true,
				},
				{
					timecode: "00:02:00;00",
					valid:    false,
				},
				{
					timecode: "00:02:00;01",
					valid:    false,
				},
				{
					timecode: "00:02:00;02",
					valid:    false,
				},
				{
					timecode: "00:02:00;03",
					valid:    false,
				},
				{
					timecode: "00:02:00;04",
					valid:    true,
				},
				{
					timecode: "00:02:00;05",
					valid:    true,
				},
				{
					timecode: "00:09:59;58",
					valid:    true,
				},
				{
					timecode: "00:09:59;59",
					valid:    true,
				},
				{
					timecode: "00:10:00;00",
					valid:    true,
				},
				{
					timecode: "00:10:00;01",
					valid:    true,
				},
				{
					timecode: "00:10:00;02",
					valid:    true,
				},
				{
					timecode: "00:10:00;03",
					valid:    true,
				},
			},
			framerate: 59.94,
		},
	}
	t.Parallel()
	for _, tt := range tests {

		for _, timecode := range tt.timecodes {
			t.Run(tt.name, func(t *testing.T) {
				tc, err := NewTimecodeFromString(timecode.timecode, tt.framerate)
				require.Nil(t, err)
				err = tc.Validate()
				if timecode.valid {
					require.Nil(t, err)
				} else {
					require.NotNil(t, err)
				}
			})
		}

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
			name:      "Errors when invalid drop frame timecode",
			error:     fmt.Errorf("00:07:00;00 is not valid drop frame timecode"),
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
