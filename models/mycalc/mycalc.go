package mycalc

import "time"

type CalcRequestModel struct {
	Input     []float64 `json:"input" validate:"required"`
	Operation string    `json:"op" validate:"required"`
}

type CalcResponseModel struct {
	ID           uint      `gorm:"primaryKey"`
	Input        string    `gorm:"type:text"` // JSON input
	Operation    string    `gorm:"type:varchar(50)"`
	Result       string    `gorm:"type:text"` // JSON result
	Error        string    `gorm:"type:text"`
	RequestTime  time.Time `gorm:"column:request_time"`
	ResponseTime time.Time `gorm:"column:response_time"`
	DurationMs   float64   `gorm:"column:duration_ms"`
}

func (CalcResponseModel) TableName() string {
	return "my_calculation_logs"
}

type CalcResponse struct {
	Result []float64 `json:"result"`
	Error  string    `json:"error,omitempty"`
}
