package service

import (
	"errors"
	"server/dao"
	"sort"
	"strconv"
	"sync"
)

type QueryType uint8

const (
	TitleQuery QueryType = iota
	AuthorQuery
	NotQuery
	PageQuery
)

type IndexScore struct {
	ID    uint64
	Score float32
}

func SearchArticle(queries map[QueryType]string, admin bool) (articles []dao.Article, err error) {
	invertedIndexes := make([]dao.InvertedIndex, 0)

	// page query got default value 1.
	page, err := strconv.Atoi(queries[PageQuery])
	if err != nil {
		return nil, errors.New("invalid page")
	}

	titleText := queries[TitleQuery]
	titleWords := queryTextToWord(titleText)
	authorText, authorTextOk := queries[AuthorQuery]
	authorWords := queryTextToWord(authorText)
	notText, notTextOk := queries[NotQuery]
	notWords := queryTextToWord(notText)

	// Try to get from cache.
	articles, ok := getCachedArticleRes(titleWords, authorWords, notWords, page, admin)
	if ok {
		return
	}

	// Get inverted-indexes by title queries.
	invertedIndexes = getInvertedIndexes(titleWords, word2article)
	if 2*len(titleWords) < len(invertedIndexes) {
		return nil, errors.New("no match records for your title queries")
	}

	// Get inverted-indexes by authors.
	if authorTextOk {
		authorInvertedIndexes := getInvertedIndexes(authorWords, word2author)

		authorIndexes, err := dao.Intersection(authorInvertedIndexes)
		if err != nil {
			return nil, err
		}

		if len(authorIndexes) == 0 {
			return nil, errors.New("no match records for your author queries")
		}

		strAuthorIndexes := make([]string, len(authorIndexes))

		for i := range strAuthorIndexes {
			strAuthorIndexes[i] = strconv.FormatUint(authorIndexes[i], 10)
		}

		_invertedIndexes := getInvertedIndexes(strAuthorIndexes, author2article)
		if err != nil {
			return nil, err
		}

		authorToArticleIndex, err := dao.Union(_invertedIndexes)
		if err != nil {
			return nil, err
		}

		invertedIndexes = append(invertedIndexes, authorToArticleIndex)
	}

	articleIndexes, err := dao.Intersection(invertedIndexes)
	if err != nil {
		return nil, err
	}

	// Del inverted-indexes by not words
	if notTextOk {
		notInvertedIndexes := getInvertedIndexes(notWords, word2article)
		notArticleIndexes, err := dao.Union(notInvertedIndexes)
		if err != nil {
			return nil, err
		}

		// Only reserve indexes which don't exist in notArticleIndexes.
		tempIndexes := articleIndexes
		articleIndexes = make([]uint64, 0)
		for _, id := range tempIndexes {
			if notArticleIndexes[id] == 0 {
				articleIndexes = append(articleIndexes, id)
			}
		}
	}

	indexScores := make([]IndexScore, len(articleIndexes))

	wg := sync.WaitGroup{}
	wg.Add(NumCpu)
	Len := len(indexScores)
	partialLen := (Len + NumCpu - 1) / NumCpu
	for i := 0; i < NumCpu; i++ {
		go func(j int) {
			defer wg.Done()
			t := j + partialLen
			for ; j < t && j < Len; j++ {
				indexScores[j].ID = articleIndexes[j]
				indexScores[j].Score = bm25(articleIndexes[j], invertedIndexes, word2article)
			}
		}(i * partialLen)
	}
	wg.Wait()

	// higher score first
	sort.Slice(indexScores, func(i, j int) bool {
		return indexScores[i].Score > indexScores[j].Score
	})

	// create a goroutine to do it.
	go cacheArticleRes(titleWords, authorWords, notWords, indexScores)

	for i := (page - 1) * 50; i < len(indexScores) && i < page*50; i++ {
		article, ok := getArticle(indexScores[i].ID, admin)
		if ok {
			articles = append(articles, article)
		}
	}

	return
}

//func GetAuthor(name string, limit, offset int) ([]Author, error) {
//	searchSQL := authorsTable.Limit(limit).Offset(offset).Where("name = ?", name)
//	_results := []dao.Author{}
//	err := searchSQL.Find(&_results).Error
//
//	if err != nil {
//		return nil, err
//	}
//
//	results := make([]Author, 0, len(_results))
//
//	// Turn work id to work title.
//	for _, _result := range _results {
//		worksBuilder := strings.Builder{}
//		workIDs := strings.Fields(_result.Articles)
//		for _, workID := range workIDs {
//			var title string
//			worksTable.Where("id = ?", workID).Select("title").Find(&title)
//			if worksBuilder.Len() != 0 {
//				worksBuilder.WriteString(", ")
//			}
//			worksBuilder.WriteString(title)
//		}
//		results = append(results, Author{_result.Name, worksBuilder.String(), _result.ArticleCount})
//	}
//	return results, nil
//}
//
//func GetTopAuthor(limit, offset int) (*[]Author, error) {
//	results := make([]Author, limit)
//	err := authorsTable.Select("name", "work_count").Order("work_count desc").Limit(limit).Offset(offset).Find(&results).Error
//	if err != nil {
//		return nil, err
//	}
//	return &results, nil
//}
