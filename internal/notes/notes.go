package notes

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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

func init() {
	var redisConnectionString string
	if os.Getenv("REDIS_URL") == "" {
		redisConnectionString = "redis://localhost:6379"
	} else {
		redisConnectionString = os.Getenv("REDIS_URL")
	}
	log.SetFormatter(&log.JSONFormatter{})
	opt, err := redis.ParseURL(redisConnectionString)
	if err != nil {
		panic(err)
	}

	rdb = redis.NewClient(opt)
}

func AddNote(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Invalid Request", 400)
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
		http.Error(w, "Invalid Request - TTL is too high", 400)
		return
	}

	if !checkAcceptableLength(n.Message) {
		log.Errorf("Message size too large: %v", len(n.Message))
		http.Error(w, "Invalid Request", 400)
		return
	}
	uuid := getId()
	rdb.Set(ctx, uuid, n.Message, time.Duration(n.Expire)*time.Hour)
	fmt.Fprint(w, uuid)

}

func getId() string {
	uuid := uuid.NewString()
	return strings.ReplaceAll(uuid, "-", "")
}

func checkAcceptableLength(m string) bool {
	return len(m) <= maxLength
}

func GetNote(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Invalid Request", 400)
		return
	}
	log.Info(r.RequestURI)
	id := strings.ReplaceAll(r.RequestURI, "/id/", "")

	val, err := rdb.Get(ctx, id).Result()
	switch {
	case err == redis.Nil:
		w.WriteHeader(404)
		if textResponse(r.UserAgent()) {
			fmt.Fprintf(w, "404 - Nothing to see here")
		} else {
			t, _ := template.ParseFiles("./templates/404.html")
			t.Execute(w, nil)
		}
		return
	case err != nil:
		fmt.Println("Get failed", err)
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
		fmt.Fprint(w, val)
		return
	}

	//rdb.Del(ctx, id)
	t, err := template.ParseFiles("./templates/note.html")
	if err != nil {
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
