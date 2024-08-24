package main

import (
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	lessDayList = []Tuple[string, string]{{"Mo", "Di"}, {"Mo", "Mi"}, {"Mo", "Do"}, {"Mo", "Fr"},
		{"Di", "Mi"}, {"Di", "Do"}, {"Di", "Fr"}, {"Mi", "Do"}, {"Mi", "Fr"}, {"Do", "Fr"}, {"Mo", "Sa"},
		{"Mo", "So"}, {"Di", "Sa"}, {"Di", "So"}, {"Mi", "Sa"}, {"Mi", "So"}, {"Do", "Sa"}, {"Do", "So"},
		{"Fr", "Sa"}, {"Fr", "So"}, {"Sa", "So"}}
)

type Lecture struct {
	Day, Time, Title, Flags, Room, Commentary, Lecturers, TextRaw, Link string
	Modules, LecturersList                                              []string
	dayPattern, timePattern, titlePattern, roomPattern,
	commentaryPattern, lecturersPattern, modulesPattern *regexp.Regexp
}

func newLecture(text string, url string) Lecture {
	//  initialize Lecture with RegExs for all the parameters
	lec := Lecture{TextRaw: text, Link: url}
	lec.TextRaw = strings.Join(strings.Fields(lec.TextRaw), " ")
	lec.dayPattern = regexp.MustCompile(`([A-Z][a-z])\..*\d+:\d+`)
	lec.timePattern = regexp.MustCompile(`\d+:\d+`)
	lec.titlePattern = regexp.MustCompile(`<h1>(.*) - Einzelansicht`)
	lec.roomPattern = regexp.MustCompile(`<a class="regular" title="Details ansehen zu Raum ([A-Z][A-Z] \d*\.?\d+).*?"`)
	lec.modulesPattern = regexp.MustCompile("BM 1|BM 2|BM 3|AM 1|AM 2|AM 3|AM 4|AM 5|VM 1|VM 2|VM 3|GM 1|GM 2|GM 3")
	lec.lecturersPattern = regexp.MustCompile(`Zust.ndigkeit.*?<a.*?> (.*?) <.a>`)
	lec.commentaryPattern = regexp.MustCompile(`Inhalt\sKommentar(.*?)\s(Leistungsnachweis|Einsortiert in)`)

	// Match parameters with scraped text

	titleSubMatch := lec.titlePattern.FindStringSubmatch(lec.TextRaw)
	if len(titleSubMatch) >= 2 {
		lec.Title = titleSubMatch[1]
		lec.Title = strings.ReplaceAll(lec.Title, "&", "&amp;")
	} else {
		lec.Title = "n.a."
	}

	timesList := lec.timePattern.FindAllString(lec.TextRaw, -1)
	if len(timesList) >= 2 {
		lec.Time = timesList[0][:2] + "-" + timesList[1][:2]
	} else {
		lec.Time = "n.a."
	}

	daySubMatch := lec.dayPattern.FindStringSubmatch(lec.TextRaw)
	if len(daySubMatch) >= 2 {
		lec.Day = daySubMatch[1]
	} else {
		lec.Day = "n.a."
	}

	roomSubMatch := lec.roomPattern.FindStringSubmatch(lec.TextRaw)
	if len(roomSubMatch) >= 2 {
		lec.Room = roomSubMatch[1]
	} else {
		lec.Room = "n.a."
	}

	lec.Commentary = lec.commentaryPattern.FindString(lec.TextRaw)
	if len(lec.Commentary) == 0 {
		lec.Commentary = "n.a."
	}

	lecturersSubMatch := lec.lecturersPattern.FindStringSubmatch(lec.TextRaw)
	if len(lecturersSubMatch) >= 2 {
		lecturersStr := lecturersSubMatch[1]
		lecturersStr = strings.ReplaceAll(lecturersStr, "&nbsp;", " ")
		if len(lecturersStr) == 0 {
			lec.LecturersList = append(lec.LecturersList, "n.a.")
		} else {
			lecturersArr := strings.Split(lecturersStr, ", ")
			if len(lecturersArr) != 0 {
				lec.LecturersList = append(lec.LecturersList, lecturersArr[0])
				if len(lecturersArr) > 3 {
					secLecturerList := strings.Split(lecturersArr[2], " ")
					lec.LecturersList = append(lec.LecturersList, secLecturerList[len(secLecturerList)-1])
				}
			}
		}

		if len(lec.LecturersList) > 1 {
			lec.Lecturers = strings.Join(lec.LecturersList, ", ")
		} else if len(lec.LecturersList) == 1 {
			lec.Lecturers = lec.LecturersList[0]
		}
	} else {
		lec.Lecturers = "n.a."
	}

	lec.Modules = lec.modulesPattern.FindAllString(lec.TextRaw, -1)
	lec.Modules = mapList(lec.Modules, func(elem string) string { return strings.ReplaceAll(elem, " ", "") })
	slices.SortFunc(lec.Modules, compareModules)
	lec.Modules = slices.Compact(lec.Modules)

	lec.Flags = "____"
	if lec.Title != "n.a." {
		lec.Flags = replaceIdx(lec.Flags, "V", 0)
	}
	if lec.Room != "n.a." {
		lec.Flags = replaceIdx(lec.Flags, "R", 1)
	}
	if len(lec.Modules) != 0 {
		lec.Flags = replaceIdx(lec.Flags, "M", 2)
	}
	if lec.Commentary != "n.a." && lec.Commentary != "..." && lec.Commentary != "" {
		lec.Flags = replaceIdx(lec.Flags, "B", 3)
	}

	return lec
}

func lessDay(day_a string, day_b string) bool {
	if day_a == "Block" || day_a == "n.a." {
		return false
	} else if day_b == "Block" || day_b == "n.a." {
		return true
	} else {
		dayTuple := Tuple[string, string]{fst: day_a, snd: day_b}
		return slices.Contains(lessDayList, dayTuple)
	}

}

func lessTime(time_a_str string, time_b_str string) bool {
	if time_a_str == "n.a." {
		return false
	} else if time_b_str == "n.a." {
		return true
	} else {
		time_a, _ := strconv.Atoi(time_a_str[0:2])
		time_b, _ := strconv.Atoi(time_b_str[0:2])
		return time_a < time_b
	}

}

func compareLecsByDays(lec_a Lecture, lec_b Lecture) int {
	if lessDay(lec_a.Day, lec_b.Day) || ((lec_a.Day == lec_b.Day) && lessTime(lec_a.Time, lec_b.Time)) {
		return -1
	} else {
		return 1
	}
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
