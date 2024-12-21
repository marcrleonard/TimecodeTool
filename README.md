# TimecodeTool

`TimecodeTool` is a simple CLI tool that does a handful of handy functions:
- Validate timecode input
- Calculate the time between two timecodes in realtime, framecounts, and timecode time.
- Timecode calculator where you can add timecodes or frames together.
- (Coming soon) Convert timecodes between different frame rates

## Usage

Modes of operation:
`TimecodeTool --fps=29.97 "00:07:00;00"` - returns: valid, frame start idx, maayyybeee time from 00:00:00:00?
`TimecodeTool 23.98 "01:00:00:00" "01:01:00:00"` - returns valid, first frame idx, last frame idx (assuming inclusive), playback time, time span in frames
`TimecodeTool 23.98 "01:00:00:00" + "00:00:01:00" + 23 - 00:00:00:10` - returns valid, first frame idx, last frame idx (assuming inclusive), time span in frames

## Building

To build the tool, simply run `go build -o TimecodeTool ./cmd/TimecodeTool/main.go` in the root directory of the project. This will create an executable named `main` in the root directory.
