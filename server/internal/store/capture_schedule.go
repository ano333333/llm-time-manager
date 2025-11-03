package store

import (
	"database/sql"
	"errors"

	datamodel "github.com/ano333333/llm-time-manager/server/internal/data-model"
)

type CaptureScheduleStore interface {
	GetActiveCaptureSchedule() (*datamodel.CaptureSchedule, error)
}

type DefaultCaptureScheduleStore struct {
	DB *sql.DB
}

func (s *DefaultCaptureScheduleStore) GetActiveCaptureSchedule() (*datamodel.CaptureSchedule, error) {
	var row datamodel.CaptureSchedule

	rows, err := s.DB.Query("SELECT id, active, interval_min, retention_max_items, retention_max_days, created_at, updated_at FROM capture_schedules WHERE active = 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		if err = rows.Err(); err != nil {
			return nil, err
		}
		return nil, nil
	}
	err = rows.Scan(&row.ID, &row.Active, &row.IntervalMin, &row.RetentionMaxItems, &row.RetentionMaxDays, &row.CreatedAt, &row.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return nil, errors.New("multiple active capture schedules found")
	}

	return &row, nil
}
