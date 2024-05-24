package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"

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

type Searcher struct {
	group  sync.WaitGroup
	mutex  sync.Mutex
	errors chan error

	pageMax   uint8
	scrapers  []SiteScraper
	FileLinks []Link
}

func NewSearcher() Searcher {
	return Searcher{
		pageMax:  3,
		scrapers: []SiteScraper{VK{}},
		errors:   make(chan error),
	}
}

func (s *Searcher) processPage(link string, scraper SiteScraper) {
	defer s.group.Done() // One less goroutine running

	document, err := GetPage(link)
	if err != nil {
		s.errors <- err
		return
	}

	files, err := scraper.ScrapeFileLinks(document)
	if err != nil {
		s.errors <- err
		return
	}

	s.mutex.Lock()
	s.FileLinks = append(s.FileLinks, files...)
	s.mutex.Unlock()
}

func (s *Searcher) getFileLinks(query string, scraper SiteScraper) error {
	formattedQuery := fmt.Sprintf("%s epub %s", query, scraper.Keywords())

	results, err := searchDuckDuckGo(formattedQuery)
	if err != nil {
		return err
	}

	pageLinks, err := filterLinks(scraper.Filter(), results)
	if err != nil {
		return err
	}

	// Clamp the number of page limits
	limit := int(s.pageMax)
	if len(pageLinks) < limit {
		limit = len(pageLinks)
	}
	pageLinks = pageLinks[:limit]

	// Process each page concurrently
	for _, link := range pageLinks {
		s.group.Add(1)
		go s.processPage(link, scraper)
	}

	// Wait for all the goroutines to finish
	// and return a potential error
	s.group.Wait()
	if len(s.errors) > 0 {
		return <-s.errors
	}

	return nil
}

func (s *Searcher) Search(query string) error {
	for _, scraper := range s.scrapers {
		if err := s.getFileLinks(query, scraper); err != nil {
			return err
		}
	}
	return nil
}
