# TimecodeTool

`TimecodeTool` is a simple CLI tool that does a handful of handy functions:
- Validate timecode input
- Calculate the time between two timecodes in realtime, framecounts, and timecode time.
- Timecode calculator where you can add timecodes or frames together.
- (Coming soon) Convert timecodes between different frame rates

## Usage

### Validate
`TimecodeTool validate "00:07:00;00" --fps=29.97` - returns: valid, frame start idx, maayyybeee time from 00:00:00:00?

### Span
`TimecodeTool span "01:00:00:00" "01:01:00:00" --fps=23.98 ` - returns valid, first frame idx, last frame idx (assuming inclusive), playback time, time span in frames

### Calculate
`TimecodeTool calculate "01:00:00:00" + "00:00:01:00" + 23 - "00:00:00:10" --fps=23.98` - returns valid, first frame idx, last frame idx (assuming inclusive), time span in frames

## Development

To build the tool, simply run `go build -o dist/TimecodeTool ./cmd/TimecodeTool/main.go` in the root directory of the project. This will create the executable of `dist/TimecodeTool`

## Todo

- Implement backwards 24 hour rollover for DF
- Create Pretty Text output
- Create Json Output