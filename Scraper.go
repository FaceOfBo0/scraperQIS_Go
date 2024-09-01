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

/* func (s *Scraper) getLectures() []Lecture {
	if s.lectures == nil {
		s.loadLecturesLinks()

		// Create a channel to handle lecture links
		//lectureChan := make(chan *Lecture, len(s.lecturesLinks))
		var wg sync.WaitGroup

		// Worker function to process lecture links
		worker := func(link string) {
			defer wg.Done()
			s.lectures = append(s.lectures, newLecture(getHtmlText(link), link))
		}

		// Launch goroutines to process lecture links
		for _, link := range s.lecturesLinks {
			wg.Add(1)
			go worker(link)
		}

		// Close the channel once all goroutines are done
		wg.Wait()
		//close(lectureChan)

		// Collect results from the channel
		// for lecture := range lectureChan {
		//	s.lectures = append(s.lectures, lecture)
		//}
	}
	return s.lectures
} */

/* func (s *Scraper) getLectures() []Lecture {

	if s.lectures == nil {
		//s.loadLecturesLinks()
		s.loadLecturesLinksGQ()
		for _, elem := range s.lecturesLinks {
			s.lectures = append(s.lectures, newLecture(s.getLectureText(elem)))
		}
	}
	return s.lectures
} */

/* func (s *Scraper) loadLecturesLinksC() {
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

	err := col.Visit(s.url)
	if err != nil {
		log.Fatal(err)
	}
} */
