package main

import "regexp"

func NewLecture(text string, url string) *Lecture {
	lec := Lecture{textRaw: text, link: url}
	lec.dayPattern, _ = regexp.Compile(`([A-Z][a-z])\.\s\d+:\d+`)
	lec.timePattern, _ = regexp.Compile(`\d+:\d+`)
	return &lec
}
