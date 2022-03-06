package health

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {

	log.Info(r.RequestURI)

	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "applications/json")
		resp := make(map[string]string)
		resp["message"] = "Status OK"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Error("healthcheck response format failed")
			http.Error(w, "healthcheck failed", 500)
		}
		w.Write(jsonResp)
		return
	} else {
		log.Errorf("%s healthcheck failed: %s", r.RequestURI, r.Method)
		http.Error(w, "404 Not Found", 404)
		return
	}

}
