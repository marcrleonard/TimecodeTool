# TimecodeTool

`TimecodeTool` is a simple CLI tool that does a handful of handy functions:
- Validate input timecode
- Calculate times spanning two timecodes in playback time, framecounts, and more.
- Timecode calculator where you can add timecodes or frames together.

## Installation

### Use as library

`go get github.com/marcrleonard/TimecodeTool@latest`

### Download binaries

Download the latest from the [releases page](https://github.com/marcrleonard/TimecodeTool/releases).

## Build from source

```shell
git clone https://github.com/marcrleonard/TimecodeTool.git
cd TimecodeTool 
make build
```

## CLI Usage

### Validate
`TimecodeTool validate "00:07:00;00" --fps=29.97`

### Span
`TimecodeTool span "01:00:00:00" "01:01:00:00" --fps=23.98`

### Calculate
`TimecodeTool calculate "01:00:00:00" + "00:00:01:00" + 23 - "00:00:00:10" --fps=23.98`

### JSON Schema outputs
`TimecodeTool schema validate`

## Contributing

### Pull Requests

All dev is done through pull rests on main. They cannot be merged unless they pass the status checks.

### Builds

Builds will only occur if the status check on main completes with a version bump. You can still merge into main without a version bump, but a build may not occur if it is not bumped.

## Todo
- Maybe introduce API in the lib to do NewTimecode and attempt to fix a broken timecode (divmod)
- Request a single value (by key) of json output. Something like `TimecodeTool validate 00:01:00:00 --json-output --key=valid` -> `true` 
- Introduce an API to convert timecodes between different frame rates