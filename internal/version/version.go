package version

import (
	"encoding/json"
	"net/http"
	"runtime"

	log "github.com/sirupsen/logrus"
)

// NOTE: these variables are injected at build time
var (
	version string = "development"
	gitSHA  string
)

type build struct {
	Version   string `json:"version,omitempty"`
	GitSHA    string `json:"git,omitempty"`
	GoVersion string `json:"goversion,omitempty"`
}

func initBuild() build {
	var b build
	b.Version = version
	if b.Version == "" {
		b.Version = "development"
	}
	if len(gitSHA) >= 7 {
		b.GitSHA = gitSHA[:7]
	} else {
		b.GitSHA = gitSHA
	}
	b.GoVersion = runtime.Version()
	return b
}

func GetVersion(w http.ResponseWriter, r *http.Request) {
	log.Info(r.RequestURI)

	if r.Method != "GET" {
		log.Errorf("%s Invalid request method: %s", r.RequestURI, r.Method)
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	b := initBuild()

	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(b)
	if err != nil {
		log.Errorf("%s version response marshal failed", r.RequestURI)
		http.Error(w, "version check failed", 500)
	}

	w.Write(resp)
}
