package util

import "github.com/syyongx/go-wordsfilter"

func IsWordsFilter(text string) bool {
	wf := wordsfilter.New()
	root, _ := wf.GenerateWithFile("sensitive.txt")
	return wf.Contains(text, root)
}
