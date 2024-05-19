package main

import "fmt"

func main() {
	query := "shakespear"
	links, err := GetVKDownloadLinks(query)
	if err != nil {
		panic(err)
	}

	for _, link := range links {
		fmt.Println(link)
	}
}
