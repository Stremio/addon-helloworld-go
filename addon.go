package main

import (
	//"fmt"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	stremio "github.com/Stremio/addon-helloworld-go/internal/stremio/types"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var CATALOG_ID = "Hello, Go"

var MANIFEST = stremio.Manifest{
	Id:          "org.stremio.helloworld.go",
	Version:     "0.0.1",
	Name:        "Hello World Go Addon",
	Description: "Sample addon made with gorilla/mux package providing a few public domain movies",
	Types:       []string{"movie", "series"},
	Catalogs:    []stremio.CatalogItem{},
	Resources:   []string{"stream", "catalog"},
}

var movieMap map[string]stremio.StreamItem
var seriesMap map[string]stremio.StreamItem

var movieMetaMap map[string]stremio.MetaItem
var seriesMetaMap map[string]stremio.MetaItem

var METAHUB_BASE_URL = "https://images.metahub.space/poster/medium/"

func init() {
	movieMap = make(map[string]stremio.StreamItem)
	seriesMap = make(map[string]stremio.StreamItem)

	// Movies
	movieMap["tt0051744"] = stremio.StreamItem{Title: "House on Haunted Hill",
		InfoHash: "9f86563ce2ed86bbfedd5d3e9f4e55aedd660960"}
	movieMap["tt1254207"] = stremio.StreamItem{Title: "Big Buck Bunny",
		Url: "http://clips.vorwaerts-gmbh.de/big_buck_bunny.mp4"}
	movieMap["tt0031051"] = stremio.StreamItem{Title: "The Arizona Kid", YtId: "m3BKVSpP80s"}
	movieMap["tt0137523"] = stremio.StreamItem{Title: "Fight Club",
		ExternalUrl: "https://www.netflix.com/watch/26004747"}

	//Series
	seriesMap["tt1748166"] = stremio.StreamItem{Title: "Pioneer One",
		InfoHash: "07a9de9750158471c3302e4e95edb1107f980fa6"}

	// Meta
	movieMetaMap = make(map[string]stremio.MetaItem)
	seriesMetaMap = make(map[string]stremio.MetaItem)

	movieMetaMap["tt0051744"] = stremio.MetaItem{Name: "House on Haunted Hill",
		Genres: []string{"Horror", "Mystery"}}
	movieMetaMap["tt1254207"] = stremio.MetaItem{Name: "Big Buck Bunny", Genres: []string{"Animation", "Short", "Comedy"},
		Poster: "https://peach.blender.org/wp-content/uploads/poster_bunny_small.jpg"}
	movieMetaMap["tt0031051"] = stremio.MetaItem{Name: "The Arizona Kid",
		Genres: []string{"Music", "War", "Western"}}
	movieMetaMap["tt0137523"] = stremio.MetaItem{Name: "Fight Club",
		Genres: []string{"Drama"}}

	//Series
	seriesMetaMap["tt1748166"] = stremio.MetaItem{Name: "Pioneer One",
		Genres: []string{"Drama"}}
}

func main() {

	MANIFEST.Catalogs = append(MANIFEST.Catalogs, stremio.CatalogItem{"movie", CATALOG_ID})
	MANIFEST.Catalogs = append(MANIFEST.Catalogs, stremio.CatalogItem{"series", CATALOG_ID})

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/manifest.json", ManifestHandler)
	r.HandleFunc("/stream/{type}/{id}.json", StreamHandler)
	r.HandleFunc("/catalog/{type}/{id}.json", CatalogHandler)
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
	log.Println("Listening: 0.0.0.0:3593")
	err := http.ListenAndServe("0.0.0.0:3593", handlers.CORS(originsOk, headersOk, methodsOk)(r))
	if err != nil {
		log.Fatalf("Listen: %s", err.Error())
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	jr, _ := json.Marshal(stremio.JsonObj{"Path": '/'})
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
	stream := stremio.StreamItem{}

	if params["type"] == "movie" {
		stream = movieMap[params["id"]]
	} else if params["type"] == "series" {
		itemIds := strings.Split(params["id"], ":")
		showID, seasonId, episodeId := itemIds[0], itemIds[1], itemIds[2]
		stream = seriesMap[showID] // XXX: season, episode
		// silence the compiler
		if seasonId+episodeId != string(stream.FileIdx) {
			log.Println("Return stream for episode 1")
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"streams": [`))
	streamJson, _ := json.Marshal(stream)
	w.Write(streamJson)
	w.Write([]byte(`]}`))
}

func CatalogHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	metaMap := make(map[string]stremio.MetaItem)

	for _, item := range MANIFEST.Catalogs {
		if params["id"] == item.Id && params["type"] == item.Type {
			switch item.Type {
			case "series":
				metaMap = seriesMetaMap
			case "movie":
				metaMap = movieMetaMap
			default:
				continue
			}
			break
		}
	}

	if len(metaMap) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metas := []stremio.MetaItemJson{}
	for metaKey, metaValue := range metaMap {
		item := stremio.MetaItemJson{
			Id:     metaKey,
			Type:   params["type"],
			Name:   metaValue.Name,
			Genres: metaValue.Genres,
			Poster: METAHUB_BASE_URL + metaKey + "/img",
		}
		if metaValue.Poster != "" {
			item.Poster = metaValue.Poster
		}
		metas = append(metas, item)
	}

	w.Header().Set("Content-Type", "application/json")
	catalogJson, _ := json.Marshal(&stremio.JsonObj{"metas": metas})
	w.Write(catalogJson)
}
