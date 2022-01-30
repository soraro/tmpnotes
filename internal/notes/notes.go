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
	var n note

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
		t, _ := template.ParseFiles("./templates/404.html")
		t.Execute(w, n)
		return
	case err != nil:
		fmt.Println("Get failed", err)
		http.Error(w, "Server Error", 500)
		return
	case val == "":
		w.WriteHeader(404)
		t, _ := template.ParseFiles("./templates/404.html")
		t.Execute(w, n)
		return
	}
	n.Message = val

	rdb.Del(ctx, id)

	t, err := template.ParseFiles("./templates/note.html")
	if err != nil {
		http.Error(w, "Error rendering note", 500)
	}
	t.Execute(w, n)
}
