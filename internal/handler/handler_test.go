package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/otakakot/sample-go-api/internal/handler"
	"github.com/otakakot/sample-go-api/pkg/api"
)

type cases[Req, Res any] struct {
	name   string
	method string
	req    Req
	status int
	res    Res
}

func TestHealth(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.Health)

	tests := []cases[api.HealthRequest, api.HealthResponse]{
		{
			name:   "GET request returns OK",
			method: http.MethodGet,
			status: http.StatusOK,
			req: api.HealthRequest{
				Message: "test",
			},
			res: api.HealthResponse{
				Message: "test",
			},
		},
		{
			name:   "POST request returns Method Not Allowed",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
		},
		{
			name:   "PUT request returns Method Not Allowed",
			method: http.MethodPut,
			status: http.StatusMethodNotAllowed,
		},
		{
			name:   "DELETE request returns Method Not Allowed",
			method: http.MethodDelete,
			status: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			buf := bytes.NewBuffer(nil)

			if err := json.NewEncoder(buf).Encode(tt.req); err != nil {
				t.Fatalf("failed to encode request: %v", err)
			}

			req := httptest.NewRequestWithContext(t.Context(), tt.method, "/health", buf)

			res := httptest.NewRecorder()

			mux.ServeHTTP(res, req)

			if status := res.Code; status != tt.status {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.status)
			}

			if tt.res.Message == "" {
				return
			}

			got := &api.HealthResponse{}

			if err := json.NewDecoder(res.Body).Decode(got); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			gotB, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("failed to marshal response: %v", err)
			}

			wantB, err := json.Marshal(tt.res)
			if err != nil {
				t.Fatalf("failed to marshal expected response: %v", err)
			}

			if !bytes.Equal(gotB, wantB) {
				t.Errorf("handler returned unexpected body: got %s want %s", gotB, wantB)
			}
		})
	}
}
