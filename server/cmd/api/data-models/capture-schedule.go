package datamodels

import "time"

type CaptureSchedule struct {
	ID                string    `json:"id"`
	Active            bool      `json:"active"`
	IntervalMin       int       `json:"interval_min"`
	RetentionMaxItems int       `json:"retention_max_items"`
	RetentionMaxDays  int       `json:"retention_max_days"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
