package notes

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type counts struct {
	NoteCount        int `json:"noteCount" redis:"noteCount,omitempty"`
	EncNoteCount     int `json:"encNoteCount" redis:"encNoteCount,omitempty"`
	CurrentNoteCount int `json:"currentNoteCount" redis:"-"`
}

func GetCounts(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Invalid Request", 400)
		return
	}
	log.Info(r.RequestURI)

	var c counts

	w.Header().Set("Content-Type", "application/json")

	pipe := rdb.Pipeline()
	countdata := pipe.HGetAll(ctx, "counts")
	curnotes := pipe.Do(ctx, "DBSIZE")
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Errorf("Error getting counts: %s", err)
	}

	c.CurrentNoteCount, _ = curnotes.Int()
	c.CurrentNoteCount -= 1 // subtract 1 to omit the "counts" key
	countdata.Scan(&c)
	resp, err := json.Marshal(c)
	if err != nil {
		log.Errorf("Error in JSON marshal. Err: %s", err)
	}

	w.Write(resp)
}
