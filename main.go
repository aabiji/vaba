package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Extract the real website url from the duckduckgo search result url
func extractWebsiteLink(rawLink string) string {
	link, _ := url.QueryUnescape(rawLink) // URL deocde the query

	// Remove the prefix
	prefix := "//duckduckgo.com/l/?uddg="
	link = strings.Replace(link, prefix, "", -1)

	// Remove everything after the regex pattern match
	pattern := `\&rut=.*`
	regex := regexp.MustCompile(pattern)
	linkParts := regex.Split(link, -1)

	return linkParts[0]
}

// Search DuckDuckGo and return the result links
func searchDuckDuckGo(query string) ([]string, error) {
	// Format the duckduckgo search page url
	query = url.QueryEscape(query)
	baseLink := "https://html.duckduckgo.com/html"
	link := fmt.Sprintf("%s/?q=%s", baseLink, query)

	document, err := GetPage(link)
	if err != nil {
		return nil, err
	}

	var links []string
	nodes := FilterHTML(document, "class", "result__url")

	for _, n := range nodes {
		href, err := GetAttribute(n, "href")
		if err != nil {
			return nil, err
		}
		links = append(links, extractWebsiteLink(href))
	}

	return links, nil
}

// Only return links that match the pattern
func filterLinks(pattern string, links []string) ([]string, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	filtered := []string{}
	for _, link := range links {
		if regex.Match([]byte(link)) {
			filtered = append(filtered, link)
		}
	}

	return filtered, nil
}

func vkGetFileLinks(document *html.Node) ([]string, error) {
	// Anchor elements with this class are links to file downloads
	// Of course, the class name's subject to change
	name := "SecondaryAttachment js-SecondaryAttachment SecondaryAttachment--interactive"
	nodes := FilterHTML(document, "class", name)

	var downloads []string
	for _, n := range nodes {
		link, err := GetAttribute(n, "href")
		if err != nil {
			return nil, err
		}

		fullLink := fmt.Sprintf("https://vk.com%s", link)
		downloads = append(downloads, fullLink)
	}

	return downloads, nil
}

func main() {
	query := "shakespear"
	fmtQuery := fmt.Sprintf("%s epub vk", query)

	links, err := searchDuckDuckGo(fmtQuery)
	if err != nil {
		panic(err)
	}

	links, err = filterLinks(`vk\.com`, links)
	if err != nil {
		panic(err)
	}

	var fileDownloads []string
	for _, link := range links {
		document, err := GetPage(link)
		if err != nil {
			panic(err)
		}

		downloads, err := vkGetFileLinks(document)
		if err != nil {
			panic(err)
		}
		fileDownloads = append(fileDownloads, downloads...)
	}

	for _, download := range fileDownloads {
		fmt.Println(download)
	}
}
