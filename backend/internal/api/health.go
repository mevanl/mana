package api

import (
	"net/http"
)

func Health_Handler(write http.ResponseWriter, read *http.Request) {
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusOK)
	write.Write([]byte(`{"status":"ok"}`))
}
