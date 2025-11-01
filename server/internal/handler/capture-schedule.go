package handler

import (
	"encoding/json"
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"schedule": captureSchedule,
	})
}
