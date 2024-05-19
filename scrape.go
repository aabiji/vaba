package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"slices"

	"golang.org/x/net/html"
)

// Return an io.Reader containing the html for the page
func getPageHTML(link string) (io.Reader, error) {
	request, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}
	// Trick the website into thinking we're a real user with a fake user Agent
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

func GetPage(link string) (*html.Node, error) {
	page, err := getPageHTML(link)
	if err != nil {
		return nil, err
	}

	root, err := html.Parse(page)
	return root, err
}

// Return the nodes in the html tree with a specific attribute
func FilterHTML(node *html.Node, attrName, attrValue string) []*html.Node {
	var allMatches []*html.Node

	attr := html.Attribute{
		Key: attrName,
		Val: attrValue,
	}

	if slices.Index(node.Attr, attr) != -1 {
		allMatches = append(allMatches, node)
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		matches := FilterHTML(child, attrName, attrValue)
		allMatches = append(allMatches, matches...)
	}

	return allMatches
}

// Get a attribute from a html node
func GetAttribute(node *html.Node, name string) (string, error) {
	for _, attr := range node.Attr {
		if attr.Key == name {
			return attr.Val, nil
		}
	}
	return "", fmt.Errorf("%s attribute not found", name)
}
