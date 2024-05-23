package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func respondWithError(err error, isServerError bool, w http.ResponseWriter) {
	response := map[string]string{
		"error": fmt.Sprintf("Internal error: %s", err),
	}
	encoded, _ := json.Marshal(response)
	w.Write(encoded)

	statusCode := 400
	if isServerError {
		statusCode = 500
	}
	w.WriteHeader(statusCode)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Get the query parameter from the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(err, true, w)
		return
	}

	var values map[string]string
	err = json.Unmarshal(body, &values)
	if err != nil {
		respondWithError(err, true, w)
		return
	}

	query, ok := values["query"]
	if !ok {
		respondWithError(errors.New("query parameter not found"), false, w)
		return
	}

	// Respond with links
	vk := VK{}
	links, err := GetFileLinks(query, vk)
	if err != nil {
		respondWithError(err, true, w)
		return
	}

	response := map[string][]Link{"links": links}
	encoded, err := json.Marshal(response)
	if err != nil {
		respondWithError(err, true, w)
		return
	}

	w.WriteHeader(200)
	w.Write(encoded)
}

func main() {
	server := http.FileServer(http.Dir("web"))
	http.Handle("/", server)
	http.HandleFunc("/search", searchHandler)

	fmt.Println("Running the server at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
