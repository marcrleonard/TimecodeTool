# TimecodeTool

`TimecodeTool` is a simple CLI tool that does a handful of handy functions:
- Validate input timecode
- Calculate times spanning two timecodes in playback time, framecounts, timecode time, and more.
- Timecode calculator where you can add timecodes or frames together.

## Usage

### Validate
`TimecodeTool validate "00:07:00;00" --fps=29.97`

### Span
`TimecodeTool span "01:00:00:00" "01:01:00:00" --fps=23.98`

### Calculate
`TimecodeTool calculate "01:00:00:00" + "00:00:01:00" + 23 - "00:00:00:10" --fps=23.98`

### JSON Schema outputs
`TimecodeTool schema validate`

## Development

### CLI
To build the tool, simply run `go build -o dist/TimecodeTool ./cmd/TimecodeTool/main.go` in the root directory of the project. This will create the executable of `dist/TimecodeTool`

### Wasm
`make build_wasm`

## Todo
- Introduce an API to convert timecodes between different frame rates