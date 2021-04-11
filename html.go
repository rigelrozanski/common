package common

import (
	"net/http"

	"golang.org/x/net/html"
)

func GetUrlTitle(url string) (title string, err error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	title, _ = propogateForTitle(doc)
	return title, nil
}

func propogateForTitle(n *html.Node) (title string, found bool) {
	if n.Type == html.ElementNode && n.Data == "title" {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title, found := propogateForTitle(c)
		if found {
			return title, found
		}
	}
	return "", false
}
