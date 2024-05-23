package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Name string
}

type SiteScraper interface {
	// Return the keywords that will be appended to the end
	// of the search query
	Keywords() string

	// Return the regex pattern used to filter search result links
	// Only links matching the filter will be considered
	Filter() string

	// Scrape the html and return links to files
	ScrapeFileLinks(*html.Node) ([]Link, error)
}

type VK struct{}

func (vk VK) Filter() string {
	return `vk\.com`
}

func (vk VK) Keywords() string {
	return "vk"
}

func (vk VK) ScrapeFileLinks(document *html.Node) ([]Link, error) {
	// Anchor elements with this class are links to file downloads
	// Of course, the class name's subject to change
	name := "SecondaryAttachment js-SecondaryAttachment SecondaryAttachment--interactive"
	nodes := FilterHTML(document, "class", name)

	var downloads []Link
	for _, n := range nodes {
		href, err := GetAttribute(n, "href")
		if err != nil {
			return nil, err
		}

		// The text associated to the link
		names := FilterHTML(n, "class", "SecondaryAttachment__childrenText")
		name := names[0].FirstChild.Data

		if !strings.Contains(href, "doc") {
			continue // Not a link to a file
		}

		fullHref := fmt.Sprintf("https://vk.com%s", href)
		link := Link{Name: name, Href: fullHref}
		downloads = append(downloads, link)
	}

	return downloads, nil
}

// Extract the real website url from the duckduckgo search result url
func extractWebsiteLink(rawLink string) string {
	link, _ := url.QueryUnescape(rawLink) // URL decode the query

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

func GetFileLinks(query string, scraper SiteScraper) ([]Link, error) {
	formattedQuery := fmt.Sprintf("%s epub %s", query, scraper.Keywords())

	links, err := searchDuckDuckGo(formattedQuery)
	if err != nil {
		return nil, err
	}

	links, err = filterLinks(scraper.Filter(), links)
	if err != nil {
		return nil, err
	}

	// Only consider the first 5 links
	limit := 5
	if len(links) < limit {
		limit = len(links)
	}
	links = links[:limit]

	var fileLinks []Link
	for _, link := range links {
		document, err := GetPage(link)
		if err != nil {
			return nil, err
		}

		files, err := scraper.ScrapeFileLinks(document)
		if err != nil {
			return nil, err
		}

		fileLinks = append(fileLinks, files...)
	}

	return fileLinks, nil
}
