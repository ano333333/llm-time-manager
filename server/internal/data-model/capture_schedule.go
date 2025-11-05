package datamodel

import "time"

type CaptureSchedule struct {
	ID          string    `json:"id"`
	Active      bool      `json:"active"`
	IntervalMin int       `json:"interval_min"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
