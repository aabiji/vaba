package main

import "fmt"

func main() {
	query := "A midsummer's dream"
	fmtQuery := fmt.Sprintf("%s epub standard ebooks", query)

	links, err := Search(fmtQuery)
	if err != nil {
		panic(err)
	}

	for _, link := range links {
		fmt.Println(link)
	}
}
