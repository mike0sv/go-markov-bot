package sliceutils

import "math/rand"

func Shuffle(list *[]string) {
	lst := *list
	rand.Shuffle(len(lst), func(i, j int) { lst[i], lst[j] = lst[j], lst[i] })
}
