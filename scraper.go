package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
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
	Semester, Url string
	//lecturesLinks []string
	lectures []Lecture
}

func (s *Scraper) createUrlOffset() {
	if s.Semester != "" && len(s.Semester) == 6 {
		year := strings.Split(s.Semester, ".")[0]
		half := strings.Split(s.Semester, ".")[1]
		offset := "&k_semester.semid=" + year + half +
			"&idcol=k_semester.semid&idval=" + year + half +
			"&purge=n&getglobal=semester"
		s.Url += offset
	}
}

func (s *Scraper) loadLectures() {

	resp, err := http.Get(s.Url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Create a channel to collect links
	lecturesChan := make(chan Lecture)
	var wg sync.WaitGroup

	// Function to process each selection and send links to the channel
	processSelection := func(sel *goquery.Selection) {
		defer wg.Done()
		link := sel.Find("a").AttrOr("href", "")
		if link != "" {
			lecturesChan <- newLecture(getHtmlText(link), link)
		}
	}

	// Find all "td" elements and launch a goroutine for each
	doc.Find("td").Each(func(_ int, sel *goquery.Selection) {
		wg.Add(1)
		go processSelection(sel)
	})

	// Close the channel once all goroutines are done
	go func() {
		wg.Wait()
		close(lecturesChan)
	}()

	// Collect results from the channel
	s.lectures = make([]Lecture, len(lecturesChan))
	for lec := range lecturesChan {
		s.lectures = append(s.lectures, lec)
	}
}
