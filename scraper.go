package main

import (
	"io"
	"log"
	"net/http"
	"slices"
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

type Semester struct {
	Name, Value string
}

type Scraper struct {
	Semester, Url, SemesterUrl string
	//lecturesLinks []string
	lectures []Lecture
}

func (s *Scraper) getSemesters() []Semester {
	doc := getDocFromUrl(s.SemesterUrl)
	selection := doc.Find("a.regular[name^='W'], a.regular[name^='S']")
	semesters := make([]Semester, 0, selection.Length())
	selection.Each(func(i int, sel *goquery.Selection) {
		semName := strings.TrimSpace(sel.Text())
		var semVal string
		if strings.HasPrefix(semName, "Winter") {
			semVal = semName[len(semName)-7:len(semName)-3] + ".2"
		} else {
			semVal = semName[len(semName)-4:] + ".1"
		}
		semesters = append(semesters, Semester{Name: semName, Value: semVal})
	})
	slices.Reverse(semesters)
	return semesters
}

func getDocFromUrl(url string) *goquery.Document {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc
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

	doc := getDocFromUrl(s.Url)

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
