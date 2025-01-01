package timecodetool

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

func newFailedValidateResponse(InputTimecode string, InputFps float64, IsDf bool, ErrorMsg string) *ValidateResponse {
	return &ValidateResponse{
		InputTimecode: InputTimecode,
		InputFps:      InputFps,
		IsDf:          IsDf,
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
