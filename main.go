package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
	"math"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var tpl = template.Must(template.ParseGlob("templates/*.html"))

// Message dsadadasd
type Message struct {
	Message string    `datastore:"message"`
	Source  string    `datastore:"source"`
	Data    string    `datastore:"data,noindex"`
	File		string		`datastore:"file"`
	When    time.Time `datastore:"when"`
}

func JSONDecode(reader io.Reader, target interface{}) {
	decoder := json.NewDecoder(reader)
	decoder.Decode(target)
}

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var value Message
		JSONDecode(r.Body, &value)

		if value.Message == "" {
			return
		}

		messageLen := math.Min(float64(len(value.Message)), float64(1500))
		value.Message = value.Message[0: int(messageLen)]

		c := appengine.NewContext(r)
		k := datastore.NewIncompleteKey(c, "messages", nil)
		if _, err := datastore.Put(c, k, &value); err != nil {
			log.Fatal(err.Error())
		} else {
			fmt.Fprint(w, "store ok")
		}
		break
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	// c := appengine.NewContext(r)
	//
	// q := datastore.NewQuery("messages").Limit(10)
	//
	// messages := make([]Message, 0, 10)
	//
	// if _, err := q.GetAll(c, &messages); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	//
	// if err := tpl.ExecuteTemplate(w, "home.html", messages); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }
	fmt.Fprint(w, "home ok")
}

func Init() {
	router := mux.NewRouter().StrictSlash(false)
	s := router.PathPrefix("/api").Subrouter()

	s.HandleFunc("/messages", ApiHandler).Methods("GET", "POST")

	http.Handle("/", cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		MaxAge: 0,
	}).Handler(router))
}

func init() {
	Init()
}

func main() {
	Init()
}
