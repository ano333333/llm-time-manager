package integratetest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	setuphandlers "github.com/ano333333/llm-time-manager/server/cmd/api/setup"
	datamodel "github.com/ano333333/llm-time-manager/server/internal/data-model"
	"github.com/stretchr/testify/assert"
)

func TestGetCaptureScheduleIntegrate(t *testing.T) {
	t.Run("GET /capture/schedule はアクティブなスケジュールがなければnullを返す", func(t *testing.T) {
		// Arrange
		db, err := BeforeEach()
		if err != nil {
			t.Fatalf("failed to set up test: %v", err)
		}
		defer AfterEach(db)
		mux := setuphandlers.SetupHandlers(db)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/capture/schedule", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err := GetResponseBodyJson(rec)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"schedule": null}`, response)
	})

	t.Run("GET /capture/schedule はアクティブなスケジュールが 1 つだけあればそれを返す", func(t *testing.T) {
		// Arrange
		db, err := BeforeEach()
		if err != nil {
			t.Fatalf("failed to set up test: %v", err)
		}
		defer AfterEach(db)
		mux := setuphandlers.SetupHandlers(db)
		now := time.Now()
		schedules := []datamodel.CaptureSchedule{
			{
				ID:          "schedule-0",
				Active:      false,
				IntervalMin: 10,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          "schedule-1",
				Active:      true,
				IntervalMin: 5,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          "schedule-2",
				Active:      false,
				IntervalMin: 15,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}
		if err := InsertCaptureSchedules(db, schedules); err != nil {
			t.Fatalf("failed to insert capture schedules: %v", err)
		}

		// Act
		req := httptest.NewRequest(http.MethodGet, "/capture/schedule", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err := GetResponseBodyJson(rec)
		assert.NoError(t, err)
		expected, err := json.Marshal(map[string]interface{}{
			"schedule": map[string]interface{}{
				"id":           schedules[1].ID,
				"active":       schedules[1].Active,
				"interval_min": schedules[1].IntervalMin,
			},
		})
		if err != nil {
			t.Fatalf("failed to marshal expected: %v", err)
		}
		assert.JSONEq(t, string(expected), response)
	})

	t.Run("GET /capture/schedule は複数のアクティブなスケジュールがあれば 500 Internal Server Error を返す", func(t *testing.T) {
		// Arrange
		db, err := BeforeEach()
		if err != nil {
			t.Fatalf("failed to set up test: %v", err)
		}
		defer AfterEach(db)
		mux := setuphandlers.SetupHandlers(db)
		now := time.Now()
		schedules := []datamodel.CaptureSchedule{
			{
				ID:          "schedule-0",
				Active:      true,
				IntervalMin: 10,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          "schedule-1",
				Active:      true,
				IntervalMin: 5,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}
		if err := InsertCaptureSchedules(db, schedules); err != nil {
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
