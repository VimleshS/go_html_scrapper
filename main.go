package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/VimleshS/html_parser/scraper"
)

func main() {
	url := flag.String("url", "", "a blog url whose content to be parsed \n \thttp://www.synerzip.com")
	directory := flag.String("d", "", "directory containing files \n \t/home/synerzip/Gate/raw_Data/Data")
	outputdir := flag.String("o", "", "directory containing files (Required)\n \t/home/synerzip/Gate/raw_Data/p/")
	flag.Parse()

	if *outputdir == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *url == "" && *directory == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *url != "" {
		res, err := http.Get(*url)
		if err != nil {
			fmt.Fprintf(os.Stdout, err.Error())
		}
		scrape := scraper.Scrape{Input: res.Body, Output: os.Stdout}
		scrape.Scrape()

	} else {
		files := []string{}
		filepath.Walk(*directory, func(path string, info os.FileInfo, err error) error {
			if len(path) > 0 {
				files = append(files, path)
			}
			return nil
		})

		var wg sync.WaitGroup
		wg.Add(len(files))
		for _, f := range files {

			go func(fpath string) {
				defer wg.Done()

				f, err := os.Open(fpath)
				if err != nil {
					fmt.Fprintf(os.Stderr, err.Error())
				}

				namepart := strings.Split(fpath, "/")
				name := namepart[len(namepart)-1]
				name = strings.Replace(name, " - Synerzip", "", -1)
				name = strings.TrimSuffix(name, filepath.Ext(name))
				o, err := os.Create(*outputdir + name + ".txt")
				if err != nil {
					fmt.Fprintf(os.Stderr, err.Error())
				}
				defer o.Close()

				scrape := scraper.Scrape{Input: f, Output: o}
				scrape.Scrape()

			}(f)

		}
		wg.Wait()
	}
}
