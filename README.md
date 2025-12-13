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
Validates a timecode and returns information about the frame.

```bash
TimecodeTool validate "00:07:00;00" --fps=29.97
```

### Span
Calculate the duration between two timecodes.

```bash
TimecodeTool span "01:00:00:00" "01:01:00:00" --fps=23.98
```

**Important:** See [Inclusive vs Exclusive Spans](#inclusive-vs-exclusive-spans) below for details on the `-e` flag.

### Calculate
Perform arithmetic operations on timecodes and frame numbers.

```bash
TimecodeTool calculate "01:00:00:00" + "00:00:01:00" + 23 - "00:00:00:10" --fps=23.98
```

### JSON Schema outputs
Get JSON schema definitions for the output formats.

```bash
TimecodeTool schema validate
TimecodeTool schema span
TimecodeTool schema calculate
```

## Inclusive vs Exclusive Spans

**IMPORTANT:** Understanding inclusive vs exclusive behavior is critical for getting expected results.

### Default Behavior (Inclusive)

By default, TimecodeTool treats spans as **INCLUSIVE of both endpoints**. This means both the start and end timecodes are included in the calculation.

**Example at 24fps:**
```bash
$ TimecodeTool span "00:00:00:00" "00:00:01:00" --fps 24

Length (Frames):      25        # Includes frames 0, 1, 2, ..., 24
Length (Seconds):     1.04      # 25 frames / 24fps = 1.041666...
Length (Timecode):    00:00:01:01  # Represents 25 frames
```

This might be surprising! A span from `00:00:00:00` to `00:00:01:00` gives you **25 frames**, not 24.

### Using the `-e` Flag (Exclusive - Industry Standard)

The `-e` or `--exclude-last-timecode` flag makes the **end timecode exclusive**, which matches the behavior of most video editing software (Adobe Premiere, Final Cut Pro, DaVinci Resolve, etc.).

**Example at 24fps:**
```bash
$ TimecodeTool span "00:00:00:00" "00:00:01:00" --fps 24 -e

Length (Frames):      24        # Includes frames 0-23, excludes frame 24
Length (Seconds):     1.00      # 24 frames / 24fps = exactly 1.0
Length (Timecode):    00:00:01:00  # Represents 24 frames (1 second)
```

With the `-e` flag, you get the **expected** result: exactly 1 second = 24 frames.

### When to Use `-e`

**Use `-e` when:**
- You want spans to match industry-standard video editing software
- You think of timecodes as "up to but not including" the end point
- You want intuitive results (e.g., 1 second = 24 frames at 24fps)

**Don't use `-e` when:**
- You specifically need inclusive behavior
- You're working with existing scripts that depend on inclusive behavior
- You want to include both exact endpoints in the calculation

### Visual Example

At 24fps, frames are numbered 0, 1, 2, 3, ..., 23, 24, ...

**Without `-e` (inclusive):**
```
Span: "00:00:00:00" to "00:00:01:00"
      [0] [1] [2] ... [23] [24]
       ^                    ^
     start                 end

Result: 25 frames (both endpoints included)
```

**With `-e` (exclusive):**
```
Span: "00:00:00:00" to "00:00:01:00"
      [0] [1] [2] ... [23] | [24]
       ^                 ^    ^
     start              end   excluded

Result: 24 frames (end excluded, like a half-open interval [start, end))
```

### Calculate Command with `-e`

The `-e` flag also affects the `calculate` command when using timecode strings (not frame numbers):

**Without `-e`:**
```bash
$ TimecodeTool calculate "01:00:00:00" + "00:00:01:00" --fps 24
# Adds 24 frames (the frame index of 00:00:01:00)
```

**With `-e`:**
```bash
$ TimecodeTool calculate "01:00:00:00" + "00:00:01:00" --fps 24 -e
# Adds 23 frames (00:00:01:00 minus one frame)
```

**Note:** Frame numbers (like `+ 24`) are unaffected by the `-e` flag.

## Contributing

### Pull Requests

All dev is done through pull rests on main. They cannot be merged unless they pass the status checks.

### Builds

Builds will only occur if the status check on main completes with a version bump. You can still merge into main without a version bump, but a build may not occur if it is not bumped.

## Todo
- Within `validate` 
  - include frame index from 00:00:00:00
- Maybe introduce API in the lib to do NewTimecode and attempt to fix a broken timecode (divmod)
- Introduce an API to convert timecodes between different frame rates