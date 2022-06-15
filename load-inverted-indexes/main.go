package main

import (
	"strconv"
	"sync"
)

var (
	wg sync.WaitGroup

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

	go articleLoader()
	go authorLoader()

	wg.Wait()
}
