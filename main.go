package main

import "fmt"

func main() {
	lec := NewLecture("re", "re")
	fmt.Println(lec.commentaryPattern.String())
	//RunTutServer()
}
