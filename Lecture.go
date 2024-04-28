package main

import (
	"regexp"
	"strings"
)

func NewLecture(text string, url string) *Lecture {
	lec := &Lecture{textRaw: text, link: url}
	lec.dayPattern, _ = regexp.Compile(`([A-Z][a-z])\.\s\d+:\d+`)
	lec.timePattern, _ = regexp.Compile(`\d+:\d+`)
	lec.tittlePattern, _ = regexp.Compile(`:\sStartseite\s(.*)\s-\sEinzelansicht`)
	lec.roomPattern, _ = regexp.Compile(`woch.*?-\s([A-Z][A-Z]\s\d*\.?\d+).*\sGruppe`)
	lec.modulesPattern, _ = regexp.Compile("BM 1|BM 2|BM 3|AM 1|AM 2|AM 3|AM 4|AM 5|VM 1|VM 2|VM 3|GM 1|GM 2|GM 3")
	lec.lecturersPattern, _ = regexp.Compile(`Zuständigkeit\s(.+?)\s(Studiengänge\sAbschluss|Zuordnung\szu)+`)
	lec.commentaryPattern, _ = regexp.Compile(`Inhalt\sKommentar(.*?)\s(Leistungsnachweis|Einsortiert in)`)
	return lec
}

type Lecture struct {
	day, time, title, flags, room, commentary, lectStr, link, textRaw string
	dayPattern, timePattern, tittlePattern, roomPattern,
	commentaryPattern, lecturersPattern, modulesPattern *regexp.Regexp
	modules, lecturers []string
}

func (l *Lecture) GetLecturers() string {
	if len(l.lecturers) > 1 {
		l.lectStr = strings.Join(l.lecturers, ", ")
	} else if len(l.lecturers) == 1 {
		l.lectStr = l.lecturers[0]
	}
	return l.lectStr
}
