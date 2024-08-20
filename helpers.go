package main

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
