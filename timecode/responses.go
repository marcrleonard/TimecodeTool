package timecode

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
