package main

import (
	"fmt"
	"net/http"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {}

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
