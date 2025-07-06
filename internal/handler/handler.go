package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/otakakot/sample-go-api/pkg/api"
)

func Health(
	w http.ResponseWriter,
	r *http.Request,
) {
	req := api.HealthRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Fprintf(w, "OK")

		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(api.HealthResponse{
		Message: req.Message,
	}); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)

		return
	}
}
