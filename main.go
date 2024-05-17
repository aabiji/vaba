package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

func getSearchResultPage(query string) (io.Reader, error) {
	link := "https://duckduckgo.com/?q=%s&ia=web"
	escapedQuery := url.QueryEscape(query)
	searchLink := fmt.Sprintf(link, escapedQuery)

	client := &http.Client{}

	request, err := http.NewRequest("GET", searchLink, nil)
	if err != nil {
		return nil, err
	}

	userAgent := "'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36'"
	request.Header.Add("User-Agent", userAgent)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New("failed to get page")
	}

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(content), nil
}

func findLinks(node *html.Node) {
	if node.Type == html.ElementNode && node.Data == "a" {
		fmt.Println(node.Data, node.Namespace, node.Attr)
	}

	for next := node.FirstChild; next != nil; next = next.NextSibling {
		findLinks(next)
	}
}

func main() {
	query := "please work"
	htmlContent, err := getSearchResultPage(query)
	if err != nil {
		panic(err)
	}

	doc, err := html.Parse(htmlContent)
	if err != nil {
		panic(err)
	}

	findLinks(doc)
}
