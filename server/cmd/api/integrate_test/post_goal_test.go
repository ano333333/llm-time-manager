package integratetest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	setuphandlers "github.com/ano333333/llm-time-manager/server/cmd/api/setup"
	"github.com/ano333333/llm-time-manager/server/internal/utils"
	"github.com/stretchr/testify/assert"
)

type responseInvalidParameterValidation struct {
	Message string `json:"message" validate:"required,eq=invalid parameter"`
	Target  string `json:"target" validate:"oneof=title description start_date end_date kpi_name kpi_target kpi_unit status"`
}

type responseGoalUnitWithKpi struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	KpiName     string  `json:"kpi_name"`
	KpiTarget   float64 `json:"kpi_target"`
	KpiUnit     string  `json:"kpi_unit"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type responseGoalUnitWithoutKpi struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type responseWithKpi struct {
	Goal responseGoalUnitWithKpi `json:"goal"`
}

type responseWithoutKpi struct {
	Goal responseGoalUnitWithoutKpi `json:"goal"`
}

type responseGetWithKpi struct {
	Goals []responseGoalUnitWithKpi `json:"goals"`
}

type responseGetWithoutKpi struct {
	Goals []responseGoalUnitWithoutKpi `json:"goals"`
}

func TestPostGoalIntegrate(t *testing.T) {
	t.Run("POST /goal はリクエストパラメータが不正な場合 400 Bad Request を返す", func(t *testing.T) {
		// Arrange
		// リクエストパターン:
		// - パラメータのいずれかが欠如
		// - パラメータのいずれかが不正な型
		// - title が空白文字のみで構成されている
		// - start_dateがYYYY-MM-DD形式でない
		// - end_dateがYYYY-MM-DD形式でない
		// - start_date が end_date より大きい
		// - kpi_name のみが非 null で他が null
		// - kpi_target のみが非 null で他が null
		// - kpi_unit のみが非 null で他が null
		// - kpi_name, kpi_target, kpi_unit がすべて非nullで、kpi_name が空白文字のみで構成されている
		// - kpi_name, kpi_target, kpi_unit がすべて非nullで、kpi_unit が空白文字のみで構成されている
		// - status が "active"|"paused"|"done" 以外の値
		db, err := BeforeEach()
		if err != nil {
			t.Fatalf("failed to set up test: %v", err)
		}
		defer AfterEach(db)
		mux := setuphandlers.SetupHandlers(db)
		validator := utils.GetValidator()

		for i, request := range badRequests {
			t.Logf("request %d: %v", i, request)
			// Act
			body, _ := json.Marshal(request)
			req := httptest.NewRequest(http.MethodPost, "/goal", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			// Assert
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
			response, err := GetResponseBodyJson(rec)
			t.Logf("response %d: %v", i, response)
			assert.NoError(t, err)
			typedResponse := responseInvalidParameterValidation{}
			if err := json.Unmarshal([]byte(response), &typedResponse); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}
			if err := validator.Struct(typedResponse); err != nil {
				t.Fatalf("failed to validate response: %v", err)
			}
			t.Logf("request %d finished", i)
		}
	})

	t.Run("POST /goal は200ステータスと新規作成したGoalオブジェクトを返す", func(t *testing.T) {
		// Arrange
		db, err := BeforeEach()
		if err != nil {
			t.Fatalf("failed to set up test: %v", err)
		}
		defer AfterEach(db)
		mux := setuphandlers.SetupHandlers(db)
		validator := utils.GetValidator()
		requestWithKpi := map[string]interface{}{
			"title":       "title1",
			"description": "description1",
			"start_date":  "2025-10-01",
			"end_date":    "2025-12-31",
			"kpi_name":    "kpi_name1",
			"kpi_target":  10.,
			"kpi_unit":    "kpi_unit1",
			"status":      "active",
		}
		requestWithoutKpi := map[string]interface{}{
			"title":       "title2",
			"description": "description2",
			"start_date":  "2026-10-01",
			"end_date":    "2026-12-31",
			"kpi_name":    nil,
			"kpi_target":  nil,
			"kpi_unit":    nil,
			"status":      "paused",
		}

		// Act
		body, _ := json.Marshal(requestWithKpi)
		req := httptest.NewRequest(http.MethodPost, "/goal", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err := GetResponseBodyJson(rec)
		assert.NoError(t, err)
		responseWithKpi := responseWithKpi{}
		if err := json.Unmarshal([]byte(response), &responseWithKpi); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if err := validator.Struct(responseWithKpi); err != nil {
			t.Fatalf("failed to validate response: %v", err)
		}
		assert.Equal(t, requestWithKpi["title"], responseWithKpi.Goal.Title)
		assert.Equal(t, requestWithKpi["description"], responseWithKpi.Goal.Description)
		assert.Equal(t, requestWithKpi["start_date"], responseWithKpi.Goal.StartDate)
		assert.Equal(t, requestWithKpi["end_date"], responseWithKpi.Goal.EndDate)
		assert.Equal(t, requestWithKpi["kpi_name"], responseWithKpi.Goal.KpiName)
		assert.Equal(t, requestWithKpi["kpi_target"], responseWithKpi.Goal.KpiTarget)
		assert.Equal(t, requestWithKpi["kpi_unit"], responseWithKpi.Goal.KpiUnit)
		assert.Equal(t, requestWithKpi["status"], responseWithKpi.Goal.Status)

		// Act
		body, _ = json.Marshal(requestWithoutKpi)
		req = httptest.NewRequest(http.MethodPost, "/goal", bytes.NewBuffer(body))
		rec = httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", strings.ToLower(rec.Header().Get("Content-Type")))
		response, err = GetResponseBodyJson(rec)
		assert.NoError(t, err)
		responseWithoutKpi := responseWithoutKpi{}
		if err := json.Unmarshal([]byte(response), &responseWithoutKpi); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if err := validator.Struct(responseWithoutKpi); err != nil {
			t.Fatalf("failed to validate response: %v", err)
		}
		assert.Equal(t, requestWithoutKpi["title"], responseWithoutKpi.Goal.Title)
		assert.Equal(t, requestWithoutKpi["description"], responseWithoutKpi.Goal.Description)
		assert.Equal(t, requestWithoutKpi["start_date"], responseWithoutKpi.Goal.StartDate)
		assert.Equal(t, requestWithoutKpi["end_date"], responseWithoutKpi.Goal.EndDate)
		assert.Equal(t, requestWithoutKpi["status"], responseWithoutKpi.Goal.Status)
	})

	t.Run("POST /goal で新規作成したGoalオブジェクトはGET /goal で取得できる", func(t *testing.T) {
		// Arrange
		db, err := BeforeEach()
		if err != nil {
			t.Fatalf("failed to set up test: %v", err)
		}
		defer AfterEach(db)
		mux := setuphandlers.SetupHandlers(db)
		requestWithKpi := map[string]interface{}{
			"title":       "title1",
			"description": "description1",
			"start_date":  "2025-10-01",
			"end_date":    "2025-12-31",
			"kpi_name":    "kpi_name1",
			"kpi_target":  10,
			"kpi_unit":    "kpi_unit1",
			"status":      "active",
		}
		requestWithoutKpi := map[string]interface{}{
			"title":       "title2",
			"description": "description2",
			"start_date":  "2026-10-01",
			"end_date":    "2026-12-31",
			"kpi_name":    nil,
			"kpi_target":  nil,
			"kpi_unit":    nil,
			"status":      "paused",
		}

		requestBodyWithKpi, _ := json.Marshal(requestWithKpi)
		reqWithKpi := httptest.NewRequest(http.MethodPost, "/goal", bytes.NewBuffer(requestBodyWithKpi))
		recWithKpi := httptest.NewRecorder()
		mux.ServeHTTP(recWithKpi, reqWithKpi)
		responseBodyWithKpi, err := GetResponseBodyJson(recWithKpi)
		if err != nil {
			t.Fatalf("failed to get response body: %v", err)
		}
		responseResultWithKpi := responseWithKpi{}
		if err := json.Unmarshal([]byte(responseBodyWithKpi), &responseResultWithKpi); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		bodyWithoutKpi, _ := json.Marshal(requestWithoutKpi)
		reqWithoutKpi := httptest.NewRequest(http.MethodPost, "/goal", bytes.NewBuffer(bodyWithoutKpi))
		recWithoutKpi := httptest.NewRecorder()
		mux.ServeHTTP(recWithoutKpi, reqWithoutKpi)
		responseBodyWithoutKpi, err := GetResponseBodyJson(recWithoutKpi)
		if err != nil {
			t.Fatalf("failed to get response body: %v", err)
		}
		responseResultWithoutKpi := responseWithoutKpi{}
		if err := json.Unmarshal([]byte(responseBodyWithoutKpi), &responseResultWithoutKpi); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		// Act
		reqGetWithKpi := httptest.NewRequest(http.MethodGet, "/goal?status=active", nil)
		recGetWithKpi := httptest.NewRecorder()
		mux.ServeHTTP(recGetWithKpi, reqGetWithKpi)
		responseBodyGetWithKpi, err := GetResponseBodyJson(recGetWithKpi)
		if err != nil {
			t.Fatalf("failed to get response body: %v", err)
		}
		t.Logf("responseBodyGetWithKpi: %v", responseBodyGetWithKpi)
		responseResultGetWithKpi := responseGetWithKpi{}
		if err := json.Unmarshal([]byte(responseBodyGetWithKpi), &responseResultGetWithKpi); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		t.Logf("responseResultGetWithKpi: %v", responseResultGetWithKpi)

		// Assert
		assert.Equal(t, http.StatusOK, recGetWithKpi.Code)
		assert.Equal(t, "application/json", strings.ToLower(recGetWithKpi.Header().Get("Content-Type")))
		assert.Equal(t, responseResultWithKpi.Goal.ID, responseResultGetWithKpi.Goals[0].ID)
		assert.Equal(t, responseResultWithKpi.Goal.Title, responseResultGetWithKpi.Goals[0].Title)
		assert.Equal(t, responseResultWithKpi.Goal.Description, responseResultGetWithKpi.Goals[0].Description)
		assert.Equal(t, responseResultWithKpi.Goal.StartDate, responseResultGetWithKpi.Goals[0].StartDate)
		assert.Equal(t, responseResultWithKpi.Goal.EndDate, responseResultGetWithKpi.Goals[0].EndDate)
		assert.Equal(t, responseResultWithKpi.Goal.KpiName, responseResultGetWithKpi.Goals[0].KpiName)
		assert.Equal(t, responseResultWithKpi.Goal.KpiTarget, responseResultGetWithKpi.Goals[0].KpiTarget)
		assert.Equal(t, responseResultWithKpi.Goal.KpiUnit, responseResultGetWithKpi.Goals[0].KpiUnit)
		assert.Equal(t, responseResultWithKpi.Goal.Status, responseResultGetWithKpi.Goals[0].Status)

		// Act
		reqGetWithoutKpi := httptest.NewRequest(http.MethodGet, "/goal?status=paused", nil)
		recGetWithoutKpi := httptest.NewRecorder()
		mux.ServeHTTP(recGetWithoutKpi, reqGetWithoutKpi)
		responseBodyGetWithoutKpi, err := GetResponseBodyJson(recGetWithoutKpi)
		if err != nil {
			t.Fatalf("failed to get response body: %v", err)
		}
		responseResultGetWithoutKpi := responseGetWithoutKpi{}
		if err := json.Unmarshal([]byte(responseBodyGetWithoutKpi), &responseResultGetWithoutKpi); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		// Assert
		assert.Equal(t, http.StatusOK, recGetWithoutKpi.Code)
		assert.Equal(t, "application/json", strings.ToLower(recGetWithoutKpi.Header().Get("Content-Type")))
		assert.Equal(t, responseResultWithoutKpi.Goal.ID, responseResultGetWithoutKpi.Goals[0].ID)
		assert.Equal(t, responseResultWithoutKpi.Goal.Title, responseResultGetWithoutKpi.Goals[0].Title)
		assert.Equal(t, responseResultWithoutKpi.Goal.Description, responseResultGetWithoutKpi.Goals[0].Description)
		assert.Equal(t, responseResultWithoutKpi.Goal.StartDate, responseResultGetWithoutKpi.Goals[0].StartDate)
		assert.Equal(t, responseResultWithoutKpi.Goal.EndDate, responseResultGetWithoutKpi.Goals[0].EndDate)
		assert.Equal(t, responseResultWithoutKpi.Goal.Status, responseResultGetWithoutKpi.Goals[0].Status)
	})
}

var badRequests = []map[string]interface{}{
	{
		"description": "...",
		"start_date":  "2025-10-01",
		"end_date":    "2025-12-31",
		"kpi_name":    "集中作業時間",
		"kpi_target":  10,
		"kpi_unit":    "時間",
		"status":      "active",
	},
	{
		"title":      "週10時間の集中作業",
		"start_date": "2025-10-01",
		"end_date":   "2025-12-31",
		"kpi_name":   "集中作業時間",
		"kpi_target": 10,
		"kpi_unit":   "時間",
		"status":     "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"end_date":    "2025-12-31",
		"kpi_name":    "集中作業時間",
		"kpi_target":  10,
		"kpi_unit":    "時間",
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"start_date":  "2025-10-01",
		"kpi_name":    "集中作業時間",
		"kpi_target":  10,
		"kpi_unit":    "時間",
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"start_date":  "2025-10-01",
		"end_date":    "2025-12-31",
		"kpi_target":  10,
		"kpi_unit":    "時間",
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"start_date":  "2025-10-01",
		"end_date":    "2025-12-31",
		"kpi_name":    "集中作業時間",
		"kpi_unit":    "時間",
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"start_date":  "2025-10-01",
		"end_date":    "2025-12-31",
		"kpi_name":    "集中作業時間",
		"kpi_target":  10,
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"start_date":  "2025-10-01",
		"end_date":    "2025-12-31",
		"kpi_name":    "集中作業時間",
		"kpi_target":  10,
		"kpi_unit":    "時間",
	},
	{
		"title":       2,
		"description": "...",
		"start_date":  "2025-10-01",
		"end_date":    "2025-12-31",
		"kpi_name":    "集中作業時間",
		"kpi_target":  10,
		"kpi_unit":    "時間",
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": 5,
		"start_date":  "2025-10-01",
		"end_date":    "2025-12-31",
		"kpi_name":    "集中作業時間",
		"kpi_target":  10,
		"kpi_unit":    "時間",
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"start_date":  0,
		"end_date":    "2025-12-31",
		"kpi_name":    "集中作業時間",
		"kpi_target":  10,
		"kpi_unit":    "時間",
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"start_date":  "2025-10-01",
		"end_date":    1000,
		"kpi_name":    "集中作業時間",
		"kpi_target":  10,
		"kpi_unit":    "時間",
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"start_date":  "2025-10-01",
		"end_date":    "2025-12-31",
		"kpi_name":    5,
		"kpi_target":  10,
		"kpi_unit":    "時間",
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"start_date":  "2025-10-01",
		"end_date":    "2025-12-31",
		"kpi_name":    "集中作業時間",
		"kpi_target":  "10",
		"kpi_unit":    "時間",
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"start_date":  "2025-10-01",
		"end_date":    "2025-12-31",
		"kpi_name":    "集中作業時間",
		"kpi_target":  10,
		"kpi_unit":    100,
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"start_date":  "2025-10-01",
		"end_date":    "2025-12-31",
		"kpi_name":    "集中作業時間",
		"kpi_target":  10,
		"kpi_unit":    "時間",
		"status":      3,
	},
	{
		"title":       "    ",
		"description": "...",
		"start_date":  "2025-10-01",
		"end_date":    "2025-12-31",
		"kpi_name":    "集中作業時間",
		"kpi_target":  10,
		"kpi_unit":    "時間",
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"start_date":  "2025/10/01",
		"end_date":    "2025-12-31",
		"kpi_name":    "集中作業時間",
		"kpi_target":  10,
		"kpi_unit":    "時間",
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"start_date":  "2025-10-01",
		"end_date":    "2025-12-31T00:00:00+09:00",
		"kpi_name":    "集中作業時間",
		"kpi_target":  10,
		"kpi_unit":    "時間",
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"start_date":  "2025-12-31",
		"end_date":    "2025-10-01",
		"kpi_name":    "集中作業時間",
		"kpi_target":  10,
		"kpi_unit":    "時間",
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"start_date":  "2025-10-01",
		"end_date":    "2025-12-31",
		"kpi_name":    nil,
		"kpi_target":  10,
		"kpi_unit":    "時間",
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"start_date":  "2025-10-01",
		"end_date":    "2025-12-31",
		"kpi_name":    "集中作業時間",
		"kpi_target":  nil,
		"kpi_unit":    "時間",
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"start_date":  "2025-10-01",
		"end_date":    "2025-12-31",
		"kpi_name":    "集中作業時間",
		"kpi_unit":    nil,
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"start_date":  "2025-10-01",
		"end_date":    "2025-12-31",
		"kpi_name":    "集中作業時間",
		"kpi_target":  10,
		"kpi_unit":    nil,
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"start_date":  "2025-10-01",
		"end_date":    "2025-12-31",
		"kpi_name":    "     ",
		"kpi_target":  10,
		"kpi_unit":    "時間",
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"start_date":  "2025-10-01",
		"end_date":    "2025-12-31",
		"kpi_name":    "集中作業時間",
		"kpi_target":  10,
		"kpi_unit":    "     ",
		"status":      "active",
	},
	{
		"title":       "週10時間の集中作業",
		"description": "...",
		"start_date":  "2025-10-01",
		"end_date":    "2025-12-31",
		"kpi_name":    "集中作業時間",
		"kpi_target":  10,
		"kpi_unit":    "時間",
		"status":      "invalid",
	},
}
