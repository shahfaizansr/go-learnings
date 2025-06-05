package mybulkcalc

type BulkCalcRequestModel struct {
	Input     string `json:"input" validate:"required"`
	Operation string `json:"op" validate:"required"`
}

type BatchBulkCalcRequestModel struct {
	Input     []float64 `json:"input" validate:"required"`
	Operation string    `json:"op" validate:"required"`
}

type BulkCalcResponseModel struct {
	NLines       int         `json:"nlines"`
	NSuccess     int         `json:"nsuccess"`
	NFailed      int         `json:"nfailed"`
	Results      [][]float64 `json:"results"`
	SuccessLines []int       `json:"successLines,omitempty"`
	FailedLines  []int       `json:"failedLines,omitempty"`
	Lines        []string    `json:"lines,omitempty"`
}
