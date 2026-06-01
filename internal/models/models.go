package models

import "time"

type Status string

const (
	StatusUp    Status = "UP"
	StatusDown  Status = "DOWN"
	StatusError Status = "ERROR"
)

type CheckResult struct {
	URL            string        `json:"url"`
	FinalURL       string        `json:"finalUrl"`
	Status         Status        `json:"status"`
	StatusCode     int           `json:"statusCode"`
	ResponseTime   time.Duration `json:"-"`
	ResponseTimeMs int64         `json:"responseTimeMs"`
	ErrorType      string        `json:"errorType,omitempty"`
	ErrorMessage   string        `json:"errorMessage,omitempty"`
}

type Summary struct {
	Total       int           `json:"total"`
	Success     int           `json:"success"`
	Failure     int           `json:"failure"`
	ElapsedTime time.Duration `json:"elapsedTime"`
}
