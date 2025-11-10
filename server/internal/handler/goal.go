package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	datamodel "github.com/ano333333/llm-time-manager/server/internal/data-model"
	"github.com/ano333333/llm-time-manager/server/internal/store"
)

type GoalHandler struct {
	GoalStore        store.GoalStore
	TransactionStore store.TransactionStore
}

func (h *GoalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
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
