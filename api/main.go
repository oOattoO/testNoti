package main

import (
	"encoding/json"
	// "io/ioutil"
	"log"
	"net/http"
	"os"
	// "time"
	"fmt"
	// "strconv"

	"github.com/rs/cors"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Post struct {
	TextID      string    `json:"textId" bson:"textId"`
	Tag string `json:tag bson:"tag"`
	Title string `json:title bson:"title"`
}

var posts *mgo.Collection

func main() {
	// Connect to mongo
	session, err := mgo.Dial("mongo:27017")
	if err != nil {
		log.Fatalln(err)
		log.Fatalln("mongo err")
		os.Exit(1)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	// Get posts collection
	posts = session.DB("app").C("posts")

	// Set up routes
	r := mux.NewRouter()

	r.HandleFunc("/posts/{id}", readPosts).
		Methods("GET")

	http.ListenAndServe(":8080", cors.AllowAll().Handler(r))
	log.Println("Listening on port 8080...")
}


func readPosts(w http.ResponseWriter, r *http.Request) {
	result := Post{}
	eventID := mux.Vars(r)["id"]
	fmt.Println(eventID)
	// Find(bson.M{"text": "test"}).One(&testting)
	if err := posts.Find(bson.M{"textId": eventID}).One(&result); err != nil {
		fmt.Println(result)
		responseError(w, err.Error(), http.StatusInternalServerError)
	} else {
		responseJSON(w, result)
	}
}

func responseError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
