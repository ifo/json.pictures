package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		setJSON(w, r)
	}
}

func setJSON(w http.ResponseWriter, r *http.Request) {
	// For now, ensure that the request wasn't larger than 2048 bytes.
	if r.ContentLength == -1 || r.ContentLength > 2048 {
		http.Error(w,
			fmt.Sprintf("Content length is %d, but must be smaller than 2048 bytes", r.ContentLength), 400)
		return
	}

	bts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "", 500)
		return
	}

	// JSON must be valid in order to be indented.
	// TODO: Find a better way to validate the JSON.
	var ignore bytes.Buffer
	err = json.Indent(&ignore, bts, "", "")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	PathResponseCache[r.URL.Path] = JSONResponse{JSON: bts}
}

func returnJSON(w http.ResponseWriter, r *http.Request) {
	jsonResp, ok := PathResponseCache[r.URL.Path]
	if !ok {
		http.Error(w, "", 404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp.JSON)
}
