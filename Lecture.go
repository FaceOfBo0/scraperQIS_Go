package main

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

const ()

type Lecture struct {
	day, time, title, flags, room, commentary, lecturers, textRaw string
	dayPattern, timePattern, titlePattern, roomPattern,
	commentaryPattern, lecturersPattern, modulesPattern *regexp.Regexp
	modules, lecturersList []string
}

func newLecture(text string) Lecture {
	//  initialize Lecture with RegExs for all the parameters
	lec := Lecture{textRaw: text}
	lec.textRaw = strings.Join(strings.Fields(lec.textRaw), " ")
	lec.dayPattern = regexp.MustCompile(`([A-Z][a-z])\..*\d+:\d+`)
	lec.timePattern = regexp.MustCompile(`\d+:\d+`)
	lec.titlePattern = regexp.MustCompile(`<h1>(.*) - Einzelansicht`)
	lec.roomPattern = regexp.MustCompile(`<a class="regular" title="Details ansehen zu Raum ([A-Z][A-Z] \d*\.?\d+).*?"`)
	lec.modulesPattern = regexp.MustCompile("BM 1|BM 2|BM 3|AM 1|AM 2|AM 3|AM 4|AM 5|VM 1|VM 2|VM 3|GM 1|GM 2|GM 3")
	lec.lecturersPattern = regexp.MustCompile(`Zust.ndigkeit.*?<a.*?> (.*?) <.a>`)
	lec.commentaryPattern = regexp.MustCompile(`Inhalt\sKommentar(.*?)\s(Leistungsnachweis|Einsortiert in)`)

	// Match parameters with scraped text

	titleSubMatch := lec.titlePattern.FindStringSubmatch(lec.textRaw)
	if len(titleSubMatch) >= 2 {
		lec.title = titleSubMatch[1]
		lec.title = strings.ReplaceAll(lec.title, "&", "&amp;")
	} else {
		lec.title = "n.a."
	}

	timesList := lec.timePattern.FindAllString(lec.textRaw, -1)
	if len(timesList) >= 2 {
		lec.time = timesList[0][:2] + "-" + timesList[1][:2]
	} else {
		lec.time = "n.a."
	}

	daySubMatch := lec.dayPattern.FindStringSubmatch(lec.textRaw)
	if len(daySubMatch) >= 2 {
		lec.day = daySubMatch[1]
	} else {
		lec.day = "n.a."
	}

	roomSubMatch := lec.roomPattern.FindStringSubmatch(lec.textRaw)
	if len(roomSubMatch) >= 2 {
		lec.room = roomSubMatch[1]
	} else {
		lec.room = "n.a."
	}

	lec.commentary = lec.commentaryPattern.FindString(lec.textRaw)
	if len(lec.commentary) == 0 {
		lec.commentary = "n.a."
	}

	lecturersSubMatch := lec.lecturersPattern.FindStringSubmatch(lec.textRaw)
	if len(lecturersSubMatch) >= 2 {
		lecturersStr := lecturersSubMatch[1]
		lecturersStr = strings.ReplaceAll(lecturersStr, "&nbsp;", " ")
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
	} else {
		lec.lecturers = "n.a."
	}
	lec.modules = lec.modulesPattern.FindAllString(lec.textRaw, -1)
	lec.modules = mapList(lec.modules, func(elem string) string { return strings.ReplaceAll(elem, " ", "") })
	slices.SortFunc(lec.modules, compareModules)
	lec.modules = slices.Compact(lec.modules)

	lec.flags = "____"
	if lec.title != "n.a." {
		lec.flags = replaceIdx(lec.flags, "V", 0)
	}
	if lec.room != "n.a." {
		lec.flags = replaceIdx(lec.flags, "R", 1)
	}
	if len(lec.modules) != 0 {
		lec.flags = replaceIdx(lec.flags, "M", 2)
	}
	if lec.commentary != "n.a." && lec.commentary != "..." && lec.commentary != "" {
		lec.flags = replaceIdx(lec.flags, "B", 3)
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
