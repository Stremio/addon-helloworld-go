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


type AddonItem struct {
    Name	string			`json:"name"`
    Types	[]string		`json:"type"`
    IdPrefixes	[]string		`json:"idPrefixes"`
}
