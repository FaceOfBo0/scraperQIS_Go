package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Logging() Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() { fmt.Println(r.URL.Path, time.Since(start)) }()

			hf(w, r)
		}
	}
}

func Chain(hf http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		hf = m(hf)
	}
	return hf
}

func RunServer() {

	scr := Scraper{
		Url: "https://qis.server.uni-frankfurt.de/qisserver/rds?state=verpublish&publishContainer=lectureInstList&publishid=80100",
		SemesterUrl: "https://qis.server.uni-frankfurt.de/qisserver/rds?state=change&type=6&moduleParameter=semesterSelect&nextdir=" +
			"change&next=SearchSelect.vm&subdir=applications&targettype=7&targetstate=change&getglobal=semester"}
	// load template files
	tmplRoot := template.Must(template.ParseFiles("templates/root.html"))
	tmplChart := template.Must(template.ParseFiles("templates/chart.html"))

	//define route handler functions
	rootRoute := func(rw http.ResponseWriter, r *http.Request) {
		tmplRoot.Execute(rw, struct {
			Semesters []Semester
		}{
			Semesters: scr.getSemesters(),
		})
	}

	chartRoute := func(rw http.ResponseWriter, r *http.Request) {
		/* 		if r.Method != http.MethodPost {
			tmplRoot.Execute(rw, nil)
			return
		} */
		scr.Semester = r.FormValue("semester")
		ordered_val := r.FormValue("ordered")
		orderDir := r.FormValue("orderDir")
		lastOrdered := r.FormValue("lastOrdered")
		searchQuery := r.FormValue("searchQuery")
		if ordered_val == "" {
			scr.createUrlOffset()
			scr.loadLectures()
			ordered_val = "0"
			orderDir = "asc"
		} else {
			if lastOrdered != ordered_val {
				orderDir = "asc"
			} else {
				if orderDir == "asc" {
					orderDir = "desc"
				} else {
					orderDir = "asc"
				}
			}
		}

		// Filter lectures based on search query
		filteredLectures := scr.lectures
		if searchQuery != "" {
			filteredLectures = []Lecture{}
			for _, lecture := range scr.lectures {
				if strings.Contains(strings.ToLower(lecture.Title), strings.ToLower(searchQuery)) ||
					strings.Contains(strings.ToLower(lecture.Lecturers), strings.ToLower(searchQuery)) ||
					strings.Contains(strings.ToLower(lecture.Room), strings.ToLower(searchQuery)) ||
					strings.Contains(strings.ToLower(lecture.Time), strings.ToLower(searchQuery)) ||
					slices.Contains(lecture.Modules, searchQuery) {
					filteredLectures = append(filteredLectures, lecture)
				}
			}
		}

		sortFunc := compareLecsFuncs(ordered_val)
		if orderDir == "desc" {
			sortFunc = func(lec_a, lec_b Lecture) int {
				return -compareLecsFuncs(ordered_val)(lec_a, lec_b)
			}
		}

		slices.SortFunc(filteredLectures, sortFunc)
		tmplChart.Execute(rw, struct {
			Lectures *[]Lecture
			Semester string
			OrderDir string
			Ordered  string
		}{
			Lectures: &filteredLectures,
			Ordered:  ordered_val,
			Semester: scr.lectures[0].Semester,
			OrderDir: orderDir,
		})
	}

	downloadRoute := func(rw http.ResponseWriter, r *http.Request) {
		lecsJson, err := json.MarshalIndent(&scr.lectures, "", "    ")
		if err != nil {
			fmt.Printf("error with json serialization: %v\n", err)
			return
		}

		err = saveStrToFile(string(lecsJson), "lectures.json")
		if err != nil {
			fmt.Printf("error saving JSON to file: %v\n", err)
			return
		}
		cmd := exec.Command("java", "-jar", "lfw.jar", "lectures.json")
		if err := cmd.Run(); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
	}

	catalogRoute := func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(rw, r, "/", http.StatusSeeOther)
			return
		}

		// Maximum upload of 10 MB files
		r.ParseMultipartForm(10 << 20)

		// Get form values
		titleColumn, _ := strconv.Atoi(r.FormValue("title"))
		olatColumn, _ := strconv.Atoi(r.FormValue("olat"))

		// Get the file from form data
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Error retrieving file: %v\n", err)
			http.Error(rw, "Error retrieving file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		catInfos := getCatInfos(file, titleColumn, olatColumn)
		catStrings := mapList(catInfos, func(ci CatalogInfo) string {
			return fmt.Sprintf("%v\n%v\n%v", ci.lecTitle, ci.lecLink, ci.olatLink)
		})
		catStr := strings.Join(catStrings, "\n\n")

		fname := strings.TrimSuffix(header.Filename, ".xlsx") + ".txt"
		if err := saveStrToFile(catStr, fname); err != nil {
			panic(err)
		}
	}

	// Bind handler functions and start server
	fmt.Println("Server is running on localhost:4567..")
	http.HandleFunc("/", rootRoute)
	http.HandleFunc("/chart", Chain(chartRoute, Logging()))
	http.HandleFunc("/download", Chain(downloadRoute, Logging()))
	http.HandleFunc("/catalog", Chain(catalogRoute, Logging()))
	http.ListenAndServe(":4567", nil)

}
