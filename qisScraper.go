package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gocolly/colly/v2"
)

func getHtmlText(url string) string {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	content, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

type Scraper struct {
	lecturesLinks []string
	offset        string
}

func newScraper(offset string) *Scraper {
	return &Scraper{offset: offset}
}

func (s *Scraper) getLectureText(url string) string {
	return getHtmlText(url + s.offset)
}

func (s *Scraper) getLecturesLinks(url string) {
	lecLinks := make([]string, 0)

	col := colly.NewCollector()

	col.OnHTML("td", func(elems *colly.HTMLElement) {
		elems.ForEach("a", func(i int, linkElem *colly.HTMLElement) {
			lecLinks = append(lecLinks, linkElem.Attr("href")+s.offset)
		})
	})

	col.OnError(func(res *colly.Response, e error) {
		fmt.Println("Error: ", e)
	})

	col.OnScraped(func(res *colly.Response) {
		s.lecturesLinks = lecLinks
	})

	col.Visit(url)
}
