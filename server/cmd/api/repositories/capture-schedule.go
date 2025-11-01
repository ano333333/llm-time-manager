package repositories

import (
	"database/sql"
	"errors"

	datamodels "github.com/ano333333/llm-time-manager/server/cmd/api/data-models"
)

type CaptureScheduleRepository interface {
	GetActiveCaptureSchedule() (*datamodels.CaptureSchedule, error)
}

type DefaultCaptureScheduleRepository struct {
	DB *sql.DB
}

func (r *DefaultCaptureScheduleRepository) GetActiveCaptureSchedule() (*datamodels.CaptureSchedule, error) {
	var row datamodels.CaptureSchedule

	rows, err := r.DB.Query("SELECT * FROM capture_schedules WHERE active = 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
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
