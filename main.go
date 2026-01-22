package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

// ====== Structures de données ======

type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Locations struct {
	Index []LocationItem `json:"index"`
}

type LocationItem struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type Dates struct {
	Index []DatesItem `json:"index"`
}

type DatesItem struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	Index []RelationItem `json:"index"`
}

type RelationItem struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type ArtistPageData struct {
	Artist Artist
	Locs   LocationItem
	Dates  DatesItem
	Rel    RelationItem
}

// ====== Templates ======

var (
    tmplHome = template.Must(template.New("home").Funcs(map[string]interface{}{
        "safeJS": func(s string) template.JS { return template.JS(s) },
    }).ParseFiles("templates/home.html"))
    tmplArtists  = template.Must(template.ParseFiles("templates/artists.html"))
    tmplArtist   = template.Must(template.ParseFiles("templates/artist.html"))
    tmplLocations = template.Must(template.ParseFiles("templates/locations.html"))
    tmplDates    = template.Must(template.ParseFiles("templates/dates.html"))
    tmplRelation = template.Must(template.ParseFiles("templates/relation.html"))
)

// ====== Handlers ======

// Accueil : autocomplete
func homeHandler(w http.ResponseWriter, r *http.Request) {
	url := "https://groupietrackers.herokuapp.com/api/artists"

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Erreur lors de la requête HTTP", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture de la réponse", http.StatusInternalServerError)
		return
	}

	var artists []Artist
	if err := json.Unmarshal(body, &artists); err != nil {
		http.Error(w, "Erreur de décodage JSON", http.StatusInternalServerError)
		return
	}

	if err := tmplHome.Execute(w, artists); err != nil {
		http.Error(w, "Erreur d’exécution du template", http.StatusInternalServerError)
		return
	}
}

// Liste artistes + filtre ?q=
func artistsHandler(w http.ResponseWriter, r *http.Request) {
	url := "https://groupietrackers.herokuapp.com/api/artists"

	response, err := http.Get(url)
	if err != nil {
		http.Error(w, "Erreur lors de la requête HTTP", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture de la réponse", http.StatusInternalServerError)
		return
	}

	var artists []Artist
	if err := json.Unmarshal(body, &artists); err != nil {
		http.Error(w, "Erreur de décodage JSON", http.StatusInternalServerError)
		return
	}

	query := r.URL.Query().Get("q")
	if query != "" {
		query = strings.ToLower(query)
		var filtered []Artist

		for _, a := range artists {
			if strings.Contains(strings.ToLower(a.Name), query) {
				filtered = append(filtered, a)
				continue
			}
			for _, m := range a.Members {
				if strings.Contains(strings.ToLower(m), query) {
					filtered = append(filtered, a)
					break
				}
			}
		}
		artists = filtered
	}

	if err := tmplArtists.Execute(w, artists); err != nil {
		http.Error(w, "Erreur d’exécution du template", http.StatusInternalServerError)
		return
	}
}

// Fiche artiste complète : infos + lieux + dates + relations
func artistHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}

	// 1) Artiste
	artistURL := "https://groupietrackers.herokuapp.com/api/artists/" + idStr
	respArtist, err := http.Get(artistURL)
	if err != nil {
		http.Error(w, "Erreur lors de la requête artiste", http.StatusInternalServerError)
		return
	}
	defer respArtist.Body.Close()

	bodyArtist, err := io.ReadAll(respArtist.Body)
	if err != nil {
		http.Error(w, "Erreur de lecture artiste", http.StatusInternalServerError)
		return
	}

	var artist Artist
	if err := json.Unmarshal(bodyArtist, &artist); err != nil {
		http.Error(w, "Erreur JSON artiste", http.StatusInternalServerError)
		return
	}

	// 2) Locations
	locURL := "https://groupietrackers.herokuapp.com/api/locations"
	respLoc, err := http.Get(locURL)
	if err != nil {
		http.Error(w, "Erreur lors de la requête locations", http.StatusInternalServerError)
		return
	}
	defer respLoc.Body.Close()

	bodyLoc, err := io.ReadAll(respLoc.Body)
	if err != nil {
		http.Error(w, "Erreur de lecture locations", http.StatusInternalServerError)
		return
	}

	var locs Locations
	if err := json.Unmarshal(bodyLoc, &locs); err != nil {
		http.Error(w, "Erreur JSON locations", http.StatusInternalServerError)
		return
	}

	var locItem LocationItem
	for _, item := range locs.Index {
		if item.ID == artist.ID {
			locItem = item
			break
		}
	}

	// 3) Dates
	datesURL := "https://groupietrackers.herokuapp.com/api/dates"
	respDates, err := http.Get(datesURL)
	if err != nil {
		http.Error(w, "Erreur lors de la requête dates", http.StatusInternalServerError)
		return
	}
	defer respDates.Body.Close()

	bodyDates, err := io.ReadAll(respDates.Body)
	if err != nil {
		http.Error(w, "Erreur de lecture dates", http.StatusInternalServerError)
		return
	}

	var dates Dates
	if err := json.Unmarshal(bodyDates, &dates); err != nil {
		http.Error(w, "Erreur JSON dates", http.StatusInternalServerError)
		return
	}

	var datesItem DatesItem
	for _, item := range dates.Index {
		if item.ID == artist.ID {
			datesItem = item
			break
		}
	}

	// 4) Relations (lieu -> dates)
	relURL := "https://groupietrackers.herokuapp.com/api/relation"
	respRel, err := http.Get(relURL)
	if err != nil {
		http.Error(w, "Erreur lors de la requête relations", http.StatusInternalServerError)
		return
	}
	defer respRel.Body.Close()

	bodyRel, err := io.ReadAll(respRel.Body)
	if err != nil {
		http.Error(w, "Erreur de lecture relations", http.StatusInternalServerError)
		return
	}

	var rel Relation
	if err := json.Unmarshal(bodyRel, &rel); err != nil {
		http.Error(w, "Erreur JSON relations", http.StatusInternalServerError)
		return
	}

	var relItem RelationItem
	for _, item := range rel.Index {
		if item.ID == artist.ID {
			relItem = item
			break
		}
	}

	data := ArtistPageData{
		Artist: artist,
		Locs:   locItem,
		Dates:  datesItem,
		Rel:    relItem,
	}

	if err := tmplArtist.Execute(w, data); err != nil {
		http.Error(w, "Erreur d’exécution du template", http.StatusInternalServerError)
		return
	}
}

// Lieux
func locationsHandler(w http.ResponseWriter, r *http.Request) {
	url := "https://groupietrackers.herokuapp.com/api/locations"

	response, err := http.Get(url)
	if err != nil {
		http.Error(w, "Erreur lors de la requête HTTP", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture de la réponse", http.StatusInternalServerError)
		return
	}

	var data Locations
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Erreur de décodage JSON", http.StatusInternalServerError)
		return
	}

	if err := tmplLocations.Execute(w, data); err != nil {
		http.Error(w, "Erreur d’exécution du template", http.StatusInternalServerError)
		return
	}
}

// Dates
func datesHandler(w http.ResponseWriter, r *http.Request) {
	url := "https://groupietrackers.herokuapp.com/api/dates"

	response, err := http.Get(url)
	if err != nil {
		http.Error(w, "Erreur lors de la requête HTTP", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture de la réponse", http.StatusInternalServerError)
		return
	}

	var data Dates
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Erreur de décodage JSON", http.StatusInternalServerError)
		return
	}

	if err := tmplDates.Execute(w, data); err != nil {
		http.Error(w, "Erreur d’exécution du template", http.StatusInternalServerError)
		return
	}
}

// Relations (page globale)
func relationHandler(w http.ResponseWriter, r *http.Request) {
	url := "https://groupietrackers.herokuapp.com/api/relation"

	response, err := http.Get(url)
	if err != nil {
		http.Error(w, "Erreur lors de la requête HTTP", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture de la réponse", http.StatusInternalServerError)
		return
	}

	var data Relation
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Erreur de décodage JSON", http.StatusInternalServerError)
		return
	}

	if err := tmplRelation.Execute(w, data); err != nil {
		http.Error(w, "Erreur d’exécution du template", http.StatusInternalServerError)
		return
	}
}

// ====== main ======

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/artists", artistsHandler)
	http.HandleFunc("/artist", artistHandler)
	http.HandleFunc("/locations", locationsHandler)
	http.HandleFunc("/dates", datesHandler)
	http.HandleFunc("/relation", relationHandler)

	fmt.Println("Serveur Go : http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
