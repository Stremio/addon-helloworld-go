package main

import (
	"encoding/json"
//	"fmt"
	"net/http"
	"log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var MANIFEST = Manifest{
	Id:		"org.stremio.helloworld.go",
	Version:	"0.0.1",
	Name:		"Hello World Go Addon",
	Description:	"Sample addon made with gorilla/mux package providing a few public domain movies",
	Types:		[]string{"movie", "series"},
	Catalogs:	[]string{},
	Resources:	[]string{ "stream" },
}

// var Resources = [
// 	Resource{
// 		Name:	"stream"
// 		Types:	[]string{ "movie", "series" }
// 		IdPrefixes:	[]string{ "tt", "hpy"}
// 	},
// 	Resource{
// 		Name:	"catalogs"
// 		Types:	[ "movie", "series"]
// 	}
// ]

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/manifest.json", ManifestHandler)
	r.HandleFunc("/stream/{type}/{id}.json", StreamHandler)
	http.Handle("/", r)

	// CORS configuration
	headersOk := handlers.AllowedHeaders([]string{
		"Content-Type",
		"X-Requested-With",
		"Accept",
		"Accept-Language",
		"Accept-Encoding",
		"Content-Language",
		"Origin",
	})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET"})
	// Listen

	err := http.ListenAndServe("0.0.0.0:3592", handlers.CORS(originsOk, headersOk, methodsOk)(r))

	if err != nil {
		log.Fatalf("Listen: %s", err.Error())
	}
}


func HomeHandler(w http.ResponseWriter, r *http.Request) {
	type jsonObj map[string]interface{}


	jr, _ := json.Marshal(jsonObj{"Path": '/'})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jr)
}

func ManifestHandler(w http.ResponseWriter, r *http.Request) {
	jr, _ := json.Marshal(MANIFEST)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jr)
}

func StreamHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if params["type"] != "movie" && params["type"] != "series" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	jr := `{"streams": [{"title": "Big Buck Bunny Go", "type": "movie", "url": "http://clips.vorwaerts-gmbh.de/big_buck_bunny.mp4" }]}`
	w.Header().Set("Content-Type", "application/json")
	w.Write([] byte(jr))
}
