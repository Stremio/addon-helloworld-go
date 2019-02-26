package main

type Manifest struct {
    Id		string		`json:"id"`
    Version	string		`json:"version"`
    Name	string		`json:"name"`
    Description	string		`json:"description"`
    Types	[]string	`json:"types"`
    Catalogs	[]string	`json:"catalogs"`
    Resources	[]string	`json:"resources"`
}

type Resource struct {
    Name	string			`json:"name"`
    Types	[]string		`json:"type"`
    IdPrefixes	[]string		`json:"idPrefixes,omitempty"`
}

type StreamItemType uint8

const (
    MOVIE   StreamItemType = 0
    SERIES  StreamItemType = 1
)

type StreamItem struct {
    Title	string			`json:"title"`
//     Type	StreamItemType		`json:"-"`
    InfoHash	string			`json:"infoHash,omitempty"`
    FileIdx	uint8			`json:"fileIdx,omitempty"`
    Url		string			`json:"url,omitempty"`
    YtId	string			`json:"ytId,omitempty"`
    ExternalUrl	string			`json:"externalUrl,omitempty"`
}
