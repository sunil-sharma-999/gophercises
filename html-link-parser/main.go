package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func main() {
	flag.Usage = func() {
		fmt.Println("Provide HTML file paths as arguments.")
		flag.PrintDefaults()

	}
	flag.Parse()

	args := flag.Args()

	for _, inputStr := range args {

		buf, err := os.ReadFile(inputStr)
		if err != nil {
			log.Fatal(err)
		}
		str := string(buf)
		localLinks := getParsedLinks(str)
		for _, link := range localLinks {
			fmt.Printf("%v\n%v\n\n", link.Href, link.Text)
		}
	}

}

func getParsedLinks(htmlStr string) []Link {
	r, err := html.Parse(strings.NewReader(htmlStr))

	if err != nil {
		log.Fatal(err)
	}

	links := []Link{}

	traverseNode(r, &links, nil)

	for i, l := range links {
		links[i].Text = strings.TrimSpace(l.Text)
	}

	return links
}

func traverseNode(n *html.Node, links *[]Link, link *Link) {

	switch n.Type {
	case html.ElementNode:
		if len(n.Attr) != 0 {
		loop:
			for _, att := range n.Attr {
				if att.Key == "href" {
					*links = append(*links, Link{Href: att.Val, Text: ""})
					link = &(*links)[len(*links)-1]
					break loop
				}
			}
		}
	case html.TextNode:
		text := strings.TrimSpace(n.Data)
		if text != "" {
			if link != nil {
				(*link).Text += n.Data
			}

		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverseNode(c, links, link)
	}
}
