package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func respondWithError(err error, isServerError bool, w http.ResponseWriter) {
	prefix := "Internal error : "
	if !isServerError {
		prefix = ""
	}

	response := map[string]string{
		"error": fmt.Sprintf("%s%s", prefix, err),
	}
	encoded, _ := json.Marshal(response)
	w.Write(encoded)

	statusCode := 400
	if isServerError {
		statusCode = 500
	}
	w.WriteHeader(statusCode)
}

func respondWithJSON(w http.ResponseWriter, response map[string]any) {
	encoded, err := json.Marshal(response)
	if err != nil {
		respondWithError(err, true, w)
		return
	}

	w.WriteHeader(200)
	w.Write(encoded)
}

func getRequestJSON(w http.ResponseWriter, r *http.Request) (map[string]string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var values map[string]string
	err = json.Unmarshal(body, &values)
	return values, err
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	values, err := getRequestJSON(w, r)
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

	if len(links) == 0 {
		err := errors.New("No file links found")
		respondWithError(err, false, w)
		return
	}

	response := map[string]any{"links": links}
	respondWithJSON(w, response)
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
