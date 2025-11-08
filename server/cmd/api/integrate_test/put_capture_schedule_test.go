package integratetest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	setuphandlers "github.com/ano333333/llm-time-manager/server/cmd/api/setup"
	datamodel "github.com/ano333333/llm-time-manager/server/internal/data-model"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type ResponseInvalidParameterValidation struct {
	Message string `json:"message" validate:"required,eq=invalid parameter"`
	Target  string `json:"target" validate:"required,oneof=active interval_min"`
}

type ResponseBadRequestValidation struct {
	Message string `json:"message" validate:"required,eq=no active capture schedule found"`
}

func TestPutCaptureScheduleIntegrate(t *testing.T) {
	t.Run("PUT /capture/schedule はリクエストパラメータが不正な場合 400 Bad Request を返す", func(t *testing.T) {
		// Arrange
		// リクエストパターン:
		// - active が欠如
		// - intervalMin が欠如
		// - active が非 boolean
		// - intervalMin が非整数
		// - intervalMin が負数
		// - intervalMin が 0
		requests := []map[string]interface{}{
			{
				"interval_min": 5,
			},
			{
				"active": true,
			},
			{
				"active":       "true",
				"interval_min": 5,
			},
			{
				"active":       true,
				"interval_min": "5",
			},
			{
				"active":       true,
				"interval_min": -1,
			},
			{
				"active":       true,
				"interval_min": 0,
			},
		}
		db, err := BeforeEach()
		if err != nil {
			t.Fatalf("failed to set up test: %v", err)
		}
		defer AfterEach(db)
		mux := setuphandlers.SetupHandlers(db)
		validator := validator.New()

		for _, request := range requests {
			// Act
			body, _ := json.Marshal(request)
			req := httptest.NewRequest(http.MethodPut, "/capture/schedule", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			// Assert
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
			response, err := GetResponseBodyJson(rec)
			assert.NoError(t, err)
			typedResponse := ResponseInvalidParameterValidation{}
			if err := json.Unmarshal([]byte(response), &typedResponse); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}
			if err := validator.Struct(typedResponse); err != nil {
				t.Fatalf("failed to validate response: %v", err)
			}
		}
	})

	t.Run("PUT /capture/schedule はスケジュールがない場合 400 Bad Request を返す", func(t *testing.T) {
		// Arrange
		db, err := BeforeEach()
		if err != nil {
			t.Fatalf("failed to set up test: %v", err)
		}
		defer AfterEach(db)
		mux := setuphandlers.SetupHandlers(db)
		validator := validator.New()

		// Act
		req := httptest.NewRequest(http.MethodPut, "/capture/schedule", bytes.NewBuffer([]byte(`{"active": true, "interval_min": 5}`)))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err := GetResponseBodyJson(rec)
		assert.NoError(t, err)
		typedResponse := ResponseBadRequestValidation{}
		if err := json.Unmarshal([]byte(response), &typedResponse); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if err := validator.Struct(typedResponse); err != nil {
			t.Fatalf("failed to validate response: %v", err)
		}
	})

	t.Run("PUT /capture/schedule はアクティブなスケジュールが無い場合 400 Bad Request を返す", func(t *testing.T) {
		// Arrange
		db, err := BeforeEach()
		if err != nil {
			t.Fatalf("failed to set up test: %v", err)
		}
		defer AfterEach(db)
		mux := setuphandlers.SetupHandlers(db)
		validator := validator.New()
		timezone := GetJSTTimezone()
		now := time.Now().In(timezone)
		schedules := []datamodel.CaptureSchedule{
			{
				ID:          "schedule-0",
				Active:      false,
				IntervalMin: 10,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}
		if err := InsertCaptureSchedules(db, schedules); err != nil {
			t.Fatalf("failed to insert capture schedules: %v", err)
		}

		// Act
		req := httptest.NewRequest(http.MethodPut, "/capture/schedule", bytes.NewBuffer([]byte(`{"active": true, "interval_min": 5}`)))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err := GetResponseBodyJson(rec)
		assert.NoError(t, err)
		typedResponse := ResponseBadRequestValidation{}
		if err := json.Unmarshal([]byte(response), &typedResponse); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if err := validator.Struct(typedResponse); err != nil {
			t.Fatalf("failed to validate response: %v", err)
		}
	})

	t.Run("PUT /capture/schedule はアクティブなスケジュールがある場合更新し更新後のスケジュールを返す", func(t *testing.T) {
		// Arrange
		db, err := BeforeEach()
		if err != nil {
			t.Fatalf("failed to set up test: %v", err)
		}
		defer AfterEach(db)
		mux := setuphandlers.SetupHandlers(db)
		timezone := GetJSTTimezone()
		now := time.Now().In(timezone)
		schedules := []datamodel.CaptureSchedule{
			{
				ID:          "schedule-0",
				Active:      true,
				IntervalMin: 10,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}
		if err := InsertCaptureSchedules(db, schedules); err != nil {
			t.Fatalf("failed to insert capture schedules: %v", err)
		}

		// Act
		req := httptest.NewRequest(http.MethodPut, "/capture/schedule", bytes.NewBuffer([]byte(`{"active": true, "interval_min": 5}`)))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err := GetResponseBodyJson(rec)
		assert.NoError(t, err)
		expected, _ := json.Marshal(map[string]interface{}{
			"schedule": map[string]interface{}{
				"id":           "schedule-0",
				"active":       true,
				"interval_min": 5,
			},
		})
		assert.JSONEq(t, string(expected), response)
		reqGet := httptest.NewRequest(http.MethodGet, "/capture/schedule", nil)
		recGet := httptest.NewRecorder()
		mux.ServeHTTP(recGet, reqGet)
		responseGet, _ := GetResponseBodyJson(recGet)
		assert.JSONEq(t, string(expected), responseGet)

		// Act
		req = httptest.NewRequest(http.MethodPut, "/capture/schedule", bytes.NewBuffer([]byte(`{"active": false, "interval_min": 5}`)))
		rec = httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err = GetResponseBodyJson(rec)
		assert.NoError(t, err)
		expected, _ = json.Marshal(map[string]interface{}{
			"schedule": nil,
		})
		assert.JSONEq(t, string(expected), response)
	})
}
