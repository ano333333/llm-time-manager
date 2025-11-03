package integratetest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	setuphandlers "github.com/ano333333/llm-time-manager/server/cmd/api/setup"
	datamodel "github.com/ano333333/llm-time-manager/server/internal/data-model"
	"github.com/stretchr/testify/assert"
)

func TestGetGoalsIntegrate(t *testing.T) {
	t.Run("GET /goal はquery parameterがない場合空配列を返す", func(t *testing.T) {
		// Arrange
		db, err := BeforeEach()
		if err != nil {
			t.Fatalf("failed to set up test: %v", err)
		}
		defer AfterEach(db)
		mux := setuphandlers.SetupHandlers(db)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/goal", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err := GetResponseBodyJson(rec)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"goals": []
		}`, response)
	})

	t.Run("GET /goal は不正なquery parameterで400を返す", func(t *testing.T) {
		// Arrange
		db, err := BeforeEach()
		if err != nil {
			t.Fatalf("failed to set up test: %v", err)
		}
		defer AfterEach(db)
		mux := setuphandlers.SetupHandlers(db)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/goal?status=invalid", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("GET /goal はquery parameterにstatusがマッチするモデルを返す", func(t *testing.T) {
		// Arrange
		db, err := BeforeEach()
		if err != nil {
			t.Fatalf("failed to set up test: %v", err)
		}
		defer AfterEach(db)
		mux := setuphandlers.SetupHandlers(db)
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
		if err := InsertGoals(db, goals); err != nil {
			t.Fatalf("failed to insert goals: %v", err)
		}

		// Act(none)
		req := httptest.NewRequest(http.MethodGet, "/goal?status=", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert(none)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err := GetResponseBodyJson(rec)
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
		response, err = GetResponseBodyJson(rec)
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
		response, err = GetResponseBodyJson(rec)
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
		response, err = GetResponseBodyJson(rec)
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
		response, err = GetResponseBodyJson(rec)
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
		response, err = GetResponseBodyJson(rec)
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
		response, err = GetResponseBodyJson(rec)
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
		response, err = GetResponseBodyJson(rec)
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
		response, err = GetResponseBodyJson(rec)
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
