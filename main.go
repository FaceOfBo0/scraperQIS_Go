package main

import "fmt"

func main() {
	// RunServer()

	const path = "source.xlsx"
	infos := getMetaInfos(path)
	for _, mi := range infos {
		fmt.Println(mi)
	}
	fmt.Println(len(infos))
}
