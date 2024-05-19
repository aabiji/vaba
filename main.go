package main

import (
	"fmt"
	"regexp"
)

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

/*
I feel like we should make filtering any list (and html node) generic.
How?

Links to files on vk have this class:
class="SecondaryAttachment js-SecondaryAttachment SecondaryAttachment--interactive"
class="SecondaryAttachment js-SecondaryAttachment SecondaryAttachment--interactive"
*/

func main() {
	query := "shakespear"
	fmtQuery := fmt.Sprintf("%s epub vk", query)

	links, err := Search(fmtQuery)
	if err != nil {
		panic(err)
	}

	links, err = filterLinks(`vk\.com`, links)
	if err != nil {
		panic(err)
	}

	for _, link := range links {
		fmt.Println(link)
	}
}
