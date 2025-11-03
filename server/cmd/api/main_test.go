package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

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

func getResponseBodyJson(rec *httptest.ResponseRecorder) (string, error) {
	body, err := io.ReadAll(rec.Result().Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %w", err)
	}
	return string(body), nil
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

func insertGoals(db *sql.DB, goals []datamodel.Goal) error {
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
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err := getResponseBodyJson(rec)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"schedule": null}`, response)
	})

	t.Run("GET /capture/schedule はアクティブなスケジュールが 1 つだけあればそれを返す", func(t *testing.T) {
		// Arrange
		db, err := beforeEach()
		if err != nil {
			t.Fatalf("failed to set up test: %v", err)
		}
		defer afterEach(db)
		mux := setupHandlers(db)
		now := time.Now()
		schedules := []datamodel.CaptureSchedule{
			{
				ID:                "schedule-0",
				Active:            false,
				IntervalMin:       10,
				RetentionMaxItems: 100,
				RetentionMaxDays:  30,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			{
				ID:                "schedule-1",
				Active:            true,
				IntervalMin:       5,
				RetentionMaxItems: 1000,
				RetentionMaxDays:  30,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			{
				ID:                "schedule-2",
				Active:            false,
				IntervalMin:       15,
				RetentionMaxItems: 10000,
				RetentionMaxDays:  30,
				CreatedAt:         now,
				UpdatedAt:         now,
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
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err := getResponseBodyJson(rec)
		assert.NoError(t, err)
		expected, err := json.Marshal(map[string]interface{}{
			"schedule": map[string]interface{}{
				"id":                  schedules[1].ID,
				"active":              schedules[1].Active,
				"interval_min":        schedules[1].IntervalMin,
				"retention_max_items": schedules[1].RetentionMaxItems,
				"retention_max_days":  schedules[1].RetentionMaxDays,
			},
		})
		if err != nil {
			t.Fatalf("failed to marshal expected: %v", err)
		}
		assert.JSONEq(t, string(expected), response)
	})

	t.Run("GET /capture/schedule は複数のアクティブなスケジュールがあれば 500 Internal Server Error を返す", func(t *testing.T) {
		// Arrange
		db, err := beforeEach()
		if err != nil {
			t.Fatalf("failed to set up test: %v", err)
		}
		defer afterEach(db)
		mux := setupHandlers(db)
		now := time.Now()
		schedules := []datamodel.CaptureSchedule{
			{
				ID:                "schedule-0",
				Active:            true,
				IntervalMin:       10,
				RetentionMaxItems: 100,
				RetentionMaxDays:  30,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			{
				ID:                "schedule-1",
				Active:            true,
				IntervalMin:       5,
				RetentionMaxItems: 1000,
				RetentionMaxDays:  30,
				CreatedAt:         now,
				UpdatedAt:         now,
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
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestGetGoalsIntegrate(t *testing.T) {
	t.Run("GET /goal はquery parameterがない場合空配列を返す", func(t *testing.T) {
		// Arrange
		db, err := beforeEach()
		if err != nil {
			t.Fatalf("failed to set up test: %v", err)
		}
		defer afterEach(db)
		mux := setupHandlers(db)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/goal", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err := getResponseBodyJson(rec)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"goals": []
		}`, response)
	})

	t.Run("GET /goal は不正なquery parameterで400を返す", func(t *testing.T) {
		// Arrange
		db, err := beforeEach()
		if err != nil {
			t.Fatalf("failed to set up test: %v", err)
		}
		defer afterEach(db)
		mux := setupHandlers(db)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/goal?status=invalid", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("GET /goal はquery parameterにstatusがマッチするモデルを返す", func(t *testing.T) {
		// Arrange
		db, err := beforeEach()
		if err != nil {
			t.Fatalf("failed to set up test: %v", err)
		}
		defer afterEach(db)
		mux := setupHandlers(db)
		createdAt := time.Date(2025, 10, 1, 0, 0, 0, 0, time.Local)
		updatedAt := time.Date(2025, 10, 2, 0, 0, 0, 0, time.Local)
		startDate := time.Date(2025, 10, 2, 0, 0, 0, 0, time.Local)
		endDate := time.Date(2025, 11, 2, 0, 0, 0, 0, time.Local)
		kpiName0 := "Kpi Name 0"
		kpiTarget0 := 100.0
		kpiUnit0 := "Kpi Unit 0"
		kpiName2 := "Kpi Name 2"
		kpiTarget2 := 100.0
		kpiUnit2 := "Kpi Unit 2"
		goals := []datamodel.Goal{
			{
				ID:          "goal-0",
				Title:       "Goal 0",
				Description: "Description 0",
				StartDate:   startDate,
				EndDate:     endDate,
				KpiName:     &kpiName0,
				KpiTarget:   &kpiTarget0,
				KpiUnit:     &kpiUnit0,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
				Status:      "active",
			},
			{
				ID:          "goal-1",
				Title:       "Goal 1",
				Description: "Description 1",
				StartDate:   startDate,
				EndDate:     endDate,
				KpiName:     nil,
				KpiTarget:   nil,
				KpiUnit:     nil,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
				Status:      "paused",
			},
			{
				ID:          "goal-2",
				Title:       "Goal 2",
				Description: "Description 2",
				StartDate:   startDate,
				EndDate:     endDate,
				KpiName:     &kpiName2,
				KpiTarget:   &kpiTarget2,
				KpiUnit:     &kpiUnit2,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
				Status:      "done",
			},
		}
		if err := insertGoals(db, goals); err != nil {
			t.Fatalf("failed to insert goals: %v", err)
		}

		// Act(none)
		req := httptest.NewRequest(http.MethodGet, "/goal?status=", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert(none)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err := getResponseBodyJson(rec)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"goals": []
		}`, response)

		// Act(active)
		req = httptest.NewRequest(http.MethodGet, "/goal?status=active", nil)
		rec = httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert(active)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err = getResponseBodyJson(rec)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"goals": [
				{
					"id": "goal-0",
					"title": "Goal 0",
					"description": "Description 0",
					"start_date": "2025-10-02",
					"end_date": "2025-11-02",
					"kpi_name": "Kpi Name 0",
					"kpi_target": 100,
					"kpi_unit": "Kpi Unit 0",
					"status": "active",
					"created_at": "2025-10-01T00:00:00+09:00",
					"updated_at": "2025-10-02T00:00:00+09:00"
				}
			]
		}`, response)

		// Act(paused)
		req = httptest.NewRequest(http.MethodGet, "/goal?status=paused", nil)
		rec = httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert(paused)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err = getResponseBodyJson(rec)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"goals": [
				{
					"id": "goal-1",
					"title": "Goal 1",
					"description": "Description 1",
					"start_date": "2025-10-02",
					"end_date": "2025-11-02",
					"kpi_name": null,
					"kpi_target": null,
					"kpi_unit": null,
					"status": "paused",
					"created_at": "2025-10-01T00:00:00+09:00",
					"updated_at": "2025-10-02T00:00:00+09:00"
				}
			]
		}`, response)

		// Act(done)
		req = httptest.NewRequest(http.MethodGet, "/goal?status=done", nil)
		rec = httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert(done)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err = getResponseBodyJson(rec)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"goals": [
				{
					"id": "goal-2",
					"title": "Goal 2",
					"description": "Description 2",
					"start_date": "2025-10-02",
					"end_date": "2025-11-02",
					"kpi_name": "Kpi Name 2",
					"kpi_target": 100,
					"kpi_unit": "Kpi Unit 2",
					"status": "done",
					"created_at": "2025-10-01T00:00:00+09:00",
					"updated_at": "2025-10-02T00:00:00+09:00"
				}
			]
		}`, response)

		// Act(active,paused)
		req = httptest.NewRequest(http.MethodGet, "/goal?status=active,paused", nil)
		rec = httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert(active,paused)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err = getResponseBodyJson(rec)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"goals": [
				{
					"id": "goal-0",
					"title": "Goal 0",
					"description": "Description 0",
					"start_date": "2025-10-02",
					"end_date": "2025-11-02",
					"kpi_name": "Kpi Name 0",
					"kpi_target": 100,
					"kpi_unit": "Kpi Unit 0",
					"status": "active",
					"created_at": "2025-10-01T00:00:00+09:00",
					"updated_at": "2025-10-02T00:00:00+09:00"
				},
				{
					"id": "goal-1",
					"title": "Goal 1",
					"description": "Description 1",
					"start_date": "2025-10-02",
					"end_date": "2025-11-02",
					"kpi_name": null,
					"kpi_target": null,
					"kpi_unit": null,
					"status": "paused",
					"created_at": "2025-10-01T00:00:00+09:00",
					"updated_at": "2025-10-02T00:00:00+09:00"
				}
			]
		}`, response)

		// Act(active,done)
		req = httptest.NewRequest(http.MethodGet, "/goal?status=active,done", nil)
		rec = httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert(active,done)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err = getResponseBodyJson(rec)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"goals": [
				{
					"id": "goal-0",
					"title": "Goal 0",
					"description": "Description 0",
					"start_date": "2025-10-02",
					"end_date": "2025-11-02",
					"kpi_name": "Kpi Name 0",
					"kpi_target": 100,
					"kpi_unit": "Kpi Unit 0",
					"status": "active",
					"created_at": "2025-10-01T00:00:00+09:00",
					"updated_at": "2025-10-02T00:00:00+09:00"
				},
				{
					"id": "goal-2",
					"title": "Goal 2",
					"description": "Description 2",
					"start_date": "2025-10-02",
					"end_date": "2025-11-02",
					"kpi_name": "Kpi Name 2",
					"kpi_target": 100,
					"kpi_unit": "Kpi Unit 2",
					"status": "done",
					"created_at": "2025-10-01T00:00:00+09:00",
					"updated_at": "2025-10-02T00:00:00+09:00"
				}
			]
		}`, response)

		// Act(paused,done)
		req = httptest.NewRequest(http.MethodGet, "/goal?status=paused,done", nil)
		rec = httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert(paused,done)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err = getResponseBodyJson(rec)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"goals": [
				{
					"id": "goal-1",	
					"title": "Goal 1",
					"description": "Description 1",
					"start_date": "2025-10-02",
					"end_date": "2025-11-02",
					"kpi_name": null,
					"kpi_target": null,
					"kpi_unit": null,
					"status": "paused",
					"created_at": "2025-10-01T00:00:00+09:00",
					"updated_at": "2025-10-02T00:00:00+09:00"
				},
				{
					"id": "goal-2",
					"title": "Goal 2",
					"description": "Description 2",
					"start_date": "2025-10-02",
					"end_date": "2025-11-02",
					"kpi_name": "Kpi Name 2",
					"kpi_target": 100,
					"kpi_unit": "Kpi Unit 2",
					"status": "done",
					"created_at": "2025-10-01T00:00:00+09:00",
					"updated_at": "2025-10-02T00:00:00+09:00"
				}
			]
		}`, response)

		// Act(active,done)
		req = httptest.NewRequest(http.MethodGet, "/goal?status=active,done", nil)
		rec = httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert(active,done)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err = getResponseBodyJson(rec)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"goals": [
				{
					"id": "goal-0",
					"title": "Goal 0",
					"description": "Description 0",
					"start_date": "2025-10-02",
					"end_date": "2025-11-02",
					"kpi_name": "Kpi Name 0",
					"kpi_target": 100,
					"kpi_unit": "Kpi Unit 0",
					"status": "active",
					"created_at": "2025-10-01T00:00:00+09:00",
					"updated_at": "2025-10-02T00:00:00+09:00"
				},
				{
					"id": "goal-2",
					"title": "Goal 2",
					"description": "Description 2",
					"start_date": "2025-10-02",
					"end_date": "2025-11-02",
					"kpi_name": "Kpi Name 2",
					"kpi_target": 100,
					"kpi_unit": "Kpi Unit 2",
					"status": "done",
					"created_at": "2025-10-01T00:00:00+09:00",
					"updated_at": "2025-10-02T00:00:00+09:00"
				}
			]
		}`, response)

		// Act(active,paused,done)
		req = httptest.NewRequest(http.MethodGet, "/goal?status=active,paused,done", nil)
		rec = httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert(active,paused,done)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err = getResponseBodyJson(rec)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"goals": [
				{
					"id": "goal-0",
					"title": "Goal 0",
					"description": "Description 0",
					"start_date": "2025-10-02",
					"end_date": "2025-11-02",
					"kpi_name": "Kpi Name 0",
					"kpi_target": 100,
					"kpi_unit": "Kpi Unit 0",
					"status": "active",
					"created_at": "2025-10-01T00:00:00+09:00",
					"updated_at": "2025-10-02T00:00:00+09:00"
				},
				{
					"id": "goal-1",
					"title": "Goal 1",
					"description": "Description 1",
					"start_date": "2025-10-02",
					"end_date": "2025-11-02",
					"kpi_name": null,
					"kpi_target": null,
					"kpi_unit": null,
					"status": "paused",
					"created_at": "2025-10-01T00:00:00+09:00",
					"updated_at": "2025-10-02T00:00:00+09:00"
				},
				{
					"id": "goal-2",
					"title": "Goal 2",
					"description": "Description 2",
					"start_date": "2025-10-02",
					"end_date": "2025-11-02",
					"kpi_name": "Kpi Name 2",
					"kpi_target": 100,
					"kpi_unit": "Kpi Unit 2",
					"status": "done",
					"created_at": "2025-10-01T00:00:00+09:00",
					"updated_at": "2025-10-02T00:00:00+09:00"
				}
			]
		}`, response)
	},
	)
}
