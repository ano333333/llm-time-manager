package handler

import (
	"encoding/json"
	"log"
	"net/http"

	repositories "github.com/ano333333/llm-time-manager/server/internal/store"
	"github.com/ano333333/llm-time-manager/server/internal/utils"
)

type CaptureScheduleHandler struct {
	CaptureScheduleStore repositories.CaptureScheduleStore
	TransactionStore     repositories.TransactionStore
}

func (h *CaptureScheduleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w)
	case "PUT":
		h.put(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *CaptureScheduleHandler) get(w http.ResponseWriter) {
	tx, err := h.TransactionStore.Begin()
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	captureSchedule, err := h.CaptureScheduleStore.GetActiveCaptureSchedule(tx)
	if err != nil {
		log.Printf("failed to get active capture schedule: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if err := tx.Commit(); err != nil {
		log.Printf("failed to commit transaction: %v", err)
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

func (h *CaptureScheduleHandler) put(w http.ResponseWriter, r *http.Request) {
	validator := utils.GetValidator()

	type RequestBodyValidation struct {
		Active      any `json:"active" validate:"required,is_boolean"`
		IntervalMin any `json:"interval_min" validate:"required,is_integer,min=1,max=1440"`
	}
	type RequestBody struct {
		Active      bool `json:"active"`
		IntervalMin int  `json:"interval_min"`
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
	requestBody := RequestBody{
		Active:      requestBodyValidation.Active.(bool),
		IntervalMin: (int)(requestBodyValidation.IntervalMin.(float64)),
	}

	tx, err := h.TransactionStore.Begin()
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	affectedRows, err := h.CaptureScheduleStore.UpdateActiveCaptureSchedule(tx, requestBody.Active, requestBody.IntervalMin)
	if err != nil {
		log.Printf("failed to update active capture schedule: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if affectedRows == 0 {
		log.Printf("no active capture schedule found")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "no active capture schedule found",
		})
		return
	}
	captureSchedule, err := h.CaptureScheduleStore.GetActiveCaptureSchedule(tx)
	if err != nil {
		log.Printf("failed to get active capture schedule: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("failed to commit transaction: %v", err)
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
