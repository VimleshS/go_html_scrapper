package scraper

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Scrape requires source and dest
type Scrape struct {
	Input  io.Reader
	Output io.WriteCloser
}

// Scrape extract text out of file
func (s *Scrape) Scrape() error {
	doc, err := html.Parse(s.Input)
	if err != nil {
		return err
	}
	s.f(doc)
	return nil
}

func (s *Scrape) f(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "wpb_column vc_column_container vc_col-sm-8" {
				s.f1(n)
				break
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		s.f(c)
	}
}
func (s *Scrape) f1(n *html.Node) {
	if n.Type == html.TextNode {
		data := strings.TrimSpace(n.Data)
		if data != "" {
			fmt.Fprintf(s.Output, data)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			s.f(c)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		s.f1(c)
	}
}
