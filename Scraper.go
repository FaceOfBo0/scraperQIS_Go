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
	lectures      []*Lecture
}

func newScraper(offset string) *Scraper {
	return &Scraper{offset: offset}
}

func (s *Scraper) getLectures() []*Lecture {

	if s.lectures == nil {
		s.getLecturesLinks("https://qis.server.uni-frankfurt.de/qisserver/rds?state=verpublish&publishContainer=lectureInstList&publishid=80100")
		for _, elem := range s.lecturesLinks {
			s.lectures = append(s.lectures, newLecture(s.getLectureText(elem)))
		}
	}
	return s.lectures
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
