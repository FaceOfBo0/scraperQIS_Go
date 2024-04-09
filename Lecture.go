package main

import "regexp"

type Lecture struct {
	day, time, title, flags, room, commentary, lectStr, link, textRaw                                        string
	dayPattern, timePattern, tittlePattern, roomPattern, commentaryPattern, lecturersPattern, modulesPattern *regexp.Regexp
	modules, lecturers                                                                                       []string
}
