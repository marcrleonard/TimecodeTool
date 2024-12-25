package internal

import (
	"testing"
)

func TestNewTimecodeSpan(t *testing.T) {
	start := &Timecode{FrameRate: 24.0, DropFrame: false, _frames: 0}
	end := &Timecode{FrameRate: 24.0, DropFrame: false, _frames: 240}

	span, err := NewTimecodeSpan(start, end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if span.StartTimecode != start {
		t.Errorf("expected StartTimecode to be %v, got %v", start, span.StartTimecode)
	}
	if span.LastTimecode != end {
		t.Errorf("expected LastTimecode to be %v, got %v", end, span.LastTimecode)
	}
	if span.Framerate != 24.0 {
		t.Errorf("expected Framerate to be 24.0, got %v", span.Framerate)
	}
	if span.Dropframe {
		t.Errorf("expected Dropframe to be false, got true")
	}
}

func TestGetTotalSeconds(t *testing.T) {
	start := &Timecode{FrameRate: 24.0, DropFrame: false, _frames: 0}
	end := &Timecode{FrameRate: 24.0, DropFrame: false, _frames: 240}
	span, _ := NewTimecodeSpan(start, end)

	expected := 10.041666666666666
	if got := span.GetTotalSeconds(); got != expected {
		t.Errorf("expected %v seconds, got %v", expected, got)
	}
}

func TestGetTotalFrames(t *testing.T) {
	start := &Timecode{FrameRate: 24.0, DropFrame: false, _frames: 0}
	end := &Timecode{FrameRate: 24.0, DropFrame: false, _frames: 240}
	span, _ := NewTimecodeSpan(start, end)

	expected := 241 // inclusive of start and end frames
	if got := span.GetTotalFrames(); got != expected {
		t.Errorf("expected %v frames, got %v", expected, got)
	}
}

func TestGetSpanTimecode(t *testing.T) {
	start := &Timecode{FrameRate: 24.0, DropFrame: false, _frames: 0}
	end := &Timecode{FrameRate: 24.0, DropFrame: false, _frames: 240}
	span, _ := NewTimecodeSpan(start, end)

	expected := "00:00:10:01" // 241 frames at 24 fps
	if got := span.GetSpanTimecode(); got != expected {
		t.Errorf("expected timecode %v, got %v", expected, got)
	}
}

func TestGetSpanRealtime(t *testing.T) {
	start := &Timecode{FrameRate: 24.0, DropFrame: false, _frames: 0}
	end := &Timecode{FrameRate: 24.0, DropFrame: false, _frames: 240}
	span, _ := NewTimecodeSpan(start, end)

	expected := "00:00:10.042" // 10 seconds in real time
	if got := span.GetSpanRealtime(); got != expected {
		t.Errorf("expected realtime %v, got %v", expected, got)
	}
}
