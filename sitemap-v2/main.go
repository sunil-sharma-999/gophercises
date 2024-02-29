package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Link struct {
	Href     string `json:"href"`
	visited  bool
	Children []*Link `json:"children"`
}

func main() {

	outputFileName := flag.String("o", "", "Output file name to save JSON output. If not provided it will use [timestamp].json")
	mainURL := flag.String("u", "", "URL to crawl.")
	maxDepth := flag.Int("depth", 100, "Max depth.")
	flag.Parse()

	if strings.TrimSpace(*mainURL) == "" {
		log.Fatal("Provide valid URL to crawl.")
	}

	if strings.TrimSpace(*outputFileName) == "" {
		*outputFileName = fmt.Sprintf("%v.json", time.Now().Nanosecond())
	}
	if !strings.HasSuffix(*outputFileName, ".json") {
		log.Fatal("Invalid json file name")
	}

	baseURL, err := url.Parse(*mainURL)
	if err != nil {
		log.Fatal(err)
	}
	baseURLString := baseURL.String()

	result := Link{
		Href:     *mainURL,
		visited:  false,
		Children: []*Link{},
	}
	rootNode := &result
	processURL(&result, maxDepth, rootNode, baseURLString)

	file, err := os.Create(*outputFileName)
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	file.Write(bytes)
	defer file.Close()
}

func processURL(linkNode *Link, maxDepth *int, rootNode *Link, baseURLString string) {
	if *maxDepth < 0 {
		return
	}
	*maxDepth -= 1

	resp, err := http.Get(linkNode.Href)
	(*linkNode).visited = true
	if err != nil {
		return
	}
	defer resp.Body.Close()
	html, err := html.Parse(resp.Body)
	if err != nil {
		return
	}
	traverseNode(html, linkNode, rootNode, baseURLString)
	for _, l := range linkNode.Children {
		if !l.visited {
			processURL(l, maxDepth, rootNode, baseURLString)
		}
	}

}

func urlHasBeenVisited(url string, node *Link) bool {
	if url == node.Href {
		return true
	} else {
		for _, l := range node.Children {
			if urlHasBeenVisited(url, l) {
				return true
			}
		}
		return false
	}

}

func traverseNode(n *html.Node, linkNode *Link, rootNode *Link, baseURLString string) {

	switch n.Type {
	case html.ElementNode:
		if len(n.Attr) != 0 {
		loop:
			for _, att := range n.Attr {
				if att.Key == "href" {
					href := att.Val
					if strings.HasPrefix(href, "/") {
						href = baseURLString + href
					}

					if urlHasBeenVisited(href, rootNode) {
						break loop
					}

					newNode := Link{
						Href:     href,
						visited:  false,
						Children: []*Link{},
					}
					linkNode.Children = append(linkNode.Children, &newNode)
					break loop
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverseNode(c, linkNode, rootNode, baseURLString)
	}

}
