package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

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
	url           string
	lectures      []Lecture
}

func newScraper(url string, offset string) Scraper {
	return Scraper{url: url, offset: offset}
}

func (s *Scraper) getLectures() []Lecture {

	if s.lectures == nil {
		s.loadLecturesLinks()
		for _, elem := range s.lecturesLinks {
			s.lectures = append(s.lectures, newLecture(s.getLectureText(elem)))
		}
	}
	return s.lectures
}

func (s *Scraper) getLecturesConc() []Lecture {
	if s.lectures == nil {
		s.loadLecturesLinks()

		// Create a channel to handle lecture links
		lectureChan := make(chan Lecture, len(s.lecturesLinks))
		var wg sync.WaitGroup

		// Worker function to process lecture links
		worker := func(link string) {
			defer wg.Done()
			lecture := newLecture(s.getLectureText(link))
			lectureChan <- lecture
		}

		// Launch goroutines to process lecture links
		for _, link := range s.lecturesLinks {
			wg.Add(1)
			go worker(link)
		}

		// Close the channel once all goroutines are done
		go func() {
			wg.Wait()
			close(lectureChan)
		}()

		// Collect results from the channel
		for lecture := range lectureChan {
			s.lectures = append(s.lectures, lecture)
		}
	}
	return s.lectures
}

func (s *Scraper) getLectureText(url string) string {
	return getHtmlText(url + s.offset)
}

func (s *Scraper) loadLecturesLinks() {
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

	col.Visit(s.url)
}
