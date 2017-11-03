package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

var PathResponseCache = map[string]JSONResponse{}

// TODO: Add more meta features, like:
// Different Indentation, to show the response in more human readable formats
// Delay duration, for testing slower apis
// Max RPS, with designated failure response
type JSONResponse struct {
	JSON []byte
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.Write([]byte("Welcome to json.pictures, you're one stop shop for JSON api mocking\nJust post your JSON to a url and you're good to go."))
		return
	}

	if r.Method == http.MethodGet {
		returnJSON(w, r)
	} else if r.Method == http.MethodPost {
		// TODO: check the body size and error on sizes above some limit (10k? 100k?)
		setJSON(w, r)
	}
}

// Post requests come here
func setJSON(w http.ResponseWriter, r *http.Request) {
	bts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "", 500)
		return
	}

	// Validate JSON
	// In order to be indented, JSON must be valid.
	// TODO: Find a better way to validate the JSON.
	var ignore bytes.Buffer
	err = json.Indent(&ignore, bts, "", "")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Save it
	PathResponseCache[r.URL.Path] = JSONResponse{JSON: bts}
}

// Get requests come here
func returnJSON(w http.ResponseWriter, r *http.Request) {
	jsonResp, ok := PathResponseCache[r.URL.Path]
	if !ok {
		http.Error(w, "", 404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp.JSON)
}
