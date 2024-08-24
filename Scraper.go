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
	semester, url string
	lecturesLinks []string
	lectures      []Lecture
}

func (s *Scraper) createUrlOffset() {
	if s.semester != "" && len(s.semester) == 6 {
		year := strings.Split(s.semester, ".")[0]
		half := strings.Split(s.semester, ".")[1]
		offset := "&k_semester.semid=" + year + half +
			"&idcol=k_semester.semid&idval=" + year + half +
			"&purge=n&getglobal=semester"
		s.url += offset
	}
}

func (s *Scraper) getLectures() []Lecture {
	if s.lectures == nil {
		//s.loadLecturesLinks()
		s.loadLecturesLinks()

		// Create a channel to handle lecture links
		lectureChan := make(chan Lecture, len(s.lecturesLinks))
		var wg sync.WaitGroup

		// Worker function to process lecture links
		worker := func(link string) {
			defer wg.Done()
			lecture := newLecture(getHtmlText(link), link)
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

func (s *Scraper) loadLecturesLinks() {
	if len(s.lecturesLinks) == 0 {
		resp, err := http.Get(s.url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		doc.Find("td").Each(func(_ int, sel *goquery.Selection) {
			link := sel.Find("a").AttrOr("href", "")
			if link != "" {
				s.lecturesLinks = append(s.lecturesLinks, link)
			}
		})
	}

}

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

/* func (s *Scraper) loadLecturesLinks() {
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
