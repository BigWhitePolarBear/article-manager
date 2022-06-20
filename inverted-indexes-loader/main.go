package main

import (
	"github.com/sajari/fuzzy"
	"sync"
)

var (
	wg sync.WaitGroup

	spellChecker *fuzzy.Model
)

func main() {
	wg.Add(2)

	spellChecker = fuzzy.NewModel()
	spellChecker.SetDepth(2)
	spellChecker.Train(SampleEnglish())

	go articleLoader()
	go authorLoader()

	wg.Wait()
}
