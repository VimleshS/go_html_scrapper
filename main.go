package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	url := flag.String("url", "http://www.synerzip.com", "a blog url whose content to be parsed")
	flag.Parse()

	// res, err := http.Get("http://www.synerzip.com/wcf-service-using-cassandra/")
	res, err := http.Get(*url)
	if err != nil {
		fmt.Fprintf(os.Stdout, err.Error())
	}
	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	fmt.Fprintf(os.Stdout, err.Error())
	// }
	// fmt.Fprintf(os.Stdout, string(body))

	doc, err := html.Parse(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stdout, err.Error())
	}

	var f func(*html.Node)
	var f1 func(*html.Node)

	f1 = func(n *html.Node) {
		if n.Type == html.TextNode {
			if !(strings.Trim(n.Data, "\n") == "" || strings.Trim(n.Data, " ") == "") {
				fmt.Println(n.Data)
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f1(c)
		}
	}

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "wpb_column vc_column_container vc_col-sm-8" {
					f1(n)
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	/*
		    var f func(*html.Node, html.NodeType, string, string)
			f = func(n *html.Node, nt html.NodeType, key string, value string) {

				// if n.Type == nt && n.Data == "div" {
				// 	for _, a := range n.Attr {
				// 		if a.Key == key && a.Val == value {
				// 			// fmt.Println(a.Val)
				// 			//f(n, html.TextNode, "")
				// 			break
				// 		}
				// 	}
				// }

				if n.Type == nt && n.Data == "div" {
					for _, a := range n.Attr {
						if a.Key == key && a.Val == value {
							if key == "id" {
								f(doc, html.ElementNode, "class", "wpb_column vc_column_container vc_col-sm-8")
							} else {
								f(n, html.TextNode, key, "")
							}
							break
						}
					}
				} else if n.Type == html.TextNode {
					if strings.Trim(n.Data, "\n") != "" {
						fmt.Println(n.Data)
					}
				}

				for c := n.FirstChild; c != nil; c = c.NextSibling {
					f(c, nt, key, value)
				}
			}
			// f(doc, html.ElementNode, "class", "wpb_column vc_column_container vc_col-sm-8")
			f(doc, html.ElementNode, "id", "main")
	*/

}
