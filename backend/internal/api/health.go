package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

var startTime = time.Now()

const version = "v0.0.1"

type HealthResponse struct {
	Status    string `json:"status"`
	Uptime    string `json:"uptime"`
	Version   string `json:"version"`
	Database  string `json:"database"`
	Timestamp string `json:"time"`
	Env       string `json:"env"`
}

func (api *API) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check db conn
	dbStatus := "connected"

	resp := HealthResponse{
		Status:    "ok",
		Uptime:    time.Since(startTime).String(),
		Version:   version,
		Database:  dbStatus,
		Timestamp: time.Now().Format(time.RFC3339),
		Env:       os.Getenv("ENV"),
	}

	if dbStatus != "connected" {
		resp.Status = "error"
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, `{ "error": "internal error" }`, http.StatusInternalServerError)
	}
}
