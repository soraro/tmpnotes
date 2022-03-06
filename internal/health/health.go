package health

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {

	log.Info(r.RequestURI)

	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp := make(map[string]string)
		resp["message"] = "200 Status OK"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Errorf("%s healthcheck response format failed", r.RequestURI)
			http.Error(w, "healthcheck failed", 500)
		}
		w.Write(jsonResp)
		return
	} else {
		log.Errorf("%s healthcheck failed: %s", r.RequestURI, r.Method)
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := make(map[string]string)
		resp["message"] = "405 Method Not Allowed"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Errorf("%s healthcheck response format failed", r.RequestURI)
			http.Error(w, "healthcheck failed", 500)
		}
		w.Write(jsonResp)
		return
	}

}
