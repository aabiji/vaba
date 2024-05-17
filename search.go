package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Search DuckDuckGo and return the result links
func Search(query string) ([]string, error) {
	htmlContent, err := getSearchResultPage(query)
	if err != nil {
		return nil, err
	}

	document, err := html.Parse(htmlContent)
	if err != nil {
		return nil, err
	}

	links := getResultLinks(document)
	return links, nil
}

// Return an io.Reader containing the html
// of a DuckDuckGo search result page given a search query
func getSearchResultPage(query string) (io.Reader, error) {
	query = url.QueryEscape(query)
	baseLink := "https://html.duckduckgo.com/html"
	link := fmt.Sprintf("%s/?q=%s", baseLink, query)

	request, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}
	// Trick DuckDuckGo into thinking we're a real
	// user with a fake user Agent
	fakeAgent := "'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36'"
	request.Header.Add("User-Agent", fakeAgent)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		msg := "status code : %d | failed to get page"
		return nil, fmt.Errorf(msg, response.StatusCode)
	}

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(content), nil
}

// Extract the real website url from the duckduckgo search result url
func extractWebsiteLink(rawLink string) string {
	link, _ := url.QueryUnescape(rawLink)

	// Remove the prefix
	prefix := "//duckduckgo.com/l/?uddg="
	link = strings.Replace(link, prefix, "", -1)

	// Remove everything after the regex pattern match
	pattern := `\&rut=.*`
	regex, _ := regexp.Compile(pattern)
	linkParts := regex.Split(link, -1)

	return linkParts[0]
}

// Scrape all the links to websites off the DuckDuckGo
// search result page
func getResultLinks(node *html.Node) []string {
	var allLinks []string

	if node.Type == html.ElementNode && node.Data == "a" {
		isResultUrl := false
		for _, attr := range node.Attr {
			if attr.Key == "class" && attr.Val == "result__url" {
				isResultUrl = true
			}

			if attr.Key == "href" && isResultUrl {
				link := extractWebsiteLink(attr.Val)
				allLinks = append(allLinks, link)
			}
		}
	}

	for next := node.FirstChild; next != nil; next = next.NextSibling {
		links := getResultLinks(next)
		allLinks = append(allLinks, links...)
	}

	return allLinks
}
