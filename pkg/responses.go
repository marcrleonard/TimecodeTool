package timecodetool

import "fmt"

type ValidateResponse struct {
	InputTimecode string  `json:"inputTimecode"`
	InputFps      float64 `json:"inputFps"`
	Valid         bool    `json:"valid"`
	ErrorMsg      string  `json:"errorMsg"`
	IsDf          bool    `json:"isDf"`
	NextTimecode  string  `json:"nextTimecode"`
}
type SpanResponse struct {
	InputFirstTimecode  string  `json:"inputFirstTimecode"`
	InputLastTimecode   string  `json:"inputLastTimecode,omitempty"`
	InputFps            float64 `json:"inputFps"`
	Valid               bool    `json:"valid"`
	ErrorMsg            string  `json:"errorMsg"`
	IsDf                bool    `json:"isDf"`
	ExcludeLastTimecode bool    `json:"excludeLastTimecode"`
	StartFrameIdx       int     `json:"startFrameIdx"`
	LastFrameIdx        int     `json:"lastFrameIdx"`
	LengthFrames        int     `json:"lengthFrames"`
	LengthTime          string  `json:"lengthTime"`
	LengthTimecode      string  `json:"lengthTimecode"`
	LengthSeconds       float64 `json:"lengthSeconds"`
	NextTimecode        string  `json:"nextTimecode"`
}

type CalculationStep struct {
	Operation string `json:"operation"` // "+" or "-"
	Timecode  string `json:"timecode"`  // Timecode for the operation
	Frames    int    `json:"frames"`    // Equivalent frames for the operation
}

type CalcResponse struct {
	SpanResponse
	LastTimecode string `json:"lastTimecode"`
	Steps        []CalculationStep
}

// newOkCalcResponse creates a new successful CalcResponse with steps
func newOkCalcResponse(
	InputFirstTimecode string,
	LastTimecode string,
	InputFps float64,
	IsDf bool,
	ExcludeLastTimecode bool,
	StartFrameIdx int,
	LastFrameIdx int,
	LengthFrames int,
	LengthTime string,
	LengthTimecode string,
	LengthSeconds float64,
	NextTimecode string,
	Steps []CalculationStep) *CalcResponse {

	// Create the SpanResponse part of the CalcResponse
	spanResponse := newOkSpanResponse(
		InputFirstTimecode,
		"",
		InputFps,
		IsDf,
		ExcludeLastTimecode,
		StartFrameIdx,
		LastFrameIdx,
		LengthFrames,
		LengthTime,
		LengthTimecode,
		LengthSeconds,
		NextTimecode,
	)

	// Return the CalcResponse with the embedded SpanResponse and the steps
	return &CalcResponse{
		LastTimecode: LastTimecode,
		SpanResponse: *spanResponse, // Unwrap the SpanResponse pointer
		Steps:        Steps,
	}
}

// newFailedCalcResponse creates a new failed CalcResponse with an error message
func newFailedCalcResponse(
	InputFirstTimecode string,
	InputLastTimecode string,
	InputFps float64,
	ExcludeLastTimecode bool,
	ErrorMsg string,
	Steps []CalculationStep) *CalcResponse {

	// Create the failed SpanResponse part
	spanResponse := newFailedSpanResponse(
		InputFirstTimecode,
		InputLastTimecode,
		InputFps,
		ExcludeLastTimecode,
		ErrorMsg,
	)

	// Return the CalcResponse with the embedded SpanResponse and the steps
	return &CalcResponse{
		SpanResponse: *spanResponse, // Unwrap the SpanResponse pointer
		Steps:        Steps,
	}
}

func newOkValidateResponse(InputTimecode string, InputFps float64, IsDf bool, NextTimecode string) *ValidateResponse {
	return &ValidateResponse{
		InputTimecode: InputTimecode,
		InputFps:      InputFps,
		Valid:         true,
		IsDf:          IsDf,
		NextTimecode:  NextTimecode,
	}
}

func newFailedValidateResponse(InputTimecode string, InputFps float64, ErrorMsg string) *ValidateResponse {
	return &ValidateResponse{
		InputTimecode: InputTimecode,
		InputFps:      InputFps,
		Valid:         false,
		ErrorMsg:      ErrorMsg,
	}
}

func newOkSpanResponse(
	InputFirstTimecode string,
	InputLastTimecode string,
	InputFps float64,
	IsDf bool,
	ExcludeLastTimecode bool,
	StartFrameIdx int,
	LastFrameIdx int,
	LengthFrames int,
	LengthTime string,
	LengthTimecode string,
	LengthSeconds float64,
	NextTimecode string) *SpanResponse {
	return &SpanResponse{
		InputFirstTimecode:  InputFirstTimecode,
		InputLastTimecode:   InputLastTimecode,
		InputFps:            InputFps,
		Valid:               true,
		ErrorMsg:            "",
		IsDf:                IsDf,
		ExcludeLastTimecode: ExcludeLastTimecode,
		StartFrameIdx:       StartFrameIdx,
		LastFrameIdx:        LastFrameIdx,
		LengthFrames:        LengthFrames,
		LengthTime:          LengthTime,
		LengthTimecode:      LengthTimecode,
		LengthSeconds:       LengthSeconds,
		NextTimecode:        NextTimecode,
	}
}

func newFailedSpanResponse(
	InputFirstTimecode string,
	InputLastTimecode string,
	InputFps float64,
	ExcludeLastTimecode bool,
	ErrorMsg string) *SpanResponse {
	return &SpanResponse{
		InputFirstTimecode:  InputFirstTimecode,
		InputLastTimecode:   InputLastTimecode,
		InputFps:            InputFps,
		Valid:               false,
		ErrorMsg:            ErrorMsg,
		ExcludeLastTimecode: ExcludeLastTimecode,
	}
}

func printSeparator() {
	fmt.Println("─────────────────────────────")
}

const title = "🎥 TimecodeTool"

// PrettyPrint displays the timecode validation results in a user-friendly format
func (r *ValidateResponse) PrettyPrint() {
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
func (r *SpanResponse) PrettyPrint() {
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

func (c *CalcResponse) PrettyPrint() {

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
