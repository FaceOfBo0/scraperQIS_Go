package main

import "fmt"

func main() {

	//fmt.Println(getHtmlText("https://qis.server.uni-frankfurt.de/qisserver/rds?state=verpublish&publishContainer=lectureInstList&publishid=80100"))

	scr := newScraper("")
	scr.getLecturesLinks("https://qis.server.uni-frankfurt.de/qisserver/rds?state=verpublish&publishContainer=lectureInstList&publishid=80100")
	//fmt.Println(scr.lecturesLinks)
	newLec := newLecture(scr.getLectureText(scr.lecturesLinks[1]))
	fmt.Printf("newLec.title: %v\n", newLec.title)

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
