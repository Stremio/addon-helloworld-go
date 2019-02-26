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

var movieMap map[string]StreamItem
var seriesMap map[string]StreamItem

func initializeStreamMaps() {
	movieMap = make( map[string]StreamItem)
	seriesMap = make( map[string]StreamItem)

	// Movies
	movieMap["tt0051744"] = StreamItem{ Title: "House on Haunted Hill", InfoHash: "9f86563ce2ed86bbfedd5d3e9f4e55aedd660960" }
	movieMap["tt1254207"] = StreamItem{ Title: "Big Buck Bunny", Url: "http://clips.vorwaerts-gmbh.de/big_buck_bunny.mp4" }
	movieMap["tt0031051"] = StreamItem{ Title: "The Arizona Kid", YtId: "m3BKVSpP80s" }
	movieMap["tt0137523"] = StreamItem{ Title: "Fight Club", ExternalUrl: "https://www.netflix.com/watch/26004747" }

	//Series
	seriesMap["tt0051744:1:1"] = StreamItem{ Title: "Pioneer One", InfoHash: "07a9de9750158471c3302e4e95edb1107f980fa6" }
}

func main() {
	initializeStreamMaps()

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
	stream := StreamItem{}

	if params["type"] == "movie" {
		stream = movieMap[params["id"]]
	} else if params["type"] == "series" {
		stream = seriesMap[params["id"]] // XXX: season, episode
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([] byte(`{"streams": [`))
	streamJson, _ := json.Marshal(stream)
	w.Write(streamJson)
	w.Write([] byte(`]}`))
}
