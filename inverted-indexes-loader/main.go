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
	spellChecker = fuzzy.NewModel()
	spellChecker.SetDepth(2)
	spellChecker.Train(SampleEnglish())

	wg.Add(1)

	go articleLoader()
	// go authorLoader()

	wg.Wait()
}
