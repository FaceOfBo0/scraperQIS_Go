package main

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

const ()

type Lecture struct {
	day, time, title, flags, room, commentary, lecturers, link, textRaw string
	dayPattern, timePattern, tittlePattern, roomPattern,
	commentaryPattern, lecturersPattern, modulesPattern *regexp.Regexp
	modules, lecturersList []string
}

func newLecture(text string, url string) *Lecture {
	//  initialize Lecture with RegExs for all the parameters
	lec := &Lecture{textRaw: text, link: url}
	lec.dayPattern, _ = regexp.Compile(`([A-Z][a-z])\.\s\d+:\d+`)
	lec.timePattern, _ = regexp.Compile(`\d+:\d+`)
	lec.tittlePattern, _ = regexp.Compile(`:\sStartseite\s(.*)\s-\sEinzelansicht`)
	lec.roomPattern, _ = regexp.Compile(`woch.*?-\s([A-Z][A-Z]\s\d*\.?\d+).*\sGruppe`)
	lec.modulesPattern, _ = regexp.Compile("BM 1|BM 2|BM 3|AM 1|AM 2|AM 3|AM 4|AM 5|VM 1|VM 2|VM 3|GM 1|GM 2|GM 3")
	lec.lecturersPattern, _ = regexp.Compile(`Zuständigkeit\s(.+?)\s(Studiengänge\sAbschluss|Zuordnung\szu)+`)
	lec.commentaryPattern, _ = regexp.Compile(`Inhalt\sKommentar(.*?)\s(Leistungsnachweis|Einsortiert in)`)

	// Match parameters with scraped text

	lec.title = lec.tittlePattern.FindString(lec.textRaw)
	lec.title = strings.ReplaceAll(lec.title, "&", "&amp;")
	if len(lec.title) == 0 {
		lec.title = "n.a."
	}

	timesList := lec.timePattern.FindAllString(lec.textRaw, -1)
	if len(timesList) >= 2 {
		lec.time = timesList[0][:2] + "-" + timesList[1][:2]
	} else {
		lec.time = "n.a."
	}

	lec.day = lec.dayPattern.FindString(lec.textRaw)
	if len(lec.day) == 0 {
		lec.day = "n.a."
	}

	lec.room = lec.roomPattern.FindString(lec.textRaw)
	if len(lec.room) == 0 {
		lec.room = "n.a."
	}

	lec.commentary = lec.commentaryPattern.FindString(lec.textRaw)
	if len(lec.commentary) == 0 {
		lec.commentary = "n.a."
	}

	lecturersStr := lec.lecturersPattern.FindString(lec.textRaw)
	if len(lecturersStr) == 0 {
		lec.lecturersList = append(lec.lecturersList, "n.a.")
	} else {
		lecturersArr := strings.Split(lecturersStr, ", ")
		if len(lecturersArr) != 0 {
			lec.lecturersList = append(lec.lecturersList, lecturersArr[0])
			if len(lecturersArr) > 3 {
				secLecturerList := strings.Split(lecturersArr[2], " ")
				lec.lecturersList = append(lec.lecturersList, secLecturerList[len(secLecturerList)-1])
			}
		}
	}

	if len(lec.lecturersList) > 1 {
		lec.lecturers = strings.Join(lec.lecturersList, ", ")
	} else if len(lec.lecturersList) == 1 {
		lec.lecturers = lec.lecturersList[0]
	}

	lec.modules = lec.modulesPattern.FindAllString(lec.textRaw, -1)
	mapList(lec.modules, func(elem string) string { return strings.ReplaceAll(elem, " ", "") })
	slices.SortFunc(lec.modules, compareModules)

	lec.flags = "V___"
	if lec.room == "n.a" {
		lec.flags = replaceIdx(lec.flags, "R", 1)
	}

	return lec
}

func lessDay(day_a string, day_b string) bool {
	return false
}

func lessTime(time_a string, time_b string) bool {
	return false
}

func compareLecByDays(lec_a *Lecture, lec_b *Lecture) int {
	lessDayList := []Tuple[string, string]{{"Mo", "Di"}, {"Mo", "Mi"}, {"Mo", "Do"}, {"Mo", "Fr"},
		{"Di", "Mi"}, {"Di", "Do"}, {"Di", "Fr"}, {"Mi", "Do"}, {"Mi", "Fr"}, {"Do", "Fr"}, {"Mo", "Sa"},
		{"Mo", "So"}, {"Di", "Sa"}, {"Di", "So"}, {"Mi", "Sa"}, {"Mi", "So"}, {"Do", "Sa"}, {"Do", "So"},
		{"Fr", "Sa"}, {"Fr", "So"}, {"Sa", "So"}}
	fmt.Println(lessDayList)
	return 0
}

func compareModules(mod_a string, mod_b string) int {
	if (strings.HasPrefix(mod_a, "BM") && strings.HasPrefix(mod_b, "AM")) || (strings.HasPrefix(mod_a, "BM") && strings.HasPrefix(mod_b, "VM")) ||
		(strings.HasPrefix(mod_a, "BM") && strings.HasPrefix(mod_b, "GM")) || (strings.HasPrefix(mod_a, "AM") && strings.HasPrefix(mod_b, "VM")) ||
		(strings.HasPrefix(mod_a, "AM") && strings.HasPrefix(mod_b, "GM")) || (strings.HasPrefix(mod_a, "VM") && strings.HasPrefix(mod_b, "GM")) {
		return -1
	} else {
		if mod_a[:2] == mod_b[:2] {
			if int(mod_a[2]) < int(mod_b[2]) {
				return -1
			}
		}
		return 1
	}
}
