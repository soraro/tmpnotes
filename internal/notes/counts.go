package notes

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	h "tmpnotes/internal/headers"
)

type counts struct {
	NoteCount        int `json:"noteCount" redis:"noteCount,omitempty"`
	EncNoteCount     int `json:"encNoteCount" redis:"encNoteCount,omitempty"`
	CurrentNoteCount int `json:"currentNoteCount" redis:"-"`
}

func GetCounts(w http.ResponseWriter, r *http.Request) {
	h.AddStandardHeaders(w.Header())
	log.Info(r.RequestURI)

	if r.Method != "GET" {
		log.Errorf("%s Invalid request method: %s", r.RequestURI, r.Method)
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var c counts

	w.Header().Set("Content-Type", "application/json")

	pipe := rdb.Pipeline()
	countdata := pipe.HGetAll(ctx, "counts")
	curnotes := pipe.Do(ctx, "DBSIZE")
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Errorf("%s Error getting counts: %s", r.RequestURI, err)
	}

	c.CurrentNoteCount, _ = curnotes.Int()
	c.CurrentNoteCount -= 1 // subtract 1 to omit the "counts" key
	countdata.Scan(&c)
	resp, err := json.Marshal(c)
	if err != nil {
		log.Errorf("%s Error in JSON marshal. Err: %s", r.RequestURI, err)
	}

	w.Write(resp)
}
