package notes

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	cfg "tmpnotes/internal/config"
	"tmpnotes/internal/crypto"
	h "tmpnotes/internal/headers"
)

const maxLength = 1000
const maxExpire = 24

var (
	ctx = context.Background()
	rdb *redis.Client
)

type note struct {
	Message string `json:"message"`
	Expire  int    `json:"expire"`
}

func RedisInit() {
	log.SetFormatter(&log.JSONFormatter{})
	opt, err := redis.ParseURL(cfg.Config.RedisUrl)
	if err != nil {
		panic(err)
	}

	rdb = redis.NewClient(opt)
}

func AddNote(w http.ResponseWriter, r *http.Request) {

	log.Info(r.RequestURI)

	if r.Method != "POST" {
		log.Errorf("%s Invalid request method: %s", r.RequestURI, r.Method)
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var n note
	err := json.NewDecoder(r.Body).Decode(&n)
	if err != nil {
		http.Error(w, "Invalid Request", 400)
		return
	}
	if n.Expire < 1 {
		// Avoid having no expiration time
		n.Expire = 1
	}

	if n.Expire > maxExpire {
		log.Errorf("%s Expiration set too high: %v", r.RequestURI, n.Expire)
		http.Error(w, "Invalid Request - TTL is too high", 400)
		return
	}

	if !checkAcceptableLength(n.Message) {
		log.Errorf("%s Message size too large: %v", r.RequestURI, len(n.Message))
		http.Error(w, "Invalid Request", 400)
		return
	}
	id, key := generateIdAndKey()

	encryptedMessage, err := encryptNote(n.Message, key)
	if err != nil {
		log.Errorf("%s Issue encrypting message: %s", r.RequestURI, err)
	}
	pipe := rdb.Pipeline()
	pipe.Set(ctx, id, encryptedMessage, time.Duration(n.Expire)*time.Hour)
	pipe.HIncrBy(ctx, "counts", noteType(n.Message), 1)
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Errorf("%s Error setting note values: %s", r.RequestURI, err)
	}

	fmt.Fprint(w, id+key)

}

// The first 8 characters will be the redis key (id)
// The last 24 characters will be the key to encrypt the note
func generateIdAndKey() (string, string) {
	uuid, _ := uuid.NewRandom()
	full := strings.ReplaceAll(uuid.String(), "-", "")
	return full[:8], full[8:]
}

func encryptNote(note, key string) (string, error) {
	out, err := crypto.Encrypt([]byte(note), keyBytes(key))
	return hex.EncodeToString(out), err
}

// hash the key and take the first 32 characters
func keyBytes(key string) *[32]byte {
	hash := crypto.Hash("TMPNOTES", []byte(key))
	hashedKey := hex.EncodeToString(hash)
	r := strings.NewReader(hashedKey[:32])
	kb := [32]byte{}
	_, err := io.ReadFull(r, kb[:])
	if err != nil {
		panic(err)
	}
	return &kb
}

func checkAcceptableLength(m string) bool {
	return len(m) <= maxLength
}

// return the type of note from the first 5 characters
func noteType(note string) string {
	if len(note) < 5 {
		return "noteCount"
	}
	if note[0:5] == "[ENC]" {
		return "encNoteCount"
	} else {
		return "noteCount"
	}

}

func GetNote(w http.ResponseWriter, r *http.Request) {

	full := strings.ReplaceAll(r.RequestURI, "/id/", "")
	id := full[:8]
	key := full[8:]
	log.Info(id)
	h.AddStandardHeaders(w.Header())

	if r.Method != "GET" {
		log.Errorf("%s Invalid request method: %s", id, r.Method)
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	val, err := rdb.Get(ctx, id).Result()
	switch {
	case err == redis.Nil:
		w.WriteHeader(404)
		if textResponse(r.UserAgent()) {
			fmt.Fprintf(w, "404 - Nothing to see here\n")
		} else {
			t, _ := template.ParseFiles("./templates/404.html")
			t.Execute(w, nil)
		}
		return
	case err != nil:
		log.Errorf("%s Redis GET failed: %s", id, err)
		http.Error(w, "Server Error", 500)
		return
	case val == "":
		w.WriteHeader(404)
		if textResponse(r.UserAgent()) {
			fmt.Fprintf(w, "404 - Nothing to see here")
		} else {
			t, _ := template.ParseFiles("./templates/404.html")
			t.Execute(w, nil)
		}
		return
	}

	if returnData(r.UserAgent(), r.Header.Get("X-Note")) {
		rdb.Del(ctx, id)
		decNote, err := decryptNote(val, key)
		if err != nil {
			http.Error(w, "👀 401 Unauthorized", 401)
			log.Errorf("Failed decrypting: %s", err)
		}
		// add a newline for text clients so your prompt wont start in the middle of the line
		if textResponse(r.UserAgent()) {
			fmt.Fprint(w, decNote+"\n")
		} else {
			fmt.Fprint(w, decNote)
		}
		return
	}

	t, err := template.ParseFiles("./templates/note.html")
	if err != nil {
		log.Errorf("%s Error rendering note: %s", id, err)
		http.Error(w, "Error rendering note", 500)
	}
	t.Execute(w, nil)
}

// Check headers to see if we should return the data or not.
// This helps make it so various link previews won't instantly burn the note
func returnData(useragent, header string) bool {
	if textResponse(useragent) {
		return true
	} else {
		// A predictable header we can use to signal the note can be returned/destroyed
		return header == "Destroy"
	}
}

// Check the user-agent for to see if we should return a text response
func textResponse(useragent string) bool {
	// add other user agents here that will burn the note right away
	acceptedUserAgents := []string{"curl", "wget"}

	for _, v := range acceptedUserAgents {
		if strings.Contains(strings.ToLower(useragent), v) {
			return true
		}
	}
	return false
}

func decryptNote(ciphernote, key string) (string, error) {
	ct, err := hex.DecodeString(ciphernote)
	if err != nil {
		return "", err
	}
	out, err := crypto.Decrypt([]byte(ct), keyBytes(key))
	if err != nil {
		return "", err
	}
	return string(out), nil
}
