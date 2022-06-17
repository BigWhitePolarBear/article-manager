package main

import (
	"github.com/sajari/fuzzy"
	"strconv"
	"sync"
)

var (
	wg sync.WaitGroup

	spellChecker *fuzzy.Model

	articleCount int64
	authorCount  int64
)

func main() {
	// Store the counts of article and author into database
	DB.Table("articles").Count(&articleCount)
	DB.Table("authors").Count(&authorCount)

	DB.Table("variables").Create(variable{Key: "ArticleCounts", Value: strconv.FormatInt(articleCount, 16)})
	DB.Table("variables").Create(variable{Key: "AuthorCounts", Value: strconv.FormatInt(authorCount, 16)})

	wg.Add(2)

	spellChecker = fuzzy.NewModel()
	spellChecker.SetDepth(2)
	spellChecker.Train(SampleEnglish())

	go articleLoader()
	go authorLoader()

	wg.Wait()
}
