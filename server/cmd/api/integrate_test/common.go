package integratetest

import (
	"database/sql"
	"fmt"
	"io"
	"net/http/httptest"
	"os"

	datamodel "github.com/ano333333/llm-time-manager/server/internal/data-model"
	"github.com/ano333333/llm-time-manager/server/internal/database"
)

const DBPathKey = "DB_PATH"

var originalDBPath = ""

const (
	dbPath        = "../../../data/test.db"
	migrationsDir = "../../../migrations"
)

func BeforeEach() (*sql.DB, error) {
	originalDBPath = os.Getenv(DBPathKey)
	os.Setenv(DBPathKey, dbPath)
	db, err := database.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err = database.RunMigrations(db, migrationsDir); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	return db, nil
}

func AfterEach(db *sql.DB) {
	db.Close()
	os.Remove(dbPath)
	os.Setenv(DBPathKey, originalDBPath)
}

func GetResponseBodyJson(rec *httptest.ResponseRecorder) (string, error) {
	body, err := io.ReadAll(rec.Result().Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %w", err)
	}
	return string(body), nil
}

func InsertCaptureSchedules(db *sql.DB, schedules []datamodel.CaptureSchedule) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	for _, schedule := range schedules {
		_, err := tx.Exec("INSERT INTO capture_schedules (id, active, interval_min, retention_max_items, retention_max_days, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)", schedule.ID, schedule.Active, schedule.IntervalMin, schedule.RetentionMaxItems, schedule.RetentionMaxDays, schedule.CreatedAt, schedule.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to insert capture schedule: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func InsertGoals(db *sql.DB, goals []datamodel.Goal) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	for _, goal := range goals {
		if goal.KpiName == nil {
			_, err := tx.Exec("INSERT INTO goals (id, status, title, description, start_date, end_date, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", goal.ID, goal.Status, goal.Title, goal.Description, goal.StartDate, goal.EndDate, goal.CreatedAt, goal.UpdatedAt)
			if err != nil {
				return fmt.Errorf("failed to insert goal: %w", err)
			}
			continue
		}
		_, err := tx.Exec("INSERT INTO goals (id, status, title, description, start_date, end_date, kpi_name, kpi_target, kpi_unit, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", goal.ID, goal.Status, goal.Title, goal.Description, goal.StartDate, goal.EndDate, *goal.KpiName, *goal.KpiTarget, *goal.KpiUnit, goal.CreatedAt, goal.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to insert goal: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
