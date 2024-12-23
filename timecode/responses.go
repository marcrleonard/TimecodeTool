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
	InputLastTimecode   string  `json:"inputLastTimecode"`
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
type CalcResponse struct {
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

// prettyPrint displays the timecode validation results in a user-friendly format
func PrettyPrint(response *ValidateResponse) {
	fmt.Println("🎥 Timecode Validation Tool")
	fmt.Println("─────────────────────────────")
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

	fmt.Println("─────────────────────────────")
}

// prettyPrintSpanResponse displays the span command results in a user-friendly format
func PrettyPrintSpanResponse(response SpanResponse) {
	fmt.Println("🎥 Timecode Span Tool")
	fmt.Println("─────────────────────────────")

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
		fmt.Printf("Start Frame Index: %d\n", response.StartFrameIdx)
		fmt.Printf("Last Frame Index:  %d\n", response.LastFrameIdx)
		fmt.Printf("Length (Frames):  %d\n", response.LengthFrames)
		fmt.Printf("Length (Time):    %s\n", response.LengthTime)
		fmt.Printf("Length (Seconds): %.2f\n", response.LengthSeconds)
		fmt.Printf("Next Timecode:    %s\n", response.NextTimecode)
	} else {
		fmt.Printf("Valid Span:       ❌ No\n")
		fmt.Printf("Error:            %s\n", response.ErrorMsg)
	}

	fmt.Println("─────────────────────────────")
}
