package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ano333333/llm-time-manager/server/cmd/api/repositories"
)

type CaptureScheduleHandler struct {
	CaptureScheduleRepository repositories.CaptureScheduleRepository
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
	captureSchedule, err := h.CaptureScheduleRepository.GetActiveCaptureSchedule()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"schedule": captureSchedule,
	})
}
