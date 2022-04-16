package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	cfg "tmpnotes/internal/config"
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

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/new", notes.AddNote)
	http.HandleFunc("/id/", notes.GetNote)
	http.HandleFunc("/counts", notes.GetCounts)
	http.HandleFunc("/healthz", health.HealthCheck)
	http.HandleFunc("/version", version.GetVersion)
	log.Info("Server listening at ", port)
	http.ListenAndServe(port, nil)
}
