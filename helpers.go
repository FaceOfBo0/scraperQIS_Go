package main

import "os"

type Tuple[V, U any] struct {
	fst V
	snd U
}

func mapList[I, O any](list []I, fn func(I) O) []O {
	outList := make([]O, len(list))
	for i, elem := range list {
		outList[i] = fn(elem)
	}
	return outList
}

func replaceIdx(str string, replacement string, index int) string {
	return str[:index] + replacement + str[index+1:]
}

// saveStrToFile saves the JSON string to a specified file.
func saveStrToFile(jsonStr, filePath string) error {
	// Open the file for writing, create it if it doesn't exist
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the JSON string to the file
	_, err = file.WriteString(jsonStr)
	if err != nil {
		return err
	}

	return nil
}
