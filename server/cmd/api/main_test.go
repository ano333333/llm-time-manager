package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	datamodel "github.com/ano333333/llm-time-manager/server/internal/data-model"
	"github.com/ano333333/llm-time-manager/server/internal/database"
	"github.com/stretchr/testify/assert"
)

const dbPathKey = "DB_PATH"

var originalDBPath = ""

const dbPath = "../../data/test.db"

func beforeEach() (*sql.DB, error) {
	originalDBPath = os.Getenv(dbPathKey)
	os.Setenv(dbPathKey, dbPath)
	db, err := database.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	migrationsDir := "../../migrations"
	if err = database.RunMigrations(db, migrationsDir); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	return db, nil
}

func afterEach(db *sql.DB) {
	db.Close()
	os.Remove(dbPath)
	os.Setenv(dbPathKey, originalDBPath)
}

func readResponseBody(rec *httptest.ResponseRecorder) (map[string]interface{}, error) {
	body, err := io.ReadAll(rec.Result().Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal body: %w", err)
	}
	return response, nil
}

func insertCaptureSchedules(db *sql.DB, schedules []datamodel.CaptureSchedule) error {
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

func TestGetCaptureScheduleIntegrate(t *testing.T) {
	t.Run("GET /capture/schedule はアクティブなスケジュールがなければnullを返す", func(t *testing.T) {
		// Arrange
		db, err := beforeEach()
		if err != nil {
			t.Fatalf("failed to set up test: %v", err)
		}
		defer afterEach(db)
		mux := setupHandlers(db)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/capture/schedule", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, rec.Code, http.StatusOK)
		response, err := readResponseBody(rec)
		assert.NoError(t, err)
		assert.Nil(t, response["schedule"])
	})

	t.Run("GET /capture/schedule はアクティブなスケジュールが 1 つだけあればそれを返す", func(t *testing.T) {
		// Arrange
		db, err := beforeEach()
		if err != nil {
			t.Fatalf("failed to set up test: %v", err)
		}
		defer afterEach(db)
		mux := setupHandlers(db)
		schedules := []datamodel.CaptureSchedule{
			{
				ID:                "schedule-0",
				Active:            false,
				IntervalMin:       10,
				RetentionMaxItems: 100,
				RetentionMaxDays:  30,
			},
			{
				ID:                "schedule-1",
				Active:            true,
				IntervalMin:       5,
				RetentionMaxItems: 1000,
				RetentionMaxDays:  30,
			},
			{
				ID:                "schedule-2",
				Active:            false,
				IntervalMin:       15,
				RetentionMaxItems: 10000,
				RetentionMaxDays:  30,
			},
		}
		if err := insertCaptureSchedules(db, schedules); err != nil {
			t.Fatalf("failed to insert capture schedules: %v", err)
		}

		// Act
		req := httptest.NewRequest(http.MethodGet, "/capture/schedule", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, rec.Code, http.StatusOK)
		response, err := readResponseBody(rec)
		assert.NoError(t, err)
		assert.NotNil(t, response["schedule"])
		assert.Equal(t, schedules[1].ID, response["schedule"].(map[string]interface{})["id"])
		assert.Equal(t, schedules[1].Active, response["schedule"].(map[string]interface{})["active"])
		assert.Equal(t, schedules[1].IntervalMin, response["schedule"].(map[string]interface{})["interval_min"])
		assert.Equal(t, schedules[1].RetentionMaxItems, response["schedule"].(map[string]interface{})["retention_max_items"])
		assert.Equal(t, schedules[1].RetentionMaxDays, response["schedule"].(map[string]interface{})["retention_max_days"])
	})

	t.Run("GET /capture/schedule は複数のアクティブなスケジュールがあれば 500 Internal Server Error を返す", func(t *testing.T) {
		// Arrange
		db, err := beforeEach()
		if err != nil {
			t.Fatalf("failed to set up test: %v", err)
		}
		defer afterEach(db)
		mux := setupHandlers(db)
		schedules := []datamodel.CaptureSchedule{
			{
				ID:                "schedule-0",
				Active:            true,
				IntervalMin:       10,
				RetentionMaxItems: 100,
				RetentionMaxDays:  30,
			},
			{
				ID:                "schedule-1",
				Active:            true,
				IntervalMin:       5,
				RetentionMaxItems: 1000,
				RetentionMaxDays:  30,
			},
		}
		if err := insertCaptureSchedules(db, schedules); err != nil {
			t.Fatalf("failed to insert capture schedules: %v", err)
		}

		// Act
		req := httptest.NewRequest(http.MethodGet, "/capture/schedule", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, rec.Code, http.StatusInternalServerError)
	})
}
