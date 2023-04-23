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

// used for data to template the expiration options available
type homeTemplate struct {
	ExpireHours []int
	UiMaxLength int
}

var ht homeTemplate

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	err := cfg.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = cfg.GetTemplates()
	if err != nil {
		log.Fatal(err)
	}
	notes.RedisInit()

	if cfg.Config.SlackToken != "" && cfg.Config.SlackSigningSecret != "" {
		notes.SlackInit()
	} else {
		log.Info("Slack secrets not defined - disabling slack")
	}

	// create slice to template index.html
	for i := 1; i <= cfg.Config.MaxExpire; i++ {
		ht.ExpireHours = append(ht.ExpireHours, i)
	}
	ht.UiMaxLength = cfg.Config.UiMaxLength
}

func addStandardHeaders(h http.Header) {
	h.Set("Content-Security-Policy", "default-src 'self'; script-src 'self' https://cdnjs.cloudflare.com; style-src 'self' https://cdn.jsdelivr.net")
	h.Set("X-Frame-Options", "DENY")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("X-XSS-Protection", "1; mode=block")
	if cfg.Config.EnableHsts {
		h.Set("Strict-Transport-Security", "max-age=15552000")
	}
}

type tmpnotesHandler func(http.ResponseWriter, *http.Request)

func (th tmpnotesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Anything we want to add for all requests can go in this handler
	if r.Method == "GET" {
		addStandardHeaders(w.Header())
	}
	th(w, r)
}

func main() {
	port := fmt.Sprint(":", cfg.Config.Port)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", serveStatic(fs))
	http.Handle("/new", tmpnotesHandler(notes.AddNote))
	http.Handle("/id/", tmpnotesHandler(notes.GetNote))
	http.Handle("/counts", tmpnotesHandler(notes.GetCounts))
	http.Handle("/healthz", tmpnotesHandler(health.HealthCheck))
	http.Handle("/version", tmpnotesHandler(version.GetVersion))
	if notes.SlackEnabled {
		http.Handle("/slack", tmpnotesHandler(notes.SlackHandler))
		http.Handle("/slack-response", tmpnotesHandler(notes.SlackResponseHandler))
	}
	log.Info("Server listening at ", port)
	http.ListenAndServe(port, nil)
}

func serveStatic(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.RequestURI)
		addStandardHeaders(w.Header())

		// template index.html instead of serving it from the fileserver
		if r.RequestURI == "/" || r.RequestURI == "/index.html" {
			cfg.Tmpl.ExecuteTemplate(w, "index.html", ht)
		} else {
			fs.ServeHTTP(w, r)
		}
	}
}
