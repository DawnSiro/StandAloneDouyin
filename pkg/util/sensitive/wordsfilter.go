package sensitive

import "github.com/syyongx/go-wordsfilter"

func IsWordsFilter(text string) bool {
	wf := wordsfilter.New()
	root, _ := wf.GenerateWithFile("./pkg/util/sensitive/sensitive.txt")
	return wf.Contains(text, root)
}
