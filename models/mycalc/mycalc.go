package mycalc

import "time"

type CalcRequestModel struct {
	Input     []float64 `json:"input" validate:"required"`
	Operation string    `json:"op" validate:"required"`
}

type CalcResponse struct {
	Result []float64 `json:"result"`
	Error  string    `json:"error,omitempty"`
}
type CalcResponseBUlk struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

type CalcLogEntry struct {
	RequestTime  time.Time `db:"request_time"`
	ResponseTime time.Time `db:"response_time"`
	DurationMs   float64   `db:"duration_ms"`
	RequestData  string    `db:"request_data"`
	ResponseData string    `db:"response_data"`
	Error        string    `db:"error"`
	Operation    string    `db:"operation"`
	Input        string    `db:"input"`
}
