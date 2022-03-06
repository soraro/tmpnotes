package main

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"tmpnotes/internal/health"
	"tmpnotes/internal/notes"
)

func main() {
	var port string
	log.SetFormatter(&log.JSONFormatter{})
	if os.Getenv("PORT") == "" {
		port = ":5000"
	} else {
		port = ":" + os.Getenv("PORT")
	}
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/new", notes.AddNote)
	http.HandleFunc("/id/", notes.GetNote)
	http.HandleFunc("/counts", notes.GetCounts)
	http.HandleFunc("/healthz", health.HealthCheck)
	log.Info("Server listening at ", port)
	http.ListenAndServe(port, nil)
}
