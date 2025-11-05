package handler

import (
	"encoding/json"
	"log"
	"net/http"

	repositories "github.com/ano333333/llm-time-manager/server/internal/store"
)

type CaptureScheduleHandler struct {
	CaptureScheduleStore repositories.CaptureScheduleStore
}

func (h *CaptureScheduleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w)
	default:
		http.NotFound(w, r)
	}
}

func (h *CaptureScheduleHandler) get(w http.ResponseWriter) {
	captureSchedule, err := h.CaptureScheduleStore.GetActiveCaptureSchedule()
	if err != nil {
		log.Printf("failed to get active capture schedule: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if captureSchedule != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"schedule": map[string]interface{}{
				"id":           captureSchedule.ID,
				"active":       captureSchedule.Active,
				"interval_min": captureSchedule.IntervalMin,
			},
		})
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"schedule": nil,
		})
	}
}
