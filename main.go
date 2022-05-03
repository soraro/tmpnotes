package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	cfg "tmpnotes/internal/config"
	h "tmpnotes/internal/headers"
	"tmpnotes/internal/health"
	"tmpnotes/internal/notes"
	"tmpnotes/internal/version"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	err := cfg.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	notes.RedisInit()
}

func main() {
	port := fmt.Sprint(":", cfg.Config.Port)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", addHeaders(fs))
	http.HandleFunc("/new", notes.AddNote)
	http.HandleFunc("/id/", notes.GetNote)
	http.HandleFunc("/counts", notes.GetCounts)
	http.HandleFunc("/healthz", health.HealthCheck)
	http.HandleFunc("/version", version.GetVersion)
	log.Info("Server listening at ", port)
	http.ListenAndServe(port, nil)
}

func addHeaders(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.RequestURI)
		h.AddStandardHeaders(w.Header())
		fs.ServeHTTP(w, r)
	}
}
