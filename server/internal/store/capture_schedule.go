package store

import (
	"database/sql"
	"errors"

	datamodel "github.com/ano333333/llm-time-manager/server/internal/data-model"
)

type CaptureScheduleStore interface {
	GetActiveCaptureSchedule(tx Transaction) (*datamodel.CaptureSchedule, error)
	// アクティブなスケジュールを更新し、更新された行数を返す。
	//
	// 更新前のアクティブなスケジュールの数をチェックしない点に注意。
	UpdateActiveCaptureSchedule(tx Transaction, active bool, intervalMin int) (int64, error)
}

type DefaultCaptureScheduleStore struct {
	DB *sql.DB
}

func (s *DefaultCaptureScheduleStore) GetActiveCaptureSchedule(tx Transaction) (*datamodel.CaptureSchedule, error) {
	defaultTx, ok := tx.(DefaultTransaction)
	if !ok {
		return nil, errors.New("transaction is not DefaultTransaction")
	}
	var row datamodel.CaptureSchedule

	rows, err := defaultTx.Tx.Query("SELECT id, active, interval_min, created_at, updated_at FROM capture_schedules WHERE active = 1")
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
	err = rows.Scan(&row.ID, &row.Active, &row.IntervalMin, &row.CreatedAt, &row.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return nil, errors.New("multiple active capture schedules found")
	}
	return &row, nil
}

func (s *DefaultCaptureScheduleStore) UpdateActiveCaptureSchedule(tx Transaction, active bool, intervalMin int) (int64, error) {
	defaultTx, ok := tx.(DefaultTransaction)
	if !ok {
		return 0, errors.New("transaction is not DefaultTransaction")
	}
	result, err := defaultTx.Tx.Exec(
		"UPDATE capture_schedules SET active = ?, interval_min = ? WHERE active = 1",
		active,
		intervalMin,
	)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}
