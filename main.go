package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "encoding/json"
    "html/template"
)

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
    ID int `json:"id"`
    Locations []string `json:"locations"`
    Dates string `json:"dates"`
}

type Dates struct {
    Index []DatesItem `json:"index"`
}

type DatesItem struct {
    ID int `json:"id"`
    Dates []string `json:"dates"`
}

type Relation struct {
    Index []RelationItem `json:"index"`
}

type RelationItem struct {
    ID             int                 `json:"id"`
    DatesLocations map[string][]string `json:"datesLocations"` //revoir les maps (pas bien compris "map[string][]string")
}

var tmplHome    = template.Must(template.ParseFiles("templates/home.html"))
var tmplArtists = template.Must(template.ParseFiles("templates/artists.html"))
var tmplLocations = template.Must(template.ParseFiles("templates/locations.html"))
var tmplDates = template.Must(template.ParseFiles("templates/dates.html"))
var tmplRelation = template.Must(template.ParseFiles("templates/relation.html")) 

func homeHandler(w http.ResponseWriter, r *http.Request) {
    tmplHome.Execute(w, nil)
}

func artistsHandler(w http.ResponseWriter, r *http.Request) {
    url := "https://groupietrackers.herokuapp.com/api/artists"

    response, err := http.Get(url)
    if err != nil {
        http.Error(w, "Erreur lors de la requête HTTP", http.StatusInternalServerError)
        return
    }
    defer response.Body.Close()

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        http.Error(w, "Erreur lors de la lecture de la réponse", http.StatusInternalServerError)
        return
    }

    var artists []Artist
    if err := json.Unmarshal(body, &artists); err != nil {
        http.Error(w, "Erreur de décodage JSON", http.StatusInternalServerError)
        return
    }

    if err := tmplArtists.Execute(w, artists); err != nil {
        http.Error(w, "Erreur d’exécution du template", http.StatusInternalServerError)
        return
    }
}

func locationsHandler(w http.ResponseWriter, r *http.Request) {
        url := "https://groupietrackers.herokuapp.com/api/locations"   

    response, err := http.Get(url)
    if err != nil {
        http.Error(w, "Erreur lors de la requête HTTP", http.StatusInternalServerError)
        return
    }
    defer response.Body.Close()
    
    body, err := ioutil.ReadAll(response.Body)
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

func datesHandler(w http.ResponseWriter, r *http.Request) {
url := "https://groupietrackers.herokuapp.com/api/dates"

    response, err := http.Get(url)
    if err != nil {
        http.Error(w, "Erreur lors de la requête HTTP", http.StatusInternalServerError)
        return
    }
    defer response.Body.Close()

    body, err := ioutil.ReadAll(response.Body)
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

func relationHandler(w http.ResponseWriter, r *http.Request) {
    url := "https://groupietrackers.herokuapp.com/api/relation"

        response, err := http.Get(url)
    if err != nil {
        http.Error(w, "Erreur lors de la requête HTTP", http.StatusInternalServerError)
        return
    }
    defer response.Body.Close()

        body, err := ioutil.ReadAll(response.Body)
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


func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/artists", artistsHandler)
    http.HandleFunc("/locations", locationsHandler)
    http.HandleFunc("/dates", datesHandler)
    http.HandleFunc("/relation", relationHandler)

    fmt.Println("Serveur Go : http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
