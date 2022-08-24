package service

import (
	"context"
	"errors"
	"log"
	"server/dao"
	"sort"
	"strconv"
	"sync"
	"time"
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
	// page query got default value 1.
	page, err := strconv.Atoi(queries[PageQuery])
	if err != nil {
		return nil, errors.New("invalid page")
	}

	invertedIndexes := make([]dao.InvertedIndex, 0)

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

func SearchAuthor(name, _page string, admin bool) (authors []dao.Author, err error) {
	// page query got default value 1.
	page, err := strconv.Atoi(_page)
	if err != nil {
		return nil, errors.New("invalid page")
	}

	invertedIndexes := make([]dao.InvertedIndex, 0)

	nameWords := queryTextToWord(name)

	// Try to get from cache.
	authors, ok := getCachedAuthorRes(nameWords, page, admin)
	if ok {
		return
	}

	// Get inverted-indexes by title queries.
	invertedIndexes = getInvertedIndexes(nameWords, word2author)
	if 2*len(nameWords) < len(invertedIndexes) {
		return nil, errors.New("no match records for your name queries")
	}

	authorIndexes, err := dao.Intersection(invertedIndexes)
	if err != nil {
		return nil, err
	}

	indexScores := make([]IndexScore, len(authorIndexes))

	wg := sync.WaitGroup{}
	wg.Add(NumCpu)
	Len := len(indexScores)
	partialLen := (Len + NumCpu - 1) / NumCpu
	for i := 0; i < NumCpu; i++ {
		go func(j int) {
			defer wg.Done()
			t := j + partialLen
			for ; j < t; j++ {
				indexScores[j].ID = authorIndexes[j]
				indexScores[j].Score = bm25(authorIndexes[j], invertedIndexes, word2author)
			}
		}(i * partialLen)
	}
	wg.Wait()

	// higher score first
	sort.Slice(indexScores, func(i, j int) bool {
		return indexScores[i].Score > indexScores[j].Score
	})

	// Create a goroutine to cache the results.
	go cacheAuthorRes(nameWords, indexScores)

	for i := (page - 1) * 50; i < len(indexScores) && i < page*50; i++ {
		author, ok := getAuthor(indexScores[i].ID, admin)
		if ok {
			authors = append(authors, author)
		}
	}

	return
}

func GetTopAuthor(page uint64, admin bool) (authors []dao.Author, err error) {
	if page*10 > uint64(dao.AuthorCnt) {
		return nil, errors.New("there are no so many authors")
	}

	sPage := strconv.FormatUint(page, 10)

	authors = make([]dao.Author, 0, 10)

	// Try to get from cache.
	IDs := dao.TopAuthorResRDB.LRange(context.Background(), sPage, 0, -1).Val()
	if len(IDs) > 0 {
		for _, _id := range IDs {
			id, err := strconv.ParseUint(_id, 10, 64)
			if err != nil {
				log.Println("service/get.go GetTopAuthor error:", err)
				continue
			}
			author, ok := getAuthor(id, admin)
			if !ok {
				continue
			}
			authors = append(authors, author)
		}
		return authors, nil
	}

	// cache missed
	offset := 10 * (page - 1)
	limit := 10
	err = dao.DB.Model(&dao.Author{}).Select("id").
		Order("article_count desc").Offset(int(offset)).Limit(int(limit)).Find(&IDs).Error
	if err != nil || len(IDs) == 0 {
		log.Println("service/get.go GetTopAuthor error:", err)
		return nil, nil
	}

	err = dao.TopAuthorResRDB.RPush(context.Background(), sPage, IDs).Err()
	if err != nil {
		log.Println("service/get.go GetTopAuthor error:", err)
	}

	var expireErr error
	expireErr = errors.New("")
	// Retry 5 times.
	for i := 0; expireErr != nil && i < 5; i++ {
		expireErr = dao.TopAuthorResRDB.Expire(context.Background(), sPage, time.Minute).Err()
	}
	if expireErr != nil {
		log.Println("service/get.go GetTopAuthor error:", expireErr)
	}

	for _, _id := range IDs {
		id, err := strconv.ParseUint(_id, 10, 64)
		if err != nil {
			log.Println("service/get.go GetTopAuthor error:", err)
			continue
		}
		author, ok := getAuthor(id, admin)
		if !ok {
			continue
		}
		authors = append(authors, author)
	}
	return authors, nil
}
