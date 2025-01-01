package main

import (
	"fmt"

	"github.com/marcrleonard/TimecodeTool/pkg"
)

func printSeparator() {
	fmt.Println("─────────────────────────────")
}

const title = "🎥 TimecodeTool"

// PrettyPrint displays the timecode validation results in a user-friendly format
func PrettyPrintValidate(r *timecodetool.ValidateResponse) {
	fmt.Println(title + " Validate")
	printSeparator()
	fmt.Printf("Input Timecode:   %s\n", r.InputTimecode)
	fmt.Printf("Frame Rate (FPS): %.2f\n", r.InputFps)

	if r.Valid {
		dfIndicator := ""
		if r.IsDf {
			dfIndicator = " (Drop Frame)"
		}
		fmt.Printf("Valid Timecode:   ✅ Yes%s\n", dfIndicator)
		fmt.Printf("Next Timecode:    %s\n", r.NextTimecode)
	} else {
		fmt.Printf("Valid Timecode:   ❌ No\n")
		fmt.Printf("Error:            %s\n", r.ErrorMsg)
	}

	printSeparator()
}

// prettyPrintSpanResponse displays the span command results in a user-friendly format
func PrettyPrintSpan(r *timecodetool.SpanResponse) {
	fmt.Println(title + " Span")
	printSeparator()

	// Helper function to handle invalid timecodes
	printInvalidTimecode := func(timecode string) string {
		if timecode == "" {
			return "❌ Invalid (Empty)"
		}
		return timecode
	}

	// Print First and Last Timecodes
	fmt.Printf("First Timecode:   %s\n", printInvalidTimecode(r.InputFirstTimecode))
	fmt.Printf("Last Timecode:    %s\n", printInvalidTimecode(r.InputLastTimecode))
	fmt.Printf("Frame Rate (FPS): %.2f\n", r.InputFps)

	// Output based on the validity of the span
	if r.Valid {
		fmt.Printf("Valid Span:       ✅ Yes\n")
		fmt.Printf("Start Frame Index:    %d\n", r.StartFrameIdx)
		fmt.Printf("Last Frame Index:     %d\n", r.LastFrameIdx)
		fmt.Printf("Length (Frames):      %d\n", r.LengthFrames)
		fmt.Printf("Length (Real Time):   %s\n", r.LengthTime)
		fmt.Printf("Length (Seconds):     %.2f\n", r.LengthSeconds)
		fmt.Printf("Length (Timecode):    %s\n", r.LengthTimecode)
		fmt.Printf("Next Timecode:        %s\n", r.NextTimecode)
	} else {
		fmt.Printf("Valid Span:        ❌ No\n")
		fmt.Printf("Error:                %s\n", r.ErrorMsg)
	}

	printSeparator()
}

func PrettyPrintCalc(c *timecodetool.CalcResponse) {

	steps := c.Steps

	fmt.Println(title + " Calculate")
	printSeparator()

	// Starting timecode and frames
	fmt.Printf(" 🎬 Starting Timecode:      %s (Index %d)\n", c.InputFirstTimecode, c.StartFrameIdx)

	// Process each step
	for _, step := range steps {
		if step.Operation == "+" {
			fmt.Printf("   ➕  Add Timecode:         %s (%d frames)\n", step.Timecode, step.Frames)
		} else if step.Operation == "-" {
			fmt.Printf("   ➖  Sub Timecode:         %s (%d frames)\n", step.Timecode, step.Frames)
		}
	}

	// Resulting timecode and frames
	printSeparator()
	fmt.Printf(" 🟰  Resulting Timecode:    %s (%d total frames)\n", c.LastTimecode, c.LengthFrames)
	fmt.Printf("%d ➡️ %d frame indexes\n", c.StartFrameIdx, c.LastFrameIdx)
	printSeparator()
}
