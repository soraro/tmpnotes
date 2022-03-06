package health

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	type args struct {
		w              http.ResponseWriter
		r              *http.Request
		expectedStatus int
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "GET method",
			args: args{
				w:              httptest.NewRecorder(),
				r:              httptest.NewRequest(http.MethodGet, "/healthz", nil),
				expectedStatus: 200,
			},
		},
		{
			name: "POST method",
			args: args{
				w:              httptest.NewRecorder(),
				r:              httptest.NewRequest(http.MethodPost, "/healthz", nil),
				expectedStatus: 405,
			},
		},
		{
			name: "PUT method",
			args: args{
				w:              httptest.NewRecorder(),
				r:              httptest.NewRequest(http.MethodPut, "/healthz", nil),
				expectedStatus: 405,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			HealthCheck(w, tt.args.r)
			res := w.Result()

			if res.StatusCode != tt.args.expectedStatus {
				t.Errorf("Unexpected Status Code")
			}
		})
	}
}
