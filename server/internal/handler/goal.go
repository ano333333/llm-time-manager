package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	datamodel "github.com/ano333333/llm-time-manager/server/internal/data-model"
	"github.com/ano333333/llm-time-manager/server/internal/store"
	"github.com/ano333333/llm-time-manager/server/internal/utils"
)

type GoalHandler struct {
	GoalStore        store.GoalStore
	TransactionStore store.TransactionStore
}

func (h *GoalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
	case "POST":
		h.post(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *GoalHandler) get(w http.ResponseWriter, r *http.Request) {
	statusRaw := r.URL.Query().Get("status")
	if statusRaw == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"goals": make([]datamodel.Goal, 0),
		})
		return
	}

	status := strings.Split(statusRaw, ",")
	for i, s := range status {
		status[i] = strings.TrimSpace(s)
		s = status[i]
		if s != "active" && s != "paused" && s != "done" {
			http.Error(w, fmt.Sprintf("invalid status: %s", s), http.StatusBadRequest)
			return
		}
	}

	tx, err := h.TransactionStore.Begin()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	goals, err := h.GoalStore.GetGoal(tx, status)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if err := tx.Commit(); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	results := make([](map[string]interface{}), 0)
	for _, goal := range goals {
		results = append(results, map[string]interface{}{
			"id":          goal.ID,
			"title":       goal.Title,
			"description": goal.Description,
			"start_date":  goal.StartDate.Format("2006-01-02"),
			"end_date":    goal.EndDate.Format("2006-01-02"),
			"kpi_name":    goal.KpiName,
			"kpi_target":  goal.KpiTarget,
			"kpi_unit":    goal.KpiUnit,
			"status":      goal.Status,
			"created_at":  goal.CreatedAt,
			"updated_at":  goal.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"goals": results,
	})
}

func (h *GoalHandler) post(w http.ResponseWriter, r *http.Request) {
	validator := utils.GetValidator()

	type RequestBodyValidation struct {
		Title       any `json:"title" validate:"required,min=1,max=255,is_string,not_consists_of_whitespaces"`
		Description any `json:"description" validate:"required,is_string"`
		StartDate   any `json:"start_date" validate:"required,is_string,datetime=2006-01-02"`
		EndDate     any `json:"end_date" validate:"required,is_string,datetime=2006-01-02"`
		KpiName     any `json:"kpi_name" validate:"required_with_all=KpiTarget KpiUnit,is_nullable_string"`
		KpiTarget   any `json:"kpi_target" validate:"required_with_all=KpiName KpiUnit,is_nullable_float64"`
		KpiUnit     any `json:"kpi_unit" validate:"required_with_all=KpiName KpiTarget,is_nullable_string"`
		Status      any `json:"status" validate:"required,is_string,oneof=active paused done"`
	}
	type RequestBody struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		StartDate   string   `json:"start_date"`
		EndDate     string   `json:"end_date"`
		KpiName     *string  `json:"kpi_name"`
		KpiTarget   *float64 `json:"kpi_target"`
		KpiUnit     *string  `json:"kpi_unit"`
		Status      string   `json:"status"`
	}
	var requestBodyValidation RequestBodyValidation
	if err := json.NewDecoder(r.Body).Decode(&requestBodyValidation); err != nil {
		log.Printf("failed to decode request body: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "invalid JSON format",
		})
		return
	}
	if err := validator.Struct(requestBodyValidation); err != nil {
		log.Printf("failed to validate request body: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "invalid parameter",
			"target":  utils.GetFirstValidationErrorTarget(err),
		})
		return
	}
	var requestBody RequestBody
	requestBody.Title = requestBodyValidation.Title.(string)
	requestBody.Description = requestBodyValidation.Description.(string)
	requestBody.StartDate = requestBodyValidation.StartDate.(string)
	requestBody.EndDate = requestBodyValidation.EndDate.(string)
	if requestBodyValidation.KpiName != nil {
		requestBody.KpiName = new(string)
		*requestBody.KpiName = requestBodyValidation.KpiName.(string)
	}
	if requestBodyValidation.KpiTarget != nil {
		requestBody.KpiTarget = new(float64)
		*requestBody.KpiTarget = requestBodyValidation.KpiTarget.(float64)
	}
	if requestBodyValidation.KpiUnit != nil {
		requestBody.KpiUnit = new(string)
		*requestBody.KpiUnit = requestBodyValidation.KpiUnit.(string)
	}
	requestBody.Status = requestBodyValidation.Status.(string)
	// validationでのタグチェックが面倒なものは自前チェック
	// validationのgtefieldはNumbersまたはtime.~にしか効かない
	if requestBody.StartDate > requestBody.EndDate {
		log.Printf("start date is after end date")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "invalid parameter",
			"target":  "end_date",
		})
		return
	}
	// kpi_nameとkpi_unitの非空白文字チェックもvalidatorだとnullableの処理が面倒
	if requestBody.KpiName != nil && strings.TrimSpace(*requestBody.KpiName) == "" {
		log.Printf("kpi name is empty")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "invalid parameter",
			"target":  "kpi_name",
		})
		return
	}
	if requestBody.KpiUnit != nil && strings.TrimSpace(*requestBody.KpiUnit) == "" {
		log.Printf("kpi unit is empty")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "invalid parameter",
			"target":  "kpi_unit",
		})
		return
	}
	// kpi_*のnull/非nullが揃っているか
	target := ""
	if requestBody.KpiName != nil {
		if requestBody.KpiTarget == nil {
			target = "kpi_target"
		}
		if requestBody.KpiUnit == nil {
			target = "kpi_unit"
		}
	} else {
		if requestBody.KpiTarget != nil {
			target = "kpi_target"
		}
		if requestBody.KpiUnit != nil {
			target = "kpi_unit"
		}
	}
	if target != "" {
		log.Printf("%s is not consistent with kpi_name", target)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "invalid parameter",
			"target":  target,
		})
		return
	}

	startDate, _ := time.Parse("2006-01-02", requestBody.StartDate)
	endDate, _ := time.Parse("2006-01-02", requestBody.EndDate)
	tx, err := h.TransactionStore.Begin()
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	goal, err := h.GoalStore.CreateGoal(tx, nil, requestBody.Title, requestBody.Description, startDate, endDate, requestBody.KpiName, requestBody.KpiTarget, requestBody.KpiUnit, requestBody.Status)
	if err != nil {
		log.Printf("failed to create goal: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if err := tx.Commit(); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	timezone := utils.GetJSTTimezone()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"goal": map[string]interface{}{
			"id":          goal.ID,
			"title":       goal.Title,
			"description": goal.Description,
			"start_date":  goal.StartDate.Format("2006-01-02"),
			"end_date":    goal.EndDate.Format("2006-01-02"),
			"kpi_name":    goal.KpiName,
			"kpi_target":  goal.KpiTarget,
			"kpi_unit":    goal.KpiUnit,
			"status":      goal.Status,
			"created_at":  goal.CreatedAt.In(timezone).Format(time.RFC3339),
			"updated_at":  goal.UpdatedAt.In(timezone).Format(time.RFC3339),
		},
	})
}
