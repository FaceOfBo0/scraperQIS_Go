package main

import "regexp"

type Lecture struct {
	day, time, title, flags, room, commentary, lectStr, link, textRaw                                        string
	dayPattern, timePattern, tittlePattern, roomPattern, commentaryPattern, lecturersPattern, modulesPattern *regexp.Regexp
	modules, lecturers                                                                                       []string
}

func newLecture(text string, url string) *Lecture {
	lec := Lecture{textRaw: text, link: url}
	lec.dayPattern, _ = regexp.Compile(`([A-Z][a-z])\.\s\d+:\d+`)
	lec.timePattern, _ = regexp.Compile(`\d+:\d+`)
	return &lec
}
