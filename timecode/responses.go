package timecode

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

// NewOkCalcResponse creates a new successful CalcResponse with steps
func NewOkCalcResponse(
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
	spanResponse := NewOkSpanResponse(
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

// NewFailedCalcResponse creates a new failed CalcResponse with an error message
func NewFailedCalcResponse(
	InputFirstTimecode string,
	InputLastTimecode string,
	InputFps float64,
	ExcludeLastTimecode bool,
	ErrorMsg string,
	Steps []CalculationStep) *CalcResponse {

	// Create the failed SpanResponse part
	spanResponse := NewFailedSpanResponse(
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

func NewOkValidateResponse(InputTimecode string, InputFps float64, IsDf bool, NextTimecode string) *ValidateResponse {
	return &ValidateResponse{
		InputTimecode: InputTimecode,
		InputFps:      InputFps,
		Valid:         true,
		IsDf:          IsDf,
		NextTimecode:  NextTimecode,
	}
}

func NewFailedValidateResponse(InputTimecode string, InputFps float64, ErrorMsg string) *ValidateResponse {
	return &ValidateResponse{
		InputTimecode: InputTimecode,
		InputFps:      InputFps,
		Valid:         false,
		ErrorMsg:      ErrorMsg,
	}
}

func NewOkSpanResponse(
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

func NewFailedSpanResponse(
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

// prettyPrint displays the timecode validation results in a user-friendly format
func PrettyPrint(response *ValidateResponse) {
	fmt.Println(title + " Validate")
	printSeparator()
	fmt.Printf("Input Timecode:   %s\n", response.InputTimecode)
	fmt.Printf("Frame Rate (FPS): %.2f\n", response.InputFps)

	if response.Valid {
		dfIndicator := ""
		if response.IsDf {
			dfIndicator = " (Drop Frame)"
		}
		fmt.Printf("Valid Timecode:   ✅ Yes%s\n", dfIndicator)
		fmt.Printf("Next Timecode:    %s\n", response.NextTimecode)
	} else {
		fmt.Printf("Valid Timecode:   ❌ No\n")
		fmt.Printf("Error:            %s\n", response.ErrorMsg)
	}

	printSeparator()
}

// prettyPrintSpanResponse displays the span command results in a user-friendly format
func PrettyPrintSpanResponse(response SpanResponse) {
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
	fmt.Printf("First Timecode:   %s\n", printInvalidTimecode(response.InputFirstTimecode))
	fmt.Printf("Last Timecode:    %s\n", printInvalidTimecode(response.InputLastTimecode))
	fmt.Printf("Frame Rate (FPS): %.2f\n", response.InputFps)

	// Output based on the validity of the span
	if response.Valid {
		fmt.Printf("Valid Span:       ✅ Yes\n")
		fmt.Printf("Start Frame Index:    %d\n", response.StartFrameIdx)
		fmt.Printf("Last Frame Index:     %d\n", response.LastFrameIdx)
		fmt.Printf("Length (Frames):      %d\n", response.LengthFrames)
		fmt.Printf("Length (Real Time):   %s\n", response.LengthTime)
		fmt.Printf("Length (Seconds):     %.2f\n", response.LengthSeconds)
		fmt.Printf("Length (Timecode):    %s\n", response.LengthTimecode)
		fmt.Printf("Next Timecode:        %s\n", response.NextTimecode)
	} else {
		fmt.Printf("Valid Span:        ❌ No\n")
		fmt.Printf("Error:                %s\n", response.ErrorMsg)
	}

	printSeparator()
}

func PrettyPrintCalculateResponse(resp CalcResponse) {

	steps := resp.Steps

	fmt.Println(title + " Calculate")
	printSeparator()

	// Starting timecode and frames
	fmt.Printf(" 🎬 Starting Timecode:   %s (Index %d)\n", resp.InputFirstTimecode, resp.StartFrameIdx)

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
	fmt.Printf(" ✅  Resulting Timecode:     %s (%d total frames)\n", resp.LastTimecode, resp.LengthFrames)

	// Summary
	printSeparator()
	fmt.Println("Details:")
	fmt.Printf(" 🟢 Start Index:      %d\n", resp.StartFrameIdx)
	for _, step := range steps {
		if step.Operation == "+" {
			fmt.Printf(" ➕ %d\n", step.Frames)
		} else if step.Operation == "-" {
			fmt.Printf(" ➖ %d\n", step.Frames)
		}
	}
	fmt.Printf(" 🔵 Last Index:     %d\n", resp.LastFrameIdx)
	printSeparator()
}
