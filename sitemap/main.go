package main

import (
	"encoding/xml"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Links = map[string]bool

type Loc struct {
	Loc string `xml:"loc"`
}

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	URL     []Loc    `xml:"url"`
}

func main() {

	var inputURL string
	var destFile string
	var depth int
	flag.StringVar(&inputURL, "url", "", "Provide URL to crawl")
	flag.StringVar(&destFile, "o", "", "Provide path for saving generated xml.")
	flag.IntVar(&depth, "depth", 100, "Provide max depth.")
	flag.Parse()

	if strings.TrimSpace(destFile) == "" {
		log.Fatal("Invalid file path")
	}

	u, err := url.ParseRequestURI(inputURL)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}

	htmlNode, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	links := Links{}

	traverseNode(htmlNode, &links, u.String())
	generateSiteMap(&links, &depth, u.String())

	file, err := os.Create(destFile)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
	<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	for link := range links {
		file.WriteString("\n\t<url>\n\t\t<loc>")
		xml.EscapeText(file, []byte(link))
		file.WriteString("</loc>\n\t</url>")
	}
	file.WriteString("\n</urlset>")
}

func generateSiteMap(links *Links, depth *int, mainURL string) {
	for link, visited := range *links {
		if *depth < 0 {
			return
		}
		if visited {
			continue
		} else {
			(*links)[link] = true
			u, err := url.ParseRequestURI(link)
			if err != nil {
				delete(*links, link)
				continue
			}

			resp, err := http.Get(u.String())
			if err != nil {
				delete(*links, link)
				continue
			}

			htmlNode, err := html.Parse(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			err = resp.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			traverseNode(htmlNode, links, mainURL)
			*depth -= 1
			generateSiteMap(links, depth, mainURL)
		}
	}

}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func traverseNode(n *html.Node, links *Links, mainURL string) {

	switch n.Type {
	case html.ElementNode:
		if len(n.Attr) != 0 {
		loop:
			for _, att := range n.Attr {
				if att.Key == "href" || att.Key == "src" {

					_, ok := (*links)[att.Val]

					if !ok {
						linkStr := att.Val
						if strings.HasPrefix(linkStr, "/") {
							linkStr = mainURL + linkStr
						}
						(*links)[linkStr] = false
					}
					break loop
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverseNode(c, links, mainURL)
	}
}
