package version

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"testing"
)

func TestGetVersion(t *testing.T) {
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
				r:              httptest.NewRequest(http.MethodGet, "/version", nil),
				expectedStatus: 200,
			},
		},
		{
			name: "POST method",
			args: args{
				w:              httptest.NewRecorder(),
				r:              httptest.NewRequest(http.MethodPost, "/version", nil),
				expectedStatus: 405,
			},
		},
		{
			name: "PUT method",
			args: args{
				w:              httptest.NewRecorder(),
				r:              httptest.NewRequest(http.MethodPut, "/version", nil),
				expectedStatus: 405,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			GetVersion(w, tt.args.r)
			res := w.Result()

			if res.StatusCode != tt.args.expectedStatus {
				t.Errorf("%s: Unexpected Status Code %v", tt.name, res.StatusCode)
			}
		})
	}
}

func Test_initBuild(t *testing.T) {
	goversion := runtime.Version()
	type compileVars struct {
		version string
		sha     string
	}
	tests := []struct {
		name  string
		want  build
		input compileVars
	}{
		{
			name: "Nothing specified",
			want: build{
				Version:   "development",
				GitSHA:    "",
				GoVersion: goversion,
			},
			input: compileVars{
				version: "",
				sha:     "",
			},
		},
		{
			name: "Long SHA",
			want: build{
				Version:   "1.2.3",
				GitSHA:    "1234567",
				GoVersion: goversion,
			},
			input: compileVars{
				version: "1.2.3",
				sha:     "123456789",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version = tt.input.version
			gitSHA = tt.input.sha
			if got := initBuild(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s: initBuild() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
