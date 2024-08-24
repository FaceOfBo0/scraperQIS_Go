package main

import (
	"fmt"
	"time"
)

func main() {

	scr := newScraper("https://qis.server.uni-frankfurt.de/qisserver/rds?state=verpublish&publishContainer=lectureInstList&publishid=80100", "")

	/* start := time.Now()
	scr.loadLecturesLinks()
	duration := time.Since(start)
	fmt.Printf("scr.lecturesLinks: %d\n", len(scr.lecturesLinks))
	fmt.Printf("duration: %v\n", duration) */

	start := time.Now()
	test := scr.getLecturesv2()
	duration := time.Since(start)
	fmt.Println(test[0].time)
	fmt.Println(duration)
	/* fmt.Print("Press any key...")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() */

	/* scr.getLecturesLinks("https://qis.server.uni-frankfurt.de/qisserver/rds?state=verpublish&publishContainer=lectureInstList&publishid=80100")
	//fmt.Println(scr.lecturesLinks)
	newLec := newLecture(scr.getLectureText(scr.lecturesLinks[3]))
	fmt.Printf("newLec.textRaw: %v\n", newLec.textRaw)
	fmt.Printf("newLec.title: %v\n", newLec.title)
	fmt.Printf("newLec.lecturers: %v\n", newLec.lecturers)
	fmt.Printf("newLec.day: %v\n", newLec.day)
	fmt.Printf("newLec.time: %v\n", newLec.time)
	fmt.Printf("newLec.room: %v\n", newLec.room)
	fmt.Printf("newLec.modules: %v\n", newLec.modules)
	fmt.Printf("newLec.commentary: %v\n", newLec.commentary)
	fmt.Printf("newLec.flags: %v\n", newLec.flags) */

	/* for _, elem := range newLec.dayPattern.FindStringSubmatch(newLec.textRaw) {
		fmt.Printf("elem: %v\n", elem)
	} */

	/* lec := newLecture("re", "re")
	fmt.Println(lec.commentaryPattern.String()) */
	//RunTutServer()

	/* 	r, _ := regexp.Compile(`\w+`)
	   	newStr := r.FindAll(content, -1)
	   	fmt.Printf("newStr: %v\n", newStr)
	   	strList := mapList(newStr, func(elem []byte) string {
	   		return string(elem)
	   	})
	   	fmt.Printf("strList: %v\n", strList) */

	/* 	patt, _ := regexp.Compile(`\d+:\d+`)
	   	testStr := "hello 122:243 3 teststaata"
	   	testMatch := patt.Find([]byte(testStr))
	   	rslt := string(testMatch)
	   	if len(rslt) == 0 {
	   		fmt.Println("No match, its empty!")
	   	} else {
	   		fmt.Println(rslt)
	   	} */

}
