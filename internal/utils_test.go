package internal

import (
	"testing"
)

func TestParseStringToTimecode(t *testing.T) {
	tests := []struct {
		name                string
		input               string
		fps                 float64
		excludeLastTimecode bool
		dropFrame           bool
		expectedFrames      int
		expectError         bool
	}{
		{"Valid Frame Count", "300", 30.0, false, false, 300, false},
		{"Frame Count With Exclude", "300", 30.0, true, false, 300, false},
		{"Valid Timecode", "00:00:10:00", 30.0, false, false, 301, false},
		{"Valid Timecode", "00:00:10:00", 30.0, true, false, 300, false},
		{"Invalid Timecode", "invalid", 30.0, false, false, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc, err := ParseStringToTimecode(tt.input, tt.fps, tt.excludeLastTimecode, tt.dropFrame)
			if (err != nil) != tt.expectError {
				t.Fatalf("Expected error: %v, got: %v", tt.expectError, err)
			}
			if err == nil && tc.GetFrameCount() != tt.expectedFrames {
				t.Errorf("Expected frames: %d, got: %d", tt.expectedFrames, tc.GetFrameCount())
			}
		})
	}
}

func TestDivmod(t *testing.T) {
	cases := []struct {
		numerator   int64
		denominator int64
		quotient    int64
		remainder   int64
	}{
		{10, 3, 3, 1},
		{100, 7, 14, 2},
		{42, 42, 1, 0},
		{0, 1, 0, 0},
	}

	for _, c := range cases {
		q, r := divmod(c.numerator, c.denominator)
		if q != c.quotient || r != c.remainder {
			t.Errorf("divmod(%d, %d) = (%d, %d); expected (%d, %d)", c.numerator, c.denominator, q, r, c.quotient, c.remainder)
		}
	}
}

func TestFormatTimecode(t *testing.T) {
	cases := []struct {
		hours       int64
		minutes     int64
		seconds     int64
		frames      int64
		isDropframe bool
		expected    string
	}{
		{1, 2, 3, 4, false, "01:02:03:04"},
		{10, 20, 30, 15, true, "10:20:30;15"},
		{0, 0, 0, 0, false, "00:00:00:00"},
	}

	for _, c := range cases {
		result := formatTimecode(c.hours, c.minutes, c.seconds, c.frames, c.isDropframe)
		if result != c.expected {
			t.Errorf("formatTimecode(%d, %d, %d, %d, %v) = %q; expected %q", c.hours, c.minutes, c.seconds, c.frames, c.isDropframe, result, c.expected)
		}
	}
}

func TestFormatTimeSpan(t *testing.T) {
	cases := []struct {
		hours    int64
		minutes  int64
		seconds  int64
		ms       string
		expected string
	}{
		{1, 2, 3, "456", "01:02:03.456"},
		{0, 0, 0, "000", "00:00:00.000"},
	}

	for _, c := range cases {
		result := formatTimeSpan(c.hours, c.minutes, c.seconds, c.ms)
		if result != c.expected {
			t.Errorf("formatTimeSpan(%d, %d, %d, %q) = %q; expected %q", c.hours, c.minutes, c.seconds, c.ms, result, c.expected)
		}
	}
}

func TestGetTimeBase(t *testing.T) {
	cases := []struct {
		framerate float64
		expected  int
	}{
		{29.97, 30},
		{59.94, 60},
		{24.0, 24},
	}

	for _, c := range cases {
		result := getTimeBase(c.framerate)
		if result != c.expected {
			t.Errorf("getTimeBase(%f) = %d; expected %d", c.framerate, result, c.expected)
		}
	}
}
